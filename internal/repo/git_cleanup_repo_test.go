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

package repo

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matrixhub-ai/hfd/pkg/repository"
	hfdstorage "github.com/matrixhub-ai/hfd/pkg/storage"
)

func TestFindOrphanedLFSShardedObjects(t *testing.T) {
	ctx := context.Background()
	store := hfdstorage.NewStorage(hfdstorage.WithRootDir(t.TempDir()))
	referencedOID := strings.Repeat("a", 64)
	orphanedOID := strings.Repeat("b", 64)
	for _, oid := range []string{referencedOID, orphanedOID} {
		objectPath := filepath.Join(store.LFSDir(), oid[:2], oid[2:4], oid[4:])
		if err := os.MkdirAll(filepath.Dir(objectPath), 0750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(objectPath, []byte("object"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(store.LFSDir(), "bb", "bb", "lfsd_tmp_x"), []byte("pending"), 0600); err != nil {
		t.Fatal(err)
	}

	repoPath := store.ResolvePath("test-project/test-model")
	if err := os.MkdirAll(filepath.Dir(repoPath), 0750); err != nil {
		t.Fatal(err)
	}
	repo, err := repository.Init(ctx, repoPath, "main")
