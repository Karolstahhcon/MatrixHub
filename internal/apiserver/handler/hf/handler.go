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
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matrixhub-ai/hfd/pkg/lfs"
	"github.com/matrixhub-ai/hfd/pkg/mirror"
	"github.com/matrixhub-ai/hfd/pkg/permission"
	"github.com/matrixhub-ai/hfd/pkg/receive"
	"github.com/matrixhub-ai/hfd/pkg/repository"
	"github.com/matrixhub-ai/hfd/pkg/storage"

	"github.com/matrixhub-ai/matrixhub/internal/domain/authz"
	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
)

// Handler handles HTTP requests for HuggingFace-compatible API endpoints, including repository management and git operations.
type Handler struct {
	storage             *storage.Storage
	root                *mux.Router
	middlewares         []mux.MiddlewareFunc
	next                http.Handler
	lfsStorage          lfs.Storage
	permissionHookFunc  permission.PermissionHookFunc
	preReceiveHookFunc  receive.PreReceiveHookFunc
	postReceiveHookFunc receive.PostReceiveHookFunc
	mirror              *mirror.Mirror
	modelService        model.IModelService
	gitRepo             git.IGitRepo
	authzService        authz.IAuthzService
}

// Option defines a functional option for configuring the Handler.
type Option func(*Handler)

// WithServices sets the services for the router.
func WithServices(model model.IModelService, git git.IGitRepo, authz authz.IAuthzService) Option {
	return func(h *Handler) {
		h.modelService = model
		h.gitRepo = git
		h.authzService = authz
	}
}

// WithMiddlewares sets the middlewares for the router.
func WithMiddlewares(middlewares ...mux.MiddlewareFunc) Option {
	return func(h *Handler) {
		h.middlewares = middlewares
	}
}

// WithStorage sets the storage backend for the handler. This is required.
func WithStorage(storage *storage.Storage) Option {
	return func(h *Handler) {
		h.storage = storage
	}
}

// WithNext sets the next http.Handler to call if the request is not handled by this handler.
func WithNext(next http.Handler) Option {
	return func(h *Handler) {
		h.next = next
	}
}

// WithPermissionHookFunc sets the permission hook for verifying operations.
func WithPermissionHookFunc(fn permission.PermissionHookFunc) Option {
	return func(h *Handler) {
		h.permissionHookFunc = fn
	}
}

// WithPreReceiveHookFunc sets the pre-receive hook called before ref changes are applied.
// If the hook returns an error, the operation is rejected.
func WithPreReceiveHookFunc(fn receive.PreReceiveHookFunc) Option {
