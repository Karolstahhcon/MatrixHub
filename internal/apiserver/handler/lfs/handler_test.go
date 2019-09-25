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
