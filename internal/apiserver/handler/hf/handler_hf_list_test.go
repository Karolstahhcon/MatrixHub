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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsSafePathComponent(t *testing.T) {
	tests := []struct {
		component string
		safe      bool
	}{
		{component: "alice", safe: true},
		{component: "user.name", safe: true},
		{component: "", safe: false},
		{component: ".", safe: false},
		{component: "..", safe: false},
		{component: "../alice", safe: false},
		{component: "alice/../bob", safe: false},
		{component: "/tmp", safe: false},
		{component: `\\tmp`, safe: false},
	}

	for _, tt := range tests {
		if got := isSafePathComponent(tt.component); got != tt.safe {
			t.Errorf("isSafePathComponent(%q) = %t, want %t", tt.component, got, tt.safe)
		}
	}
}

func TestResolveContainedPath(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "namespace"), 0o755); err != nil {
		t.Fatal(err)
	}

	resolved, err := resolveContainedPath(root, "namespace")
	if err != nil {
		t.Fatalf("resolveContainedPath() error = %v", err)
	}
	want, err := filepath.EvalSymlinks(filepath.Join(root, "namespace"))
	if err != nil {
		t.Fatal(err)
	}
	if resolved != want {
		t.Fatalf("resolveContainedPath() = %q, want %q", resolved, want)
	}

	for _, relativePath := range []string{"../outside", "namespace/../../outside"} {
		if _, err := resolveContainedPath(root, relativePath); err == nil {
			t.Errorf("resolveContainedPath(%q) succeeded, want containment error", relativePath)
		}
	}
}

func TestResolveContainedPathRejectsExternalSymlink(t *testing.T) {
	root := t.TempDir()
	external := t.TempDir()
	link := filepath.Join(root, "namespace")
	if err := os.Symlink(external, link); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	if _, err := resolveContainedPath(root, "namespace"); err == nil {
		t.Fatal("resolveContainedPath() followed external symlink")
	}
}

func TestDiscoverReposInNamespaceRejectsUnsafeNamespace(t *testing.T) {
	root := t.TempDir()

	for _, namespace := range []string{"", ".", "..", "../outside", "namespace/../outside", "/tmp", `\\tmp`} {
		if entries := discoverReposInNamespace(root, namespace); len(entries) != 0 {
			t.Errorf("discoverReposInNamespace(%q) returned %d entries, want none", namespace, len(entries))
		}
	}
}

func TestDiscoverReposInNamespace(t *testing.T) {
	root := t.TempDir()
	repoPath := filepath.Join(root, "namespace", "model.git")
	if err := os.MkdirAll(repoPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repoPath, "HEAD"), []byte("ref: refs/heads/main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	entries := discoverReposInNamespace(root, "namespace")
	if len(entries) != 1 {
		t.Fatalf("discoverReposInNamespace() returned %d entries, want 1", len(entries))
	}
	if entries[0].fullName != "namespace/model" {
		t.Errorf("fullName = %q, want %q", entries[0].fullName, "namespace/model")
	}
	resolvedNamespace, err := filepath.EvalSymlinks(filepath.Join(root, "namespace"))
	if err != nil {
		t.Fatal(err)
	}
	wantRepoPath := filepath.Join(resolvedNamespace, "model.git")
	if entries[0].repoPath != wantRepoPath {
		t.Errorf("repoPath = %q, want %q", entries[0].repoPath, wantRepoPath)
	}
}

func TestHandleListModelsEmpty(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	resp, err := http.Get(endpoint + "/api/models")
	if err != nil {
		t.Fatalf("Failed to list models: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected 200, got %d: %s", resp.StatusCode, respBody)
	}

	var items []repoListItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("Expected 0 models, got %d", len(items))
	}
}

func TestHandleListModels(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	// Create two model repos
	for _, body := range []string{
		`{"type":"model","name":"alpha-model","organization":"user-a"}`,
		`{"type":"model","name":"beta-model","organization":"user-b"}`,
	} {
		resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create repo: %v", err)
		}
		_ = resp.Body.Close()
	}

	// Also create a dataset to ensure it's NOT listed
	resp, err := http.Post(endpoint+"/api/repos/create", "application/json",
		strings.NewReader(`{"type":"dataset","name":"my-dataset","organization":"user-a"}`))
	if err != nil {
		t.Fatalf("Failed to create dataset: %v", err)
	}
	_ = resp.Body.Close()

	// List all models
	resp, err = http.Get(endpoint + "/api/models")
	if err != nil {
		t.Fatalf("Failed to list models: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected 200, got %d: %s", resp.StatusCode, respBody)
	}

	var items []repoListItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("Expected 2 models, got %d", len(items))
	}

	// Verify modelId is set for models
	for _, item := range items {
		if item.ModelID == "" {
			t.Errorf("Expected modelId to be set for model %s", item.RepoID)
		}
	}
}

func TestHandleListModelsFilterByAuthor(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	// Create repos under different authors
	for _, body := range []string{
		`{"type":"model","name":"model-1","organization":"alice"}`,
		`{"type":"model","name":"model-2","organization":"bob"}`,
	} {
		resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create repo: %v", err)
		}
		_ = resp.Body.Close()
	}

	resp, err := http.Get(endpoint + "/api/models?author=alice")
	if err != nil {
		t.Fatalf("Failed to list models: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var items []repoListItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("Expected 1 model, got %d", len(items))
	}
	if items[0].RepoID != "alice/model-1" {
		t.Errorf("Expected alice/model-1, got %s", items[0].RepoID)
	}
}

func TestHandleListModelsSearch(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	for _, body := range []string{
		`{"type":"model","name":"llama-7b","organization":"meta"}`,
		`{"type":"model","name":"bert-base","organization":"google"}`,
	} {
		resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create repo: %v", err)
		}
		_ = resp.Body.Close()
	}

	resp, err := http.Get(endpoint + "/api/models?search=llama")
	if err != nil {
		t.Fatalf("Failed to list models: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var items []repoListItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("Expected 1 model, got %d", len(items))
	}
	if items[0].RepoID != "meta/llama-7b" {
		t.Errorf("Expected meta/llama-7b, got %s", items[0].RepoID)
	}
}

func TestHandleListModelsLimit(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	for _, body := range []string{
		`{"type":"model","name":"model-a","organization":"org"}`,
		`{"type":"model","name":"model-b","organization":"org"}`,
		`{"type":"model","name":"model-c","organization":"org"}`,
	} {
		resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create repo: %v", err)
		}
		_ = resp.Body.Close()
	}

	resp, err := http.Get(endpoint + "/api/models?limit=2")
	if err != nil {
		t.Fatalf("Failed to list models: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var items []repoListItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("Expected 2 models, got %d", len(items))
	}
}

func TestHandleListModelsResponseFormat(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	// Create a model repo and commit a file with metadata
	createBody := `{"type":"model","name":"test-model","organization":"test-org"}`
	resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(createBody))
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	_ = resp.Body.Close()

	// Commit a README with pipeline_tag and library_name
	ndjson := `{"key":"header","value":{"summary":"Initial commit"}}` + "\n" +
		`{"key":"file","value":{"content":"---\npipeline_tag: text-generation\nlibrary_name: transformers\n---\n# Test\n","path":"README.md","encoding":"utf-8"}}` + "\n"
	resp, err = http.Post(endpoint+"/api/models/test-org/test-model/commit/main", "application/x-ndjson", strings.NewReader(ndjson))
	if err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
	_ = resp.Body.Close()

	// List models and verify response format
	resp, err = http.Get(endpoint + "/api/models")
	if err != nil {
		t.Fatalf("Failed to list models: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var items []repoListItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("Expected 1 model, got %d", len(items))
	}

	item := items[0]
	if item.RepoID != "test-org/test-model" {
		t.Errorf("Expected id=test-org/test-model, got %s", item.RepoID)
	}
	if item.ModelID != "test-org/test-model" {
		t.Errorf("Expected modelId=test-org/test-model, got %s", item.ModelID)
	}
	if item.PipelineTag != "text-generation" {
		t.Errorf("Expected pipeline_tag=text-generation, got %s", item.PipelineTag)
	}
	if item.LibraryName != "transformers" {
		t.Errorf("Expected library_name=transformers, got %s", item.LibraryName)
	}
}

func TestHandleListDatasetsEmpty(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	resp, err := http.Get(endpoint + "/api/datasets")
	if err != nil {
		t.Fatalf("Failed to list datasets: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected 200, got %d: %s", resp.StatusCode, respBody)
	}

	var items []repoListItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("Expected 0 datasets, got %d", len(items))
	}
}

func TestHandleListDatasets(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	// Create datasets
	for _, body := range []string{
		`{"type":"dataset","name":"dataset-a","organization":"org-a"}`,
		`{"type":"dataset","name":"dataset-b","organization":"org-b"}`,
	} {
		resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create repo: %v", err)
		}
		_ = resp.Body.Close()
	}

	// Also create a model to ensure it's NOT listed in datasets
	resp, err := http.Post(endpoint+"/api/repos/create", "application/json",
		strings.NewReader(`{"type":"model","name":"my-model","organization":"org-a"}`))
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	_ = resp.Body.Close()

	// List all datasets
	resp, err = http.Get(endpoint + "/api/datasets")
	if err != nil {
		t.Fatalf("Failed to list datasets: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var items []repoListItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("Expected 2 datasets, got %d", len(items))
	}

	// Verify modelId is NOT set for datasets
	for _, item := range items {
		if item.ModelID != "" {
			t.Errorf("Expected modelId to be empty for dataset %s, got %s", item.RepoID, item.ModelID)
		}
	}
}

func TestHandleListDatasetsFilterByAuthor(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	for _, body := range []string{
		`{"type":"dataset","name":"ds-1","organization":"alice"}`,
		`{"type":"dataset","name":"ds-2","organization":"bob"}`,
	} {
		resp, err := http.Post(endpoint+"/api/repos/create", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create repo: %v", err)
		}
		_ = resp.Body.Close()
	}

	resp, err := http.Get(endpoint + "/api/datasets?author=alice")
	if err != nil {
		t.Fatalf("Failed to list datasets: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var items []repoListItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("Expected 1 dataset, got %d", len(items))
	}
	if items[0].RepoID != "alice/ds-1" {
		t.Errorf("Expected alice/ds-1, got %s", items[0].RepoID)
	}
}

func TestHandleListSpacesEmpty(t *testing.T) {
	server, _ := setupTestServer(t)
	endpoint := server.URL

	resp, err := http.Get(endpoint + "/api/spaces")
	if err != nil {
		t.Fatalf("Failed to list spaces: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected 200, got %d: %s", resp.StatusCode, respBody)
	}

	var items []repoListItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
