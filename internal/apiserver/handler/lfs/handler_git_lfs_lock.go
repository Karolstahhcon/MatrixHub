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

	limit := 0
	if limitStr := r.FormValue("limit"); limitStr != "" {
		strtLimit, err := strconv.Atoi(limitStr)
		if err != nil || strtLimit < 0 {
			responseJSON(w, &lfs.LockList{Message: "invalid limit parameter"}, http.StatusBadRequest)
			return
		}
		limit = strtLimit
	}

	locks, nextCursor, err := h.locksStorage.Filtered(repoName,
		r.FormValue("path"),
		r.FormValue("cursor"),
		limit,
	)

	if err != nil {
		ll.Message = err.Error()
	} else {
		ll.Locks = locks
		ll.NextCursor = nextCursor
	}

	responseJSON(w, ll, http.StatusOK)
}

func (h *Handler) handleLocksVerify(w http.ResponseWriter, r *http.Request) {
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

	user := getUserFromRequest(r)

	dec := json.NewDecoder(r.Body)

	w.Header().Set("Content-Type", metaMediaType)

	reqBody := &lfs.VerifiableLockRequest{}
	if err := dec.Decode(reqBody); err != nil {
		responseJSON(w, &lfs.VerifiableLockList{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	// Limit is optional
	limit := reqBody.Limit
	if limit == 0 {
		limit = 100
	}

	ll := &lfs.VerifiableLockList{}
	locks, nextCursor, err := h.locksStorage.Filtered(repoName,
		"",
		reqBody.Cursor,
		limit,
	)
	if err != nil {
		ll.Message = err.Error()
	} else {
		ll.NextCursor = nextCursor

		for _, l := range locks {
			if l.Owner.Name == user {
				ll.Ours = append(ll.Ours, l)
			} else {
				ll.Theirs = append(ll.Theirs, l)
			}
		}
	}

	responseJSON(w, ll, http.StatusOK)
}

func (h *Handler) handleCreateLock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo"]

	if h.permissionHookFunc != nil {
		op := permission.OperationUpdateRepo
		if ok, err := h.permissionHookFunc(r.Context(), op, repoName, permission.Context{}); err != nil {
			responseJSON(w, &lfs.VerifiableLockList{Message: err.Error()}, http.StatusInternalServerError)
			return
		} else if !ok {
			responseJSON(w, &lfs.VerifiableLockList{Message: "permission denied"}, http.StatusForbidden)
			return
		}
	}

