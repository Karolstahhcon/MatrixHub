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

package hf

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	backendhttp "github.com/matrixhub-ai/hfd/pkg/backend/http"
	backendlfs "github.com/matrixhub-ai/hfd/pkg/backend/lfs"
	"github.com/matrixhub-ai/hfd/pkg/repository"
	"github.com/matrixhub-ai/hfd/pkg/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	authzmocks "github.com/matrixhub-ai/matrixhub/internal/domain/authz/mocks"
	domaingit "github.com/matrixhub-ai/matrixhub/internal/domain/git"
	gitmocks "github.com/matrixhub-ai/matrixhub/internal/domain/git/mocks"
	modelmocks "github.com/matrixhub-ai/matrixhub/internal/domain/model/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
)

func createTestCommit(ctx context.Context, testStorage *storage.Storage, repoType, project, name, revision string, commit *domaingit.Commit, ops []domaingit.CommitOperation) (string, error) {
	prefix := ""
	if repoType == "datasets" || repoType == "spaces" {
		prefix = repoType + "/"
	}
	repo, err := repository.Open(testStorage.ResolvePath(prefix + project + "/" + name))
	if err != nil {
		return "", err
	}
	repoOps := make([]repository.CommitOperation, len(ops))
	for i, op := range ops {
		repoOps[i] = repository.CommitOperation{
			Type:    repository.CommitOperationType(op.Type),
			Path:    op.Path,
			Content: op.Content,
		}
	}
	return repo.CreateCommit(ctx, revision, commit.Message, commit.AuthorName, commit.AuthorEmail, repoOps, commit.ParentCommit)
}

func setupTestServer(t *testing.T) (*httptest.Server, string) {
	t.Helper()

	dataDir := t.TempDir()
	ctrl := gomock.NewController(t)
	authzService := authzmocks.NewMockIAuthzService(ctrl)
	authzService.EXPECT().
		VerifyProjectPermissionByName(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil).
		AnyTimes()

	testStorage := storage.NewStorage(storage.WithRootDir(dataDir))
	gitRepo := gitmocks.NewMockIGitRepo(ctrl)
	gitRepo.EXPECT().
		CreateCommit(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, repoType, project, name, revision string, commit *domaingit.Commit, ops []domaingit.CommitOperation) (string, error) {
			return createTestCommit(ctx, testStorage, repoType, project, name, revision, commit, ops)
		}).
		AnyTimes()
	modelService := modelmocks.NewMockIModelService(ctrl)
	modelService.EXPECT().
		CreateModelCommit(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, project, name, revision string, commit *domaingit.Commit, ops []domaingit.CommitOperation) (string, error) {
			return createTestCommit(ctx, testStorage, "models", project, name, revision, commit, ops)
		}).
		AnyTimes()

	// Set up handler chain (same order as main.go)
	var handler http.Handler

	handler = NewHandler(
		WithStorage(testStorage),
		WithServices(modelService, gitRepo, authzService),
	)

	handler = backendlfs.NewHandler(
		backendlfs.WithStorage(testStorage),
		backendlfs.WithNext(handler),
	)

	handler = backendhttp.NewHandler(
		backendhttp.WithStorage(testStorage),
		backendhttp.WithNext(handler),
	)

	server := httptest.NewServer(handler)
	t.Cleanup(func() { server.Close() })

	return server, dataDir
}

func TestRepositoryMutationAuthorization(t *testing.T) {
	authzErr := errors.New("authorization unavailable")
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		permission role.Permission
		allowed    bool
		authzErr   error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "create model denied",
			method:     http.MethodPost,
			path:       "/api/repos/create",
			body:       `{"type":"model","name":"model","organization":"project"}`,
			permission: role.ModelPush,
			wantStatus: http.StatusForbidden,
			wantBody:   "permission denied",
		},
		{
			name:       "create dataset authorization error",
			method:     http.MethodPost,
			path:       "/api/repos/create",
			body:       `{"type":"dataset","name":"dataset","organization":"project"}`,
			permission: role.DatasetPush,
			authzErr:   authzErr,
			wantStatus: http.StatusInternalServerError,
			wantBody:   authzErr.Error(),
		},
		{
			name:       "delete model denied",
			method:     http.MethodDelete,
			path:       "/api/repos/delete",
			body:       `{"type":"model","name":"model","organization":"project"}`,
			permission: role.ModelPush,
			wantStatus: http.StatusForbidden,
			wantBody:   "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			authzService := authzmocks.NewMockIAuthzService(ctrl)
			authzService.EXPECT().
				VerifyProjectPermissionByName(gomock.Any(), "project", tt.permission).
				Return(tt.allowed, tt.authzErr)

			handler := NewHandler(
				WithStorage(storage.NewStorage(storage.WithRootDir(t.TempDir()))),
				WithServices(nil, nil, authzService),
			)
			request := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			require.Equal(t, tt.wantStatus, response.Code)
			require.Contains(t, response.Body.String(), tt.wantBody)
		})
	}
}

func TestHuggingFaceCreateRepo(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	body := `{"type":"model","name":"test-model","organization":"test-user"}`
	resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected 200, got %d: %s", resp.StatusCode, respBody)
	}

	var result createRepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if result.URL == "" {
		t.Error("Expected url in response")
	}

	// Creating the same repo again should succeed (exist_ok behavior)
	resp2, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create repo again: %v", err)
	}
	_ = resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp2.Body)
		t.Fatalf("Expected 200 for existing repo, got %d: %s", resp2.StatusCode, respBody)
	}
}

func TestHuggingFacePreupload(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	// Create repo first
	createBody := `{"type":"model","name":"test-model","organization":"test-user"}`
	resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(createBody))
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	_ = resp.Body.Close()

	// Test preupload
	preuploadBody := `{"files":[{"path":"README.md","size":20,"sample":""},{"path":"large.bin","size":20000000,"sample":""}]}`
	resp, err = http.Post(endpoint+"/api/models/test-user/test-model/preupload/main", "application/json", strings.NewReader(preuploadBody))
	if err != nil {
		t.Fatalf("Failed to preupload: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected 200, got %d: %s", resp.StatusCode, respBody)
	}

	var result preuploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
