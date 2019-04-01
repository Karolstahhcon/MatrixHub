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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/matrixhub-ai/hfd/pkg/hf"
	"github.com/matrixhub-ai/hfd/pkg/repository"
)

// handleList is the unified handler for listing repositories of different types (models, datasets, spaces).
func (h *Handler) handleList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoType := vars["repoType"]
	h.handleListRepos(w, r, repoType)
}

// repoListFilter holds the parsed query parameters for listing repositories.
type repoListFilter struct {
	search     string
	author     string
	filterTags []string
	sortField  string
	limit      int
	offset     int
}

// parseRepoListFilter extracts and parses the list query parameters from the request.
func parseRepoListFilter(r *http.Request) repoListFilter {
	query := r.URL.Query()
	f := repoListFilter{
		search:     query.Get("search"),
		author:     query.Get("author"),
		filterTags: query["filter"],
		sortField:  query.Get("sort"),
	}
	if v, err := strconv.Atoi(query.Get("limit")); err == nil && v > 0 {
		f.limit = v
	}
	if cursor := query.Get("cursor"); cursor != "" {
		f.offset = decodeCursorOffset(cursor)
	}
	return f
}

// handleListRepos is the unified handler for listing models, datasets, or spaces.
func (h *Handler) handleListRepos(w http.ResponseWriter, r *http.Request, repoType string) {
	f := parseRepoListFilter(r)

	reposDir := h.storage.RepositoriesDir()
	isModel := repoType == "models"

	baseDir := reposDir
	if !isModel {
		baseDir = filepath.Join(reposDir, repoType)
	}

	entries := discoverRepos(baseDir, isModel, f.author)
	items := buildRepoListItems(entries, isModel, f.search, f.filterTags)

	// Sort results
	sortRepoItems(items, f.sortField)

	// Apply cursor offset
	if f.offset > 0 && f.offset < len(items) {
		items = items[f.offset:]
	} else if f.offset >= len(items) {
		items = nil
	}

	// Apply limit and set Link header if there are more results
	if f.limit > 0 && len(items) > f.limit {
		items = items[:f.limit]

		nextCursor := encodeCursorOffset(f.offset + f.limit)
		nextURL := buildNextURL(r, nextCursor)
		w.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"next\"", nextURL))
	}

	if items == nil {
		items = []repoListItem{}
	}

	responseJSON(w, items, http.StatusOK)
}

// repoEntry represents a discovered repository on disk.
type repoEntry struct {
	fullName string // "namespace/repo"
	repoPath string // absolute path to the .git directory
}

// isSafePathComponent reports whether component is a single, non-traversing path component.
func isSafePathComponent(component string) bool {
	if component == "" || component == "." || component == ".." {
		return false
	}
	if filepath.IsAbs(component) {
		return false
	}
	return !strings.ContainsAny(component, `/\\`)
}

// resolveContainedPath resolves a safe path component below root and rejects
// paths that escape root, including through symbolic links.
func resolveContainedPath(root, component string) (string, error) {
	if !isSafePathComponent(component) {
		return "", fmt.Errorf("invalid path component %q", component)
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	resolvedRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(filepath.Join(root, component))
	if err != nil {
		return "", err
	}
	resolvedPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(resolvedRoot, resolvedPath)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path %q escapes root %q", component, root)
	}
	return resolvedPath, nil
}
