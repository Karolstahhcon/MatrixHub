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

package handler

import (
	"context"
	"errors"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/auth"
	authzmocks "github.com/matrixhub-ai/matrixhub/internal/domain/authz/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
	modelmocks "github.com/matrixhub-ai/matrixhub/internal/domain/model/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
)

// Helpers and fixtures below are shared by the handler tests in this package.

var (
	handlerTestCreatedAt = time.Date(2026, 3, 1, 8, 30, 0, 0, time.UTC)
	handlerTestUpdatedAt = time.Date(2026, 3, 2, 9, 45, 0, 0, time.UTC)

	handlerTestInsecure = grpc.WithTransportCredentials(insecure.NewCredentials())
)

const (
	handlerTestCreatedAtText = "2026-03-01T08:30:00Z"
	handlerTestUpdatedAtText = "2026-03-02T09:45:00Z"
)

type handlerTestErrorCase struct {
	name     string
	err      error
	wantCode codes.Code
	// wantMsg is only checked when set.
	wantMsg string
}

func requireHandlerTestStatus(t *testing.T, err error, wantCode codes.Code, wantMsg string) {
	t.Helper()
	require.Equal(t, wantCode, status.Code(err), "unexpected error: %v", err)
	if wantMsg != "" {
		require.Equal(t, wantMsg, status.Convert(err).Message())
	}
}

// requireProtoEqual uses proto.Equal because reflect-based equality also
// compares the internal state of generated messages.
func requireProtoEqual(t *testing.T, want, got proto.Message) {
	t.Helper()
	require.Truef(t, proto.Equal(want, got), "proto mismatch\nwant: %v\ngot:  %v", want, got)
}

// registerHandlerForTest registers h on a fresh gRPC server and grpc-gateway
// mux the way APIServer does, with gRPC served over an in-memory listener.
// RegisterToServer dials the gateway client with context.Background(), so the
// generated gateway code never closes it; the helper captures that client from
// the calls it makes and closes it when the test ends.
func registerHandlerForTest(t *testing.T, h IHandler, dialOpts ...grpc.DialOption) (*grpc.Server, *runtime.ServeMux) {
	t.Helper()
	listener := bufconn.Listen(1 << 20)
	grpcServer := grpc.NewServer()
	gatewayMux := runtime.NewServeMux()
	dialer := grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
		return listener.DialContext(ctx)
	})
	var (
		clientsMu sync.Mutex
		clients   = map[*grpc.ClientConn]struct{}{}
	)
	trackClients := grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		clientsMu.Lock()
		clients[cc] = struct{}{}
		clientsMu.Unlock()
		return invoker(ctx, method, req, reply, cc, opts...)
	})
	h.RegisterToServer(&ServerOptions{
		GatewayMux:  gatewayMux,
		GRPCServer:  grpcServer,
		GRPCAddr:    "passthrough:///bufconn",
		GRPCDialOpt: append([]grpc.DialOption{dialer, trackClients}, dialOpts...),
	})
	go func() {
		_ = grpcServer.Serve(listener)
	}()
	t.Cleanup(func() {
		clientsMu.Lock()
		defer clientsMu.Unlock()
		for cc := range clients {
			_ = cc.Close()
		}
		grpcServer.Stop()
	})
	return grpcServer, gatewayMux
}

func serveHandlerTestRequest(t *testing.T, handler http.Handler, method, target string) *httptest.ResponseRecorder {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequestWithContext(ctx, method, target, nil))
	return rec
}

func handlerTestCommit() *git.Commit {
	return &git.Commit{
		ID:             "8f3c2a1",
		Message:        "Add data shard",
		AuthorName:     "Alice",
		AuthorEmail:    "alice@example.com",
		AuthorDate:     time.Date(2026, 3, 4, 1, 2, 3, 0, time.UTC),
		CommitterName:  "Bob",
		CommitterEmail: "bob@example.com",
		CommitterDate:  time.Date(2026, 3, 4, 4, 5, 6, 0, time.UTC),
		Diff:           "diff --git a/config.json b/config.json",
		ParentCommit:   "1a2b3c4",
		CreatedAt:      time.Date(2026, 3, 4, 7, 8, 9, 0, time.UTC),
	}
}

func handlerTestCommitProto() *v1alpha1.Commit {
	return &v1alpha1.Commit{
		Id:             "8f3c2a1",
		Message:        "Add data shard",
		AuthorName:     "Alice",
		AuthorEmail:    "alice@example.com",
		AuthorDate:     "2026-03-04T01:02:03Z",
		CommitterName:  "Bob",
		CommitterEmail: "bob@example.com",
		CommitterDate:  "2026-03-04T04:05:06Z",
		Diff:           "diff --git a/config.json b/config.json",
		CreatedAt:      "2026-03-04T07:08:09Z",
	}
}

func handlerTestTreeEntries() []*git.TreeEntry {
	return []*git.TreeEntry{
		{Name: "data", Type: git.FileTypeDir, Path: "data"},
		{
			Name:   "shard-00000.bin",
			Type:   git.FileTypeFile,
			Size:   4_294_967_296,
			Path:   "data/shard-00000.bin",
			Hash:   "4f6b2c3d",
			IsLFS:  true,
			URL:    "/demo/sample/resolve/main/data/shard-00000.bin",
			Commit: handlerTestCommit(),
		},
	}
}

func handlerTestTreeFiles() []*v1alpha1.File {
	return []*v1alpha1.File{
		{Name: "data", Type: v1alpha1.FileType_DIR, Path: "data"},
		{
			Name:   "shard-00000.bin",
			Type:   v1alpha1.FileType_FILE,
			Size:   4_294_967_296,
			Path:   "data/shard-00000.bin",
			Sha256: "4f6b2c3d",
			Lfs:    true,
			Url:    "/demo/sample/resolve/main/data/shard-00000.bin",
			Commit: handlerTestCommitProto(),
		},
	}
}

const (
	modelHandlerTestProject = "demo"
	modelHandlerTestName    = "llama"
)

func modelHandlerTestSetup(t *testing.T) (*ModelHandler, *modelmocks.MockIModelService, *authzmocks.MockIAuthzService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := modelmocks.NewMockIModelService(ctrl)
	authzService := authzmocks.NewMockIAuthzService(ctrl)
	h, ok := NewModelHandler(svc, authzService).(*ModelHandler)
	require.True(t, ok)
	return h, svc, authzService
}

func modelHandlerTestUserContext() context.Context {
	return auth.WithIdentity(context.Background(), user.NewUserIdentity(42, "alice"))
}

// modelHandlerTestAuthorized returns a handler whose caller holds perm on the
// test project, so a test can focus on what happens after authorization.
func modelHandlerTestAuthorized(t *testing.T, perm role.Permission) (context.Context, *ModelHandler, *modelmocks.MockIModelService) {
	t.Helper()
	h, svc, authzService := modelHandlerTestSetup(t)
	ctx := modelHandlerTestUserContext()
	authzService.EXPECT().VerifyProjectPermissionByName(ctx, modelHandlerTestProject, perm).Return(true, nil)
	return ctx, h, svc
}

func modelHandlerTestModel() *model.Model {
	syncedAt := time.Date(2026, 3, 3, 10, 0, 0, 0, time.UTC)
	return &model.Model{
		ID:             7,
		Name:           modelHandlerTestName,
		ProjectID:      3,
		ProjectName:    modelHandlerTestProject,
		Size:           1536,
		DefaultBranch:  "main",
		ParameterCount: 7_000_000_000,
		ReadmeContent:  "# Llama",
