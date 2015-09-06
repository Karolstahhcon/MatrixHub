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
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/matrixhub-ai/hfd/pkg/authenticate"
	"github.com/matrixhub-ai/hfd/pkg/lfs"
	"github.com/matrixhub-ai/hfd/pkg/permission"
)

func (h *Handler) handleGetLock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo"]

	if h.permissionHookFunc != nil {
		op := permission.OperationReadRepo
		if ok, err := h.permissionHookFunc(r.Context(), op, repoName, permission.Context{}); err != nil {
			responseJSON(w, &lfs.VerifiableLockList{Message: err.Error()}, http.StatusInternalServerError)
			return
		} else if !ok {
			responseJSON(w, &lfs.VerifiableLockList{Message: "permission denied"}, http.StatusForbidden)
			return
		}
	}

	ll := &lfs.LockList{}

	w.Header().Set("Content-Type", metaMediaType)
