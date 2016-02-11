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

