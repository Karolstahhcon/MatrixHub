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

package lfs_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matrixhub-ai/hfd/pkg/authenticate"
	hfdlfs "github.com/matrixhub-ai/hfd/pkg/lfs"
	"github.com/matrixhub-ai/hfd/pkg/permission"
	"github.com/matrixhub-ai/hfd/pkg/storage"
	"github.com/stretchr/testify/require"

	backend "github.com/matrixhub-ai/matrixhub/internal/apiserver/handler/lfs"
)

const lfsMetaMediaType = "application/vnd.git-lfs+json"

func newLFSHandler(t *testing.T, opts ...backend.Option) *backend.Handler {
	t.Helper()
	dataDir := t.TempDir()
	base := []backend.Option{
		backend.WithStorage(storage.NewStorage(storage.WithRootDir(dataDir))),
		backend.WithLFSStorage(hfdlfs.NewLocal(dataDir)),
	}
	return backend.NewHandler(append(base, opts...)...)
}

func lfsRequest(t *testing.T, handler http.Handler, method, target, body, user string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Accept", lfsMetaMediaType)
	if user != "" {
		req = req.WithContext(authenticate.WithContext(req.Context(), authenticate.UserInfo{User: user}))
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
}

func objectOID(content []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(content))
}

func TestLFSObjectLifecycle(t *testing.T) {
	handler := newLFSHandler(t)
	content := []byte("matrixhub-lfs-content")
	oid := objectOID(content)

	batchBody := fmt.Sprintf(`{"operation":"upload","objects":[{"oid":%q,"size":%d}]}`, oid, len(content))
	response := lfsRequest(t, handler, http.MethodPost, "http://example.com/org/repo.git/info/lfs/objects/batch", batchBody, "")
	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, lfsMetaMediaType, response.Header().Get("Content-Type"))
	var batch struct {
		Objects []struct {
			OID     string `json:"oid"`
			Actions map[string]struct {
				Href string `json:"href"`
			} `json:"actions"`
		} `json:"objects"`
	}
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &batch))
	require.Len(t, batch.Objects, 1)
	require.Equal(t, oid, batch.Objects[0].OID)
	require.Equal(t, "http://example.com/objects/"+oid, batch.Objects[0].Actions["upload"].Href)
	require.Equal(t, "http://example.com/objects/"+oid+"/verify", batch.Objects[0].Actions["verify"].Href)

	put := httptest.NewRequest(http.MethodPut, "/objects/"+oid, bytes.NewReader(content))
	put.ContentLength = int64(len(content))
