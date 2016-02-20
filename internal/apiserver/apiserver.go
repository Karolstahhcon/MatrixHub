// Copyright The MatrixHub Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apiserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/matrixhub-ai/hfd/pkg/authenticate"
	"github.com/matrixhub-ai/hfd/pkg/lfs"
	"github.com/matrixhub-ai/hfd/pkg/mirror"
	"github.com/matrixhub-ai/hfd/pkg/permission"
	"github.com/matrixhub-ai/hfd/pkg/receive"
	hfdssh "github.com/matrixhub-ai/hfd/pkg/ssh"
	gitstorage "github.com/matrixhub-ai/hfd/pkg/storage"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/matrixhub-ai/matrixhub/internal/apiserver/handler"
	backendhf "github.com/matrixhub-ai/matrixhub/internal/apiserver/handler/hf"
	backendhttp "github.com/matrixhub-ai/matrixhub/internal/apiserver/handler/http"
	backendlfs "github.com/matrixhub-ai/matrixhub/internal/apiserver/handler/lfs"
	backendssh "github.com/matrixhub-ai/matrixhub/internal/apiserver/handler/ssh"
	"github.com/matrixhub-ai/matrixhub/internal/apiserver/middleware"
	"github.com/matrixhub-ai/matrixhub/internal/domain/authz"
	"github.com/matrixhub-ai/matrixhub/internal/domain/cleanup"
	"github.com/matrixhub-ai/matrixhub/internal/domain/dataset"
	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
	"github.com/matrixhub-ai/matrixhub/internal/domain/preset"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registrydiscovery"
	"github.com/matrixhub-ai/matrixhub/internal/domain/syncjob"
	"github.com/matrixhub-ai/matrixhub/internal/domain/syncpolicy"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/config"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
	"github.com/matrixhub-ai/matrixhub/internal/jobserver"
	"github.com/matrixhub-ai/matrixhub/internal/jobserver/canceller"
	"github.com/matrixhub-ai/matrixhub/internal/jobserver/logstore"
	"github.com/matrixhub-ai/matrixhub/internal/repo"
	hfdiscovery "github.com/matrixhub-ai/matrixhub/internal/repo/registrydiscovery/hf"
	mhdiscovery "github.com/matrixhub-ai/matrixhub/internal/repo/registrydiscovery/matrixhub"
)

const maxGrpcMsgSize = 100 * 1024 * 1024

type APIServer struct {
	config     *config.Config
	debug      bool
	cmux       cmux.CMux
	httpServer *http.Server
	sshServer  *backendssh.Server
	engine     *gin.Engine
	gatewayMux *runtime.ServeMux
	grpcServer *grpc.Server
	port       int

	gitHooks   gitHooks
	gitAuth    gitAuth
	gitStorage gitStorage

	repos    *repo.Repos
	services *Services
	handlers []handler.IHandler

	jobServer *jobserver.JobServer
	jobCancel context.CancelFunc
	jobWait   sync.WaitGroup
}

func NewAPIServer(config *config.Config) *APIServer {
	if config.APIServer == nil {
		panic("apiserver config is nil")
	}

	engine := gin.New()
	engine.Use(
		gin.Recovery(),
	)

	gatewayMux := runtime.NewServeMux(
		runtime.WithForwardResponseOption(middleware.ResponseHeaderLocation),
		runtime.WithOutgoingHeaderMatcher(middleware.HeaderMatcher),
	)

	httpServer := &http.Server{
		Handler:           engine,
		ReadHeaderTimeout: 30 * time.Second,
	}

	server := &APIServer{
		config:     config,
		debug:      config.Debug,
		httpServer: httpServer,
		engine:     engine,
		gatewayMux: gatewayMux,
		port:       config.APIServer.Port,
	}

	server.initMirrorHooks()
	server.initGitStorage()
	server.initHandlersServicesRepos()
	server.initGitHooks()
	server.initGitAuth()
	server.initSSHBackend()

	// Register authn + authz middleware (must be after initHandlersServicesRepos)
	streamMiddleware := []grpc.StreamServerInterceptor{
		grpc_recovery.StreamServerInterceptor(),
	}
	unaryMiddleware := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(middleware.RecoveryFunc)),
		middleware.AuthInterceptor(server.repos.Session, server.repos.User, server.repos.AccessToken, server.repos.Robot),
		middleware.AuthzInterceptor(server.services.Authz.VerifyPlatformPermission),
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			streamMiddleware...,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			unaryMiddleware...,
		)),
		grpc.MaxSendMsgSize(maxGrpcMsgSize),
		grpc.MaxRecvMsgSize(maxGrpcMsgSize),
	)
	server.grpcServer = grpcServer

	server.httpServer.Handler = server.initBackends(server.httpServer.Handler)
	server.registerRoutersAndHandlers()
	server.initPresets()

	return server
}

func (server *APIServer) initPresets() {
	if !server.config.Presets.IsEnabled() {
		return
	}
	service := preset.NewService(true, preset.DefaultManifest(), server.repos.Preset,
		server.repos.Project, server.repos.Model, server.repos.Git)
	ctx := context.Background()
	service.Apply(ctx)
	service.Reconcile(ctx)
}

type gitHooks struct {
	permissionHookFunc    func(ctx context.Context, op permission.Operation, repoName string, opCtx permission.Context) (bool, error)
	preReceiveHookFunc    func(ctx context.Context, repoName string, updates []receive.RefUpdate) (bool, error)
	postReceiveHookFunc   func(ctx context.Context, repoName string, updates []receive.RefUpdate) error
	mirrorSourceFunc      func(ctx context.Context, repoName string) (string, bool, error)
	mirrorDestinationFunc func(ctx context.Context, repoName string) (string, bool, error)
	mirrorRefFilterFunc   func(ctx context.Context, repoName string, remoteRefs []string) ([]string, error)
}

func (server *APIServer) initGitHooks() {
	permissionHookFunc := middleware.NewRepoEnforcer(server.services.Authz)

	preReceiveHookFunc := func(ctx context.Context, repoName string, updates []receive.RefUpdate) (bool, error) {
		repoType, project, name, ok := utils.ParseFromRepoName(repoName)
		if !ok {
			return false, nil
		}
		if repoType == "models" {
			_, err := server.services.Model.EnsureModel(ctx, project, name)
			return err == nil, err
		}

		return false, nil
	}

	// Refresh database metadata (README, size, labels) after repository content
	// changed. Metadata lives in the database, so without this a pushed or
	// mirrored repository shows up with an empty description and zero size.
	// Failures must not abort the push or the mirror sync, so they are logged only.
	postReceiveHookFunc := func(ctx context.Context, repoName string, updates []receive.RefUpdate) error {
		repoType, project, name, ok := utils.ParseFromRepoName(repoName)
		if !ok {
			return nil
		}

		var err error
		switch repoType {
		case "models":
			err = server.services.Model.SyncMetadata(ctx, project, name)
		case "datasets":
			err = server.services.Dataset.SyncMetadata(ctx, project, name)
		default:
			return nil
		}
		if err != nil {
			log.Warnw("sync metadata after receive failed", "repo", repoName, "error", err)
		}
		return nil
	}

	server.gitHooks.permissionHookFunc = permissionHookFunc
	server.gitHooks.preReceiveHookFunc = preReceiveHookFunc
	server.gitHooks.postReceiveHookFunc = postReceiveHookFunc

}

func (server *APIServer) initMirrorHooks() {
	mirrorSourceFunc := func(ctx context.Context, repoName string) (string, bool, error) {
		// return baseURL + "/" + repoName, true, nil
		return "", false, nil
	}

	mirrorDestinationFunc := func(ctx context.Context, repoName string) (string, bool, error) {
		return "", false, nil
	}

	mirrorRefFilterFunc := func(ctx context.Context, repoName string, remoteRefs []string) ([]string, error) {
		filteredRefs := []string{}
		for _, ref := range remoteRefs {
			if strings.HasPrefix(ref, "refs/heads/") || strings.HasPrefix(ref, "refs/tags/") {
				filteredRefs = append(filteredRefs, ref)
			}
		}
		return filteredRefs, nil
	}
	server.gitHooks.mirrorSourceFunc = mirrorSourceFunc
	server.gitHooks.mirrorDestinationFunc = mirrorDestinationFunc
	server.gitHooks.mirrorRefFilterFunc = mirrorRefFilterFunc
}

type gitAuth struct {
	basicAuthValidator authenticate.BasicAuthValidator
	publicKeyValidator authenticate.PublicKeyValidator
	tokenValidator     authenticate.TokenValidator
	tokenSignValidator authenticate.TokenSignValidator
}

func (server *APIServer) initGitAuth() {
	basicAuthValidator := middleware.GitBasicAuthAuthn(server.repos.AccessToken, server.repos.User, server.repos.Robot)
	publicKeyValidator := middleware.GitPublicKeyAuthn(server.repos.SSHKey, server.repos.User)
	tokenValidator := middleware.GitHTTPAuthn(server.repos.AccessToken, server.repos.User, server.repos.Robot)

	// TODO: Use a proper secret management solution to manage the token signing secret.
	// Generating and validating temporary tokens.
	// Currently only used to provide http lfs download in ssh ports.
	tmpTokenSecret := []byte("secret-xxxxxx")
	tokenSignValidator := authenticate.NewTokenSignValidator(tmpTokenSecret)

	server.gitAuth.basicAuthValidator = basicAuthValidator
	server.gitAuth.publicKeyValidator = publicKeyValidator
	server.gitAuth.tokenValidator = tokenValidator
	server.gitAuth.tokenSignValidator = tokenSignValidator
}

type gitStorage struct {
	storage      *gitstorage.Storage
	lfsStorage   lfs.Storage
	sharedMirror *mirror.Mirror
}

func (server *APIServer) initGitStorage() {
	storage := gitstorage.NewStorage(
		gitstorage.WithRootDir(server.config.DataDir),
	)

	lfsStorage := lfs.NewLocal(storage.LFSDir())

	mirrorSourceFunc := server.gitHooks.mirrorSourceFunc
	mirrorDestinationFunc := server.gitHooks.mirrorDestinationFunc
	mirrorRefFilterFunc := server.gitHooks.mirrorRefFilterFunc
	preReceiveHookFunc := server.gitHooks.preReceiveHookFunc

	sharedMirror := mirror.NewMirror(
		mirror.WithMirrorSourceFunc(mirrorSourceFunc),
		mirror.WithMirrorDestinationFunc(mirrorDestinationFunc),
		mirror.WithMirrorRefFilterFunc(mirrorRefFilterFunc),
		mirror.WithPreReceiveHookFunc(preReceiveHookFunc),
		mirror.WithLFSStorage(lfsStorage),
		mirror.WithPushXET(true),
		mirror.WithPullXET(false),
		mirror.WithConcurrency(2),
		mirror.WithCacheDir(storage.TmpDir()),
	)

	server.gitStorage.storage = storage
	server.gitStorage.lfsStorage = lfsStorage
	server.gitStorage.sharedMirror = sharedMirror
}

func (server *APIServer) initBackends(handler http.Handler) http.Handler {
	storage := server.gitStorage.storage
	lfsStorage := server.gitStorage.lfsStorage
	sharedMirror := server.gitStorage.sharedMirror
	permissionHookFunc := server.gitHooks.permissionHookFunc
	preReceiveHookFunc := server.gitHooks.preReceiveHookFunc
	postReceiveHookFunc := server.gitHooks.postReceiveHookFunc
	basicAuthValidator := server.gitAuth.basicAuthValidator
	tokenValidator := server.gitAuth.tokenValidator
	tokenSignValidator := server.gitAuth.tokenSignValidator

	handler = backendhf.NewHandler(
		backendhf.WithStorage(storage),
		backendhf.WithNext(handler),
		backendhf.WithMirror(sharedMirror),
		backendhf.WithPermissionHookFunc(permissionHookFunc),
		backendhf.WithPreReceiveHookFunc(preReceiveHookFunc),
		backendhf.WithPostReceiveHookFunc(postReceiveHookFunc),
		backendhf.WithLFSStorage(lfsStorage),
		backendhf.WithMiddlewares(
			middleware.HFAuthnMiddleware(server.repos.AccessToken, server.repos.Session, server.repos.User, server.repos.Robot),
		),
		backendhf.WithServices(server.services.Model, server.repos.Git, server.services.Authz),
	)

