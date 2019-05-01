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
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	backendhttp "github.com/matrixhub-ai/hfd/pkg/backend/http"
	hfdlfs "github.com/matrixhub-ai/hfd/pkg/lfs"
	"github.com/matrixhub-ai/hfd/pkg/mirror"
	"github.com/matrixhub-ai/hfd/pkg/repository"
	hfdstorage "github.com/matrixhub-ai/hfd/pkg/storage"
	"github.com/stretchr/testify/require"

	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
)

const (
	repoTestProject = "test-project"
	repoTestName    = "test-model"
	// repoTestMissingOID references an LFS object that no test stores locally.
	repoTestMissingOID = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
)

func TestIsCommitSHA(t *testing.T) {
	for _, test := range []struct {
		name string
		rev  string
		want bool
	}{
		{name: "lowercase full SHA", rev: "0123456789abcdef0123456789abcdef01234567", want: true},
		{name: "mixed case full SHA", rev: "0123456789ABCDEF0123456789abcdefABCDEF01", want: true},
		{name: "empty", rev: ""},
		{name: "short SHA", rev: "0123456"},
		{name: "41 hex characters", rev: "0123456789abcdef0123456789abcdef012345678"},
		{name: "40 characters with non-hex letter", rev: "0123456789abcdef0123456789abcdef0123456g"},
		{name: "40 character branch name", rev: "release/2026-10-30/candidate-build-00001"},
	} {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, isCommitSHA(test.rev))
		})
	}
}

func TestGitRepoBuildURL(t *testing.T) {
	for _, test := range []struct {
		name     string
		repoType string
		revision string
		path     string
		want     string
	}{
		{
			name:     "empty revision defaults to main",
			repoType: "model",
			path:     "README.md",
			want:     "/test-project/test-model/resolve/main/README.md",
		},
		{
			name:     "short branch name",
			repoType: "models",
			revision: "dev",
			path:     "weights/model.safetensors",
			want:     "/test-project/test-model/resolve/dev/weights/model.safetensors",
		},
		{
			name:     "qualified branch ref",
			repoType: "model",
			revision: "refs/heads/dev",
			path:     "README.md",
			want:     "/test-project/test-model/resolve/dev/README.md",
		},
		{
			name:     "qualified tag ref",
			repoType: "model",
			revision: "refs/tags/v1.0",
			path:     "README.md",
			want:     "/test-project/test-model/resolve/v1.0/README.md",
		},
		{
			name:     "commit SHA",
			repoType: "model",
			revision: "0123456789abcdef0123456789abcdef01234567",
			path:     "README.md",
			want:     "/test-project/test-model/resolve/0123456789abcdef0123456789abcdef01234567/README.md",
		},
		{
			name:     "dataset prefix",
			repoType: "dataset",
			revision: "main",
			path:     "data/train.csv",
			want:     "/datasets/test-project/test-model/resolve/main/data/train.csv",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := (&gitRepo{}).buildURL(test.repoType, repoTestProject, repoTestName, test.revision, test.path)
			require.Equal(t, test.want, got)
		})
	}
}

// refTestRepo is a repository where the short name "same" is both a branch
// pointing at the first commit and a tag pointing at the second one.
type refTestRepo struct {
	git     *gitRepo
	repo    *repository.Repository
	initial string
	first   string
	second  string
}

func newRefTestRepo(t *testing.T) refTestRepo {
	t.Helper()
	ctx := context.Background()
	store := hfdstorage.NewStorage(hfdstorage.WithRootDir(t.TempDir()))
	g := NewGitDB(store, nil).(*gitRepo)
	require.NoError(t, g.CreateRepository(ctx, "model", repoTestProject, repoTestName))

	commit := func(message string, files map[string]string) string {
		ops := make([]git.CommitOperation, 0, len(files))
		for path, content := range files {
			ops = append(ops, git.CommitOperation{Type: git.CommitOperationAdd, Path: path, Content: []byte(content)})
		}
		hash, err := g.CreateCommit(ctx, "model", repoTestProject, repoTestName, "main", &git.Commit{
			Message:     message,
			AuthorName:  "Test",
			AuthorEmail: "test@example.com",
		}, ops)
