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

package backend_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matrixhub-ai/hfd/pkg/authenticate"
	"github.com/matrixhub-ai/hfd/pkg/permission"
	"github.com/matrixhub-ai/hfd/pkg/storage"
	"github.com/stretchr/testify/require"

	backend "github.com/matrixhub-ai/matrixhub/internal/apiserver/handler/http"
)

func TestHTTPHandlerOptionsAndFallback(t *testing.T) {
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Middleware", "true")
			next.ServeHTTP(w, r)
		})
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	handler := backend.NewHandler(
		backend.WithMiddlewares(middleware),
		backend.WithStorage(storage.NewStorage(storage.WithRootDir(t.TempDir()))),
		backend.WithNext(next),
		backend.WithPermissionHookFunc(nil),
		backend.WithPreReceiveHookFunc(nil),
		backend.WithPostReceiveHookFunc(nil),
		backend.WithMirror(nil),
		backend.WithServices(nil, nil, nil),
	)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/not-git", nil))
	require.Equal(t, http.StatusTeapot, response.Code)

	response = httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/org/repo.git/info/refs", nil))
	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Equal(t, "true", response.Header().Get("X-Test-Middleware"))
}

func TestHTTPHandlerRejectsInvalidDiscoveryRequests(t *testing.T) {
	handler := backend.NewHandler(
		backend.WithStorage(storage.NewStorage(storage.WithRootDir(t.TempDir()))),
	)
	tests := []struct {
		name       string
		method     string
		target     string
		wantStatus int
		wantBody   string
	}{
