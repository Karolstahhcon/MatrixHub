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
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// createRepoAndCommit creates a repo and commits a file, returning the commit SHA.
func createRepoAndCommit(t *testing.T, endpoint, repoType, org, name string) string {
	t.Helper()
	// Create repo
	body := `{"type":"` + repoType + `","name":"` + name + `","organization":"` + org + `"}`
	resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	_ = resp.Body.Close()

	// Determine API prefix
	apiPrefix := "/api/models"
	switch repoType {
	case "dataset":
		apiPrefix = "/api/datasets"
	case "space":
		apiPrefix = "/api/spaces"
	}

	// Commit a file
	ndjson := "{\"key\":\"header\",\"value\":{\"summary\":\"Initial commit\"}}\n" +
		"{\"key\":\"file\",\"value\":{\"content\":\"# Test\\n\",\"path\":\"README.md\",\"encoding\":\"utf-8\"}}\n"

	resp, err = http.Post(endpoint+apiPrefix+"/"+org+"/"+name+"/commit/main", "application/x-ndjson", strings.NewReader(ndjson))
	if err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var commitResult commitResponse
	if err := json.NewDecoder(resp.Body).Decode(&commitResult); err != nil {
		t.Fatalf("Failed to decode commit response: %v", err)
	}
	return commitResult.CommitOid
}

func TestHuggingFaceDeleteRepo(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	// Create repo
	createBody := `{"type":"model","name":"delete-me","organization":"test-user"}`
	resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(createBody))
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	_ = resp.Body.Close()

	// Verify repo exists
	resp, err = http.Get(endpoint + "/api/models/test-user/delete-me")
	if err != nil {
		t.Fatalf("Failed to get repo info: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected repo to exist, got %d", resp.StatusCode)
	}

	// Delete repo
	deleteBody := `{"type":"model","name":"delete-me","organization":"test-user"}`
	req, _ := http.NewRequest(http.MethodDelete, endpoint+"/api/repos/delete", strings.NewReader(deleteBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete repo: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	// Verify repo no longer exists
	resp, err = http.Get(endpoint + "/api/models/test-user/delete-me")
	if err != nil {
		t.Fatalf("Failed to check repo: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404 after delete, got %d", resp.StatusCode)
	}
}

func TestHuggingFaceDeleteRepoNotFound(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	deleteBody := `{"type":"model","name":"nonexistent","organization":"test-user"}`
	req, _ := http.NewRequest(http.MethodDelete, endpoint+"/api/repos/delete", strings.NewReader(deleteBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete repo: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404, got %d", resp.StatusCode)
	}
}

func TestHuggingFaceDeleteDatasetRepo(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	// Create dataset repo
	createBody := `{"type":"dataset","name":"delete-dataset","organization":"test-user"}`
	resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(createBody))
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	_ = resp.Body.Close()

	// Delete dataset repo
	deleteBody := `{"type":"dataset","name":"delete-dataset","organization":"test-user"}`
	req, _ := http.NewRequest(http.MethodDelete, endpoint+"/api/repos/delete", strings.NewReader(deleteBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete repo: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	// Verify repo no longer exists
	resp, err = http.Get(endpoint + "/api/datasets/test-user/delete-dataset")
	if err != nil {
		t.Fatalf("Failed to check repo: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404 after delete, got %d", resp.StatusCode)
	}
}

func TestHuggingFaceMoveRepo(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	// Create and populate a repo
	createRepoAndCommit(t, endpoint, "model", "old-ns", "my-model")

	// Move repo
	moveBody := `{"fromRepo":"old-ns/my-model","toRepo":"new-ns/my-model","type":"model"}`
	resp, err := http.Post(endpoint+"/api/repos/move", "application/json", strings.NewReader(moveBody))
	if err != nil {
		t.Fatalf("Failed to move repo: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	// Verify old location no longer exists
	resp, err = http.Get(endpoint + "/api/models/old-ns/my-model")
	if err != nil {
		t.Fatalf("Failed to check old repo: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected 404 for old location, got %d", resp.StatusCode)
	}

	// Verify new location exists and has the file
	resp, err = http.Get(endpoint + "/new-ns/my-model/resolve/main/README.md")
	if err != nil {
		t.Fatalf("Failed to get file at new location: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 for new location, got %d", resp.StatusCode)
	}
	content, _ := io.ReadAll(resp.Body)
	if string(content) != "# Test\n" {
		t.Errorf("Unexpected content: %q", content)
	}
}

func TestHuggingFaceRepoSettings(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	// Create repo
	createBody := `{"type":"model","name":"settings-model","organization":"test-user"}`
	resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(createBody))
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	_ = resp.Body.Close()

	// Update settings
	settingsBody := `{"private":true}`
	req, _ := http.NewRequest(http.MethodPut, endpoint+"/api/models/test-user/settings-model/settings", strings.NewReader(settingsBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to update settings: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	// Update gated settings
	gatedBody := `{"gated":"auto"}`
	req, _ = http.NewRequest(http.MethodPut, endpoint+"/api/models/test-user/settings-model/settings", strings.NewReader(gatedBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to update gated setting: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestHuggingFaceCreateAndDeleteBranch(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	// Create and populate a repo
	createRepoAndCommit(t, endpoint, "model", "test-user", "branch-model")

	// Create a branch
	req, _ := http.NewRequest(http.MethodPost, endpoint+"/api/models/test-user/branch-model/branch/dev", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to create branch: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		t.Fatalf("Expected 200 for create branch, got %d: %s", resp.StatusCode, respBody)
	}
	_ = resp.Body.Close()

	// Verify the branch exists via refs endpoint
	resp, err = http.Get(endpoint + "/api/models/test-user/branch-model/refs")
	if err != nil {
		t.Fatalf("Failed to list refs: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 for refs, got %d", resp.StatusCode)
	}

	var refs gitRefs
	if err := json.NewDecoder(resp.Body).Decode(&refs); err != nil {
		t.Fatalf("Failed to decode refs: %v", err)
	}

	foundDev := false
	for _, b := range refs.Branches {
		if b.Name == "dev" {
			foundDev = true
			break
		}
	}
	if !foundDev {
		t.Errorf("Expected to find branch 'dev' in refs, got %+v", refs.Branches)
	}

	// Creating the same branch again should return conflict
	req, _ = http.NewRequest(http.MethodPost, endpoint+"/api/models/test-user/branch-model/branch/dev", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to create branch again: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("Expected 409 for duplicate branch, got %d", resp.StatusCode)
	}

	// Delete the branch
	req, _ = http.NewRequest(http.MethodDelete, endpoint+"/api/models/test-user/branch-model/branch/dev", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete branch: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 for delete branch, got %d", resp.StatusCode)
	}

	// Verify branch is gone
	resp, err = http.Get(endpoint + "/api/models/test-user/branch-model/refs")
	if err != nil {
		t.Fatalf("Failed to list refs after delete: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var refs2 gitRefs
	if err := json.NewDecoder(resp.Body).Decode(&refs2); err != nil {
		t.Fatalf("Failed to decode refs: %v", err)
	}

	for _, b := range refs2.Branches {
		if b.Name == "dev" {
			t.Errorf("Branch 'dev' should have been deleted, but still found in refs")
		}
	}
}

func TestHuggingFaceCreateBranchFromRevision(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	// Create and populate a repo
	commitSHA := createRepoAndCommit(t, endpoint, "model", "test-user", "branch-rev-model")

	// Create branch from specific revision
	body := `{"startingPoint":"` + commitSHA + `"}`
	req, _ := http.NewRequest(http.MethodPost, endpoint+"/api/models/test-user/branch-rev-model/branch/feature", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to create branch: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	// Verify branch points to the correct commit
	resp, err = http.Get(endpoint + "/api/models/test-user/branch-rev-model/refs")
	if err != nil {
		t.Fatalf("Failed to list refs: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var refs gitRefs
