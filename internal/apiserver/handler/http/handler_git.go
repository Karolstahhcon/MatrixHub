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

package backend

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/matrixhub-ai/hfd/pkg/authenticate"
	"github.com/matrixhub-ai/hfd/pkg/permission"
	"github.com/matrixhub-ai/hfd/pkg/receive"
	"github.com/matrixhub-ai/hfd/pkg/repository"

	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
)

type repoInformation struct {
	RepoType string

	FullName  string
	Namespace string
	Name      string
}

var repoTypePrefixes = []struct {
	prefix   string
	repoType string
}{
	{"datasets/", "datasets"},
	{"spaces/", "spaces"},
}

// getRepoInformation returns the repository information extracted from the request, including repo type, storage path, namespace, and name.
func getRepoInformation(r *http.Request) repoInformation {
	vars := mux.Vars(r)
	name := vars["repo"]

	repoType := "models"
	namespacedName := name
	for _, p := range repoTypePrefixes {
		if strings.HasPrefix(name, p.prefix) {
			repoType = p.repoType
			namespacedName = strings.TrimPrefix(name, p.prefix)
			break
		}
	}

	namespace, repoName, ok := strings.Cut(namespacedName, "/")
	if !ok {
		return repoInformation{}
	}

	return repoInformation{
		RepoType:  repoType,
		FullName:  name,
		Namespace: namespace,
		Name:      repoName,
	}
}

// gitProtocolEnv returns a GIT_PROTOCOL environment variable derived from the
// request's Git-Protocol header if the value is present and valid, or nil otherwise.
func gitProtocolEnv(r *http.Request) []string {
	value := r.Header.Get("Git-Protocol")
	if value == "" || !repository.IsValidGitProtocol(value) {
		return nil
	}
	return []string{"GIT_PROTOCOL=" + value}
}

// handleInfoRefs handles the /info/refs endpoint for git service discovery.
func (h *Handler) handleInfoRefs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo"]

	service := r.URL.Query().Get("service")
	if service == "" {
		responseText(w, "service parameter is required", http.StatusBadRequest)
		return
	}

	if service != repository.GitUploadPack && service != repository.GitReceivePack {
		responseText(w, "unsupported service", http.StatusForbidden)
		return
	}

	if h.permissionHookFunc != nil {
		op := permission.OperationReadRepo
		if service == repository.GitReceivePack {
			op = permission.OperationUpdateRepo
		}
		if ok, err := h.permissionHookFunc(r.Context(), op, repoName, permission.Context{}); err != nil {
			responseText(w, err.Error(), http.StatusInternalServerError)
			return
		} else if !ok {
			userinfo, ok := authenticate.GetUserInfo(r.Context())
			if ok && userinfo.User == authenticate.Anonymous {
				responseText(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			responseText(w, "permission denied", http.StatusForbidden)
			return
		}
	}

	repoPath := h.storage.ResolvePath(repoName)
	if repoPath == "" {
		responseText(w, fmt.Sprintf("repository %q not found", repoName), http.StatusNotFound)
		return
	}

	repoInfo := getRepoInformation(r)
	repo, err := h.openRepo(r.Context(), repoPath, repoInfo, service)
	if err != nil {
		if errors.Is(err, repository.ErrRepositoryNotExists) {
			responseText(w, fmt.Sprintf("repository %q not found", repoName), http.StatusNotFound)
			return
		}
		responseText(w, fmt.Sprintf("Failed to open repository %q: %v", repoName, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
	w.Header().Set("Cache-Control", "no-cache")

	err = repo.Stateless(r.Context(), w, nil, service, true, gitProtocolEnv(r)...)
	if err != nil {
		responseText(w, fmt.Sprintf("Failed to get info refs for %q: %v", repoName, err), http.StatusInternalServerError)
		return
	}
}

// handleUploadPack handles the git-upload-pack endpoint (fetch/clone).
func (h *Handler) handleUploadPack(w http.ResponseWriter, r *http.Request) {
	h.handleService(w, r, repository.GitUploadPack)
}

// handleReceivePack handles the git-receive-pack endpoint (push).
func (h *Handler) handleReceivePack(w http.ResponseWriter, r *http.Request) {
	h.handleService(w, r, repository.GitReceivePack)
}

// handleService handles a git service request.
func (h *Handler) handleService(w http.ResponseWriter, r *http.Request, service string) {
	vars := mux.Vars(r)
