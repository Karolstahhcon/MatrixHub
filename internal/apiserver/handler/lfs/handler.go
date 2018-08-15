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

package lfs

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matrixhub-ai/hfd/pkg/authenticate"
	"github.com/matrixhub-ai/hfd/pkg/lfs"
	"github.com/matrixhub-ai/hfd/pkg/mirror"
	"github.com/matrixhub-ai/hfd/pkg/permission"
	"github.com/matrixhub-ai/hfd/pkg/storage"
)

// Handler handles HTTP requests for Git LFS API endpoints, including batch operations and object content management.
type Handler struct {
	storage            *storage.Storage
	root               *mux.Router
	next               http.Handler
	middlewares        []mux.MiddlewareFunc
	lfsStorage         lfs.Storage
	locksStorage       *lfs.LockStorage
	permissionHookFunc permission.PermissionHookFunc
	tokenSignValidator authenticate.TokenSignValidator
	mirror             *mirror.Mirror
}

// Option defines a functional option for configuring the Handler.
type Option func(*Handler)

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

// WithPermissionHookFunc sets the authentication hook for verifying operations.
func WithPermissionHookFunc(fn permission.PermissionHookFunc) Option {
	return func(h *Handler) {
		h.permissionHookFunc = fn
	}
}

// WithTokenSignValidator sets the token signer for signing LFS batch response action headers.
func WithTokenSignValidator(signer authenticate.TokenSignValidator) Option {
	return func(h *Handler) {
		h.tokenSignValidator = signer
	}
}

// WithLFSStorage configures the LFS storage backend.
func WithLFSStorage(storage lfs.Storage) Option {
	return func(h *Handler) {
		h.lfsStorage = storage
	}
}

// WithMirror sets the mirror to use for repository synchronization. If not provided,
// a mirror will be created when mirrorSourceFunc is set.
func WithMirror(m *mirror.Mirror) Option {
	return func(h *Handler) {
		h.mirror = m
	}
}

// NewHandler creates a new Handler with the given repository directory.
func NewHandler(opts ...Option) *Handler {
	h := &Handler{
		root:         mux.NewRouter(),
		locksStorage: lfs.NewLock(),
	}

	for _, opt := range opts {
		opt(h)
	}

	h.register()
	return h
}

// ServeHTTP implements the http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.root.ServeHTTP(w, r)
}

func (h *Handler) register() {
	// Git LFS protocol endpoints
	h.registryLFS(h.root)
	h.registryLFSLock(h.root)

	h.root.Use(h.middlewares...)
