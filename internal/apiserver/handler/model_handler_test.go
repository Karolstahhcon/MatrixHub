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
		IsPopular:      true,
		Labels: []model.Label{
			{ID: 1, Name: "text-generation", Category: "task", Scope: "model", CreatedAt: handlerTestCreatedAt, UpdatedAt: handlerTestUpdatedAt},
			{ID: 2, Name: "apache-2.0", Category: "license", Scope: "model", CreatedAt: handlerTestCreatedAt, UpdatedAt: handlerTestUpdatedAt},
		},
		SyncedAt:  &syncedAt,
		CreatedAt: handlerTestCreatedAt,
		UpdatedAt: handlerTestUpdatedAt,
	}
}

func modelHandlerTestModelProto() *v1alpha1.Model {
	return &v1alpha1.Model{
		Id:            7,
		Name:          modelHandlerTestName,
		DefaultBranch: "main",
		SyncedAt:      "2026-03-03T10:00:00Z",
		CreatedAt:     handlerTestCreatedAtText,
		UpdatedAt:     handlerTestUpdatedAtText,
		CloneUrls:     &v1alpha1.CloneUrls{},
		Labels: []*v1alpha1.Label{
			{Id: 1, Name: "text-generation", Category: v1alpha1.Category_TASK, CreatedAt: handlerTestCreatedAtText, UpdatedAt: handlerTestUpdatedAtText},
			{Id: 2, Name: "apache-2.0", Category: v1alpha1.Category_LICENSE, CreatedAt: handlerTestCreatedAtText, UpdatedAt: handlerTestUpdatedAtText},
		},
		Project:        modelHandlerTestProject,
		ReadmeContent:  "# Llama",
		Size:           "1536",
		ParameterCount: "7000000000",
		Popular:        true,
	}
}

func TestNewModelHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	svc := modelmocks.NewMockIModelService(ctrl)
	authzService := authzmocks.NewMockIAuthzService(ctrl)

	h, ok := NewModelHandler(svc, authzService).(*ModelHandler)
	require.True(t, ok)
	require.Same(t, svc, h.ms)
	require.Same(t, authzService, h.authzService)
}

func TestFormatSizeAndParameterCount(t *testing.T) {
	cases := []struct {
		name  string
		value int64
		want  string
	}{
		{name: "zero", value: 0, want: "0"},
		{name: "kilobyte range is not converted", value: 1536, want: "1536"},
		{name: "billions are not abbreviated", value: 7_000_000_000, want: "7000000000"},
		{name: "max int64", value: math.MaxInt64, want: "9223372036854775807"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, formatSize(tc.value))
			require.Equal(t, tc.want, formatParameterCount(tc.value))
		})
	}
}

func TestLabelCategoryToProto(t *testing.T) {
	cases := []struct {
		category string
		want     v1alpha1.Category
	}{
		{category: "task", want: v1alpha1.Category_TASK},
		{category: "library", want: v1alpha1.Category_LIBRARY},
		{category: "license", want: v1alpha1.Category_LICENSE},
		{category: "language", want: v1alpha1.Category_LANGUAGE},
		{category: "other", want: v1alpha1.Category_OTHER},
		{category: "", want: v1alpha1.Category_OTHER},
		{category: "quantization", want: v1alpha1.Category_OTHER},
	}
	for _, tc := range cases {
		t.Run("category="+tc.category, func(t *testing.T) {
			require.Equal(t, tc.want, labelCategoryToProto(tc.category))
		})
	}
}

func TestModelToProto(t *testing.T) {
	t.Run("maps every exposed field", func(t *testing.T) {
		requireProtoEqual(t, modelHandlerTestModelProto(), modelToProto(modelHandlerTestModel()))
	})

	t.Run("keeps optional fields empty", func(t *testing.T) {
		got := modelToProto(&model.Model{
			ID:          8,
			Name:        "bare",
			ProjectName: modelHandlerTestProject,
			CreatedAt:   handlerTestCreatedAt,
			UpdatedAt:   handlerTestUpdatedAt,
		})

		requireProtoEqual(t, &v1alpha1.Model{
			Id:             8,
			Name:           "bare",
			CreatedAt:      handlerTestCreatedAtText,
			UpdatedAt:      handlerTestUpdatedAtText,
			CloneUrls:      &v1alpha1.CloneUrls{},
			Project:        modelHandlerTestProject,
			Size:           "0",
			ParameterCount: "0",
		}, got)
		require.Empty(t, got.Labels)
	})
}

func TestRevisionsToProto(t *testing.T) {
	t.Run("keeps branch and tag order", func(t *testing.T) {
		got := revisionsToProto(&git.Revisions{
			Branches: []*git.Revision{{Name: "main"}, {Name: "dev"}},
			Tags:     []*git.Revision{{Name: "v1.0"}},
		})

		requireProtoEqual(t, &v1alpha1.Revisions{
			Branches: []*v1alpha1.Revision{{Name: "main"}, {Name: "dev"}},
			Tags:     []*v1alpha1.Revision{{Name: "v1.0"}},
		}, got)
	})

	t.Run("handles a repository without refs", func(t *testing.T) {
		got := revisionsToProto(&git.Revisions{})
		require.Empty(t, got.Branches)
		require.Empty(t, got.Tags)
	})
}

func TestTreeEntryToProtoFile(t *testing.T) {
	entries := handlerTestTreeEntries()
	files := handlerTestTreeFiles()

	t.Run("maps a file with its last commit", func(t *testing.T) {
		requireProtoEqual(t, files[1], treeEntryToProtoFile(entries[1]))
	})

	t.Run("maps a directory without commit", func(t *testing.T) {
		got := treeEntryToProtoFile(entries[0])
		requireProtoEqual(t, files[0], got)
		require.Nil(t, got.Commit)
	})
}

func TestModelHandlerRegisterToServer(t *testing.T) {
	t.Run("serves models over gRPC and the HTTP gateway", func(t *testing.T) {
		h, svc, authzService := modelHandlerTestSetup(t)
		authzService.EXPECT().
			VerifyProjectPermissionByName(gomock.Any(), modelHandlerTestProject, role.ModelGet).
			Return(true, nil)
		svc.EXPECT().
			GetModel(gomock.Any(), modelHandlerTestProject, modelHandlerTestName).
			Return(modelHandlerTestModel(), nil)

		grpcServer, gatewayMux := registerHandlerForTest(t, h, handlerTestInsecure)
		require.Contains(t, grpcServer.GetServiceInfo(), "matrixhub.v1alpha1.Models")

		rec := serveHandlerTestRequest(t, gatewayMux, http.MethodGet, "/api/v1alpha1/models/demo/llama")
		require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
		got := &v1alpha1.Model{}
		require.NoError(t, protojson.Unmarshal(rec.Body.Bytes(), got))
		requireProtoEqual(t, modelHandlerTestModelProto(), got)
	})

	t.Run("keeps the gRPC service when gateway registration fails", func(t *testing.T) {
		h, _, _ := modelHandlerTestSetup(t)

		// Without transport credentials the gateway cannot create its client.
		grpcServer, gatewayMux := registerHandlerForTest(t, h)
		require.Contains(t, grpcServer.GetServiceInfo(), "matrixhub.v1alpha1.Models")

		rec := serveHandlerTestRequest(t, gatewayMux, http.MethodGet, "/api/v1alpha1/models/demo/llama")
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestModelHandlerGetAccessibleProjectIDs(t *testing.T) {
	errAuthz := errors.New("authz backend unavailable")
	cases := []struct {
		name    string
		ctx     context.Context
		expect  func(ctx context.Context, authzService *authzmocks.MockIAuthzService)
		want    []int
		wantErr error
	}{
		{
			name: "anonymous caller is limited to public projects",
			ctx:  context.Background(),
			expect: func(ctx context.Context, authzService *authzmocks.MockIAuthzService) {
				authzService.EXPECT().GetUserAccessibleProjectIDs(ctx, 0).Return([]int{1, 2}, nil)
			},
			want: []int{1, 2},
		},
		{
			name: "anonymous lookup failure is returned",
			ctx:  context.Background(),
			expect: func(ctx context.Context, authzService *authzmocks.MockIAuthzService) {
				authzService.EXPECT().GetUserAccessibleProjectIDs(ctx, 0).Return(nil, errAuthz)
			},
			wantErr: errAuthz,
		},
		{
			name: "platform administrator is not restricted",
			ctx:  modelHandlerTestUserContext(),
			expect: func(ctx context.Context, authzService *authzmocks.MockIAuthzService) {
				authzService.EXPECT().VerifyPlatformPermission(ctx, role.Permission("*.*")).Return(true, nil)
			},
			want: nil,
		},
		{
			name: "regular user is limited to accessible projects",
			ctx:  modelHandlerTestUserContext(),
			expect: func(ctx context.Context, authzService *authzmocks.MockIAuthzService) {
				authzService.EXPECT().VerifyPlatformPermission(ctx, role.Permission("*.*")).Return(false, nil)
				authzService.EXPECT().GetUserAccessibleProjectIDs(ctx, 42).Return([]int{3, 7}, nil)
			},
			want: []int{3, 7},
		},
		{
			name: "platform permission failure is returned",
			ctx:  modelHandlerTestUserContext(),
			expect: func(ctx context.Context, authzService *authzmocks.MockIAuthzService) {
				authzService.EXPECT().VerifyPlatformPermission(ctx, role.Permission("*.*")).Return(false, errAuthz)
			},
			wantErr: errAuthz,
		},
		{
			name: "accessible project lookup failure is returned",
			ctx:  modelHandlerTestUserContext(),
			expect: func(ctx context.Context, authzService *authzmocks.MockIAuthzService) {
				authzService.EXPECT().VerifyPlatformPermission(ctx, role.Permission("*.*")).Return(false, nil)
				authzService.EXPECT().GetUserAccessibleProjectIDs(ctx, 42).Return(nil, errAuthz)
			},
			wantErr: errAuthz,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, _, authzService := modelHandlerTestSetup(t)
			tc.expect(tc.ctx, authzService)

			got, err := h.getAccessibleProjectIDs(tc.ctx)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestModelHandlerListModelTaskLabels(t *testing.T) {
	ctx := context.Background()

	t.Run("returns task labels", func(t *testing.T) {
		h, svc, _ := modelHandlerTestSetup(t)
		svc.EXPECT().ListModelTaskLabels(ctx).Return([]*model.Label{
			{ID: 1, Name: "text-generation", Category: "task", Scope: "model", CreatedAt: handlerTestCreatedAt, UpdatedAt: handlerTestUpdatedAt},
			{ID: 2, Name: "image-classification", Category: "task", Scope: "model", CreatedAt: handlerTestCreatedAt, UpdatedAt: handlerTestUpdatedAt},
		}, nil)

		got, err := h.ListModelTaskLabels(ctx, &v1alpha1.ListModelTaskLabelsRequest{})
		require.NoError(t, err)
		requireProtoEqual(t, &v1alpha1.ListModelTaskLabelsResponse{Items: []*v1alpha1.Label{
			{Id: 1, Name: "text-generation", Category: v1alpha1.Category_TASK, CreatedAt: handlerTestCreatedAtText, UpdatedAt: handlerTestUpdatedAtText},
			{Id: 2, Name: "image-classification", Category: v1alpha1.Category_TASK, CreatedAt: handlerTestCreatedAtText, UpdatedAt: handlerTestUpdatedAtText},
		}}, got)
	})

	t.Run("hides service errors", func(t *testing.T) {
		h, svc, _ := modelHandlerTestSetup(t)
		svc.EXPECT().ListModelTaskLabels(ctx).Return(nil, errors.New("database is locked"))

		got, err := h.ListModelTaskLabels(ctx, &v1alpha1.ListModelTaskLabelsRequest{})
		require.Nil(t, got)
		requireHandlerTestStatus(t, err, codes.Internal, "failed to list task labels")
	})
}

func TestModelHandlerListModelFrameLabels(t *testing.T) {
	ctx := context.Background()

	t.Run("returns library labels", func(t *testing.T) {
		h, svc, _ := modelHandlerTestSetup(t)
		svc.EXPECT().ListModelFrameLabels(ctx).Return([]*model.Label{
			{ID: 3, Name: "transformers", Category: "library", Scope: "model", CreatedAt: handlerTestCreatedAt, UpdatedAt: handlerTestUpdatedAt},
			{ID: 4, Name: "pytorch", Category: "library", Scope: "model", CreatedAt: handlerTestCreatedAt, UpdatedAt: handlerTestUpdatedAt},
		}, nil)

		got, err := h.ListModelFrameLabels(ctx, &v1alpha1.ListModelFrameLabelsRequest{})
		require.NoError(t, err)
		requireProtoEqual(t, &v1alpha1.ListModelFrameLabelsResponse{Items: []*v1alpha1.Label{
			{Id: 3, Name: "transformers", Category: v1alpha1.Category_LIBRARY, CreatedAt: handlerTestCreatedAtText, UpdatedAt: handlerTestUpdatedAtText},
			{Id: 4, Name: "pytorch", Category: v1alpha1.Category_LIBRARY, CreatedAt: handlerTestCreatedAtText, UpdatedAt: handlerTestUpdatedAtText},
		}}, got)
	})

	t.Run("hides service errors", func(t *testing.T) {
		h, svc, _ := modelHandlerTestSetup(t)
		svc.EXPECT().ListModelFrameLabels(ctx).Return(nil, errors.New("database is locked"))

		got, err := h.ListModelFrameLabels(ctx, &v1alpha1.ListModelFrameLabelsRequest{})
		require.Nil(t, got)
		requireHandlerTestStatus(t, err, codes.Internal, "failed to list library labels")
	})
}

func TestModelHandlerListModels(t *testing.T) {
	popular, notPopular := true, false
	llama := modelHandlerTestModel()
	qwen := &model.Model{ID: 9, Name: "qwen", ProjectName: "public", Size: 2048, ParameterCount: 500, CreatedAt: handlerTestCreatedAt, UpdatedAt: handlerTestUpdatedAt}
	qwenProto := &v1alpha1.Model{
		Id:             9,
		Name:           "qwen",
		CreatedAt:      handlerTestCreatedAtText,
		UpdatedAt:      handlerTestUpdatedAtText,
		CloneUrls:      &v1alpha1.CloneUrls{},
		Project:        "public",
		Size:           "2048",
		ParameterCount: "500",
	}
	request := &v1alpha1.ListModelsRequest{
		Labels:   []string{"text-generation", "transformers"},
		Search:   "lla",
		Sort:     "updated_at",
		Project:  modelHandlerTestProject,
		Page:     2,
		PageSize: 5,
		Popular:  true,
	}
	requestFilter := func(projectIDs []int) *model.Filter {
		return &model.Filter{
			Project:    modelHandlerTestProject,
			ProjectIDs: projectIDs,
			Label:      []string{"text-generation", "transformers"},
			Search:     "lla",
			Sort:       "updated_at",
			Page:       2,
			PageSize:   5,
			Popular:    &popular,
		}
	}

	t.Run("platform administrator lists without project restriction", func(t *testing.T) {
		h, svc, authzService := modelHandlerTestSetup(t)
		ctx := modelHandlerTestUserContext()
		authzService.EXPECT().VerifyPlatformPermission(ctx, role.Permission("*.*")).Return(true, nil)
		svc.EXPECT().ListModels(ctx, requestFilter(nil)).Return([]*model.Model{llama, qwen}, int64(12), nil)

		got, err := h.ListModels(ctx, request)
		require.NoError(t, err)
		requireProtoEqual(t, &v1alpha1.ListModelsResponse{
			Items:      []*v1alpha1.Model{modelHandlerTestModelProto(), qwenProto},
			Pagination: &v1alpha1.Pagination{Total: 12, Page: 2, PageSize: 5, Pages: 3},
		}, got)
	})

	t.Run("regular user only sees accessible projects", func(t *testing.T) {
		h, svc, authzService := modelHandlerTestSetup(t)
		ctx := modelHandlerTestUserContext()
		authzService.EXPECT().VerifyPlatformPermission(ctx, role.Permission("*.*")).Return(false, nil)
		authzService.EXPECT().GetUserAccessibleProjectIDs(ctx, 42).Return([]int{3, 7}, nil)
		svc.EXPECT().ListModels(ctx, requestFilter([]int{3, 7})).Return([]*model.Model{llama}, int64(6), nil)

		got, err := h.ListModels(ctx, request)
		require.NoError(t, err)
		requireProtoEqual(t, &v1alpha1.ListModelsResponse{
			Items:      []*v1alpha1.Model{modelHandlerTestModelProto()},
			Pagination: &v1alpha1.Pagination{Total: 6, Page: 2, PageSize: 5, Pages: 2},
		}, got)
	})

	t.Run("anonymous caller only sees public projects", func(t *testing.T) {
		h, svc, authzService := modelHandlerTestSetup(t)
		ctx := context.Background()
		authzService.EXPECT().GetUserAccessibleProjectIDs(ctx, 0).Return([]int{5}, nil)
		svc.EXPECT().ListModels(ctx, requestFilter([]int{5})).Return([]*model.Model{qwen}, int64(1), nil)

		got, err := h.ListModels(ctx, request)
		require.NoError(t, err)
		requireProtoEqual(t, &v1alpha1.ListModelsResponse{
			Items:      []*v1alpha1.Model{qwenProto},
			Pagination: &v1alpha1.Pagination{Total: 1, Page: 2, PageSize: 5, Pages: 1},
		}, got)
	})

	t.Run("applies default pagination", func(t *testing.T) {
		h, svc, authzService := modelHandlerTestSetup(t)
		ctx := modelHandlerTestUserContext()
		authzService.EXPECT().VerifyPlatformPermission(ctx, role.Permission("*.*")).Return(true, nil)
		svc.EXPECT().
			ListModels(ctx, &model.Filter{Page: 1, PageSize: 10, Popular: &notPopular}).
			Return([]*model.Model{llama}, int64(25), nil)

		got, err := h.ListModels(ctx, &v1alpha1.ListModelsRequest{})
		require.NoError(t, err)
		requireProtoEqual(t, &v1alpha1.ListModelsResponse{
			Items:      []*v1alpha1.Model{modelHandlerTestModelProto()},
			Pagination: &v1alpha1.Pagination{Total: 25, Page: 1, PageSize: 10, Pages: 3},
		}, got)
	})

	t.Run("returns an empty page without querying when no project is accessible", func(t *testing.T) {
		h, _, authzService := modelHandlerTestSetup(t)
		ctx := modelHandlerTestUserContext()
		authzService.EXPECT().VerifyPlatformPermission(ctx, role.Permission("*.*")).Return(false, nil)
		authzService.EXPECT().GetUserAccessibleProjectIDs(ctx, 42).Return([]int{}, nil)

		got, err := h.ListModels(ctx, request)
		require.NoError(t, err)
		require.NotNil(t, got.Items)
		require.Empty(t, got.Items)
		requireProtoEqual(t, &v1alpha1.Pagination{Total: 0, Page: 2, PageSize: 5, Pages: 0}, got.Pagination)
	})

	t.Run("hides accessible project lookup errors", func(t *testing.T) {
		h, _, authzService := modelHandlerTestSetup(t)
		ctx := modelHandlerTestUserContext()
		authzService.EXPECT().VerifyPlatformPermission(ctx, role.Permission("*.*")).Return(false, errors.New("database is locked"))

		got, err := h.ListModels(ctx, request)
		require.Nil(t, got)
		requireHandlerTestStatus(t, err, codes.Internal, "failed to get accessible projects")
	})

	t.Run("hides service errors", func(t *testing.T) {
		h, svc, authzService := modelHandlerTestSetup(t)
		ctx := modelHandlerTestUserContext()
		authzService.EXPECT().VerifyPlatformPermission(ctx, role.Permission("*.*")).Return(true, nil)
		svc.EXPECT().ListModels(ctx, requestFilter(nil)).Return(nil, int64(0), errors.New("database is locked"))

		got, err := h.ListModels(ctx, request)
		require.Nil(t, got)
		requireHandlerTestStatus(t, err, codes.Internal, "failed to list models")
	})
}

func TestModelHandlerRejectsInvalidRequests(t *testing.T) {
	cases := []struct {
		name  string
		field string
		call  func(context.Context, *ModelHandler) error
	}{
		{
			name:  "GetModel without project",
			field: "GetModelRequest.Project",
			call: func(ctx context.Context, h *ModelHandler) error {
				_, err := h.GetModel(ctx, &v1alpha1.GetModelRequest{Name: modelHandlerTestName})
				return err
			},
		},
		{
			name:  "GetModel with reserved name",
			field: "GetModelRequest.Name",
			call: func(ctx context.Context, h *ModelHandler) error {
				_, err := h.GetModel(ctx, &v1alpha1.GetModelRequest{Project: modelHandlerTestProject, Name: "models"})
				return err
			},
		},
		{
			name:  "CreateModel with path-like name",
			field: "CreateModelRequest.Name",
			call: func(ctx context.Context, h *ModelHandler) error {
				_, err := h.CreateModel(ctx, &v1alpha1.CreateModelRequest{Project: modelHandlerTestProject, Name: "../escape"})
				return err
			},
		},
		{
			name:  "CreateModel with too long name",
			field: "CreateModelRequest.Name",
			call: func(ctx context.Context, h *ModelHandler) error {
				_, err := h.CreateModel(ctx, &v1alpha1.CreateModelRequest{Project: modelHandlerTestProject, Name: strings.Repeat("a", 252)})
				return err
			},
		},
		{
			name:  "CreateModel with reserved name",
			field: "CreateModelRequest.Name",
			call: func(ctx context.Context, h *ModelHandler) error {
				_, err := h.CreateModel(ctx, &v1alpha1.CreateModelRequest{Project: modelHandlerTestProject, Name: "datasets"})
				return err
			},
