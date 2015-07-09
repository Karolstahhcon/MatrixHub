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
		require.NoError(t, err)
		return hash
	}
	first := commit("first", map[string]string{"README.md": "v1", "weights/model.bin": "weights"})
	second := commit("second", map[string]string{"README.md": "second"})

	repo, err := repository.Open(store.ResolvePath(repoTestProject + "/" + repoTestName))
	require.NoError(t, err)
	evictRepoTestRepository(t, repo)
	initial, err := repo.ResolveRevision(first + "~1")
	require.NoError(t, err)
	require.NoError(t, repo.CreateBranch("same", first))
	require.NoError(t, repo.CreateTag("same", second))
	require.NoError(t, repo.CreateTag("v1.0", first))

	return refTestRepo{git: g, repo: repo, initial: initial, first: first, second: second}
}

func TestResolveRef(t *testing.T) {
	fixture := newRefTestRepo(t)
	// Like git, a full SHA names the commit even when a branch has that name.
	require.NoError(t, fixture.repo.CreateBranch(fixture.second, fixture.first))

	for _, test := range []struct {
		name string
		rev  string
		want string
	}{
		{name: "empty revision", rev: "", want: ""},
		{name: "qualified branch ref", rev: "refs/heads/main", want: "refs/heads/main"},
		{name: "qualified ref is not looked up", rev: "refs/tags/missing", want: "refs/tags/missing"},
		{name: "full commit SHA", rev: fixture.first, want: fixture.first},
		{name: "full commit SHA wins over a branch with that name", rev: fixture.second, want: fixture.second},
		{name: "branch", rev: "main", want: "refs/heads/main"},
		{name: "tag", rev: "v1.0", want: "refs/tags/v1.0"},
		{name: "branch wins over tag with the same name", rev: "same", want: "refs/heads/same"},
		{name: "short commit SHA", rev: fixture.first[:7], want: fixture.first[:7]},
		{name: "unknown revision", rev: "missing", want: "missing"},
	} {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, resolveRef(fixture.repo, test.rev))
		})
	}
}

func TestGitRepoListCommitsResolvesRevisions(t *testing.T) {
	ctx := context.Background()
	fixture := newRefTestRepo(t)

	for _, test := range []struct {
		name      string
		revision  string
		page      int
		pageSize  int
		wantIDs   []string
		wantTotal int64
		wantErr   string
	}{
		{
			name:      "branch wins over tag with the same name",
			revision:  "same",
			page:      1,
			pageSize:  10,
			wantIDs:   []string{fixture.first, fixture.initial},
			wantTotal: 2,
		},
		{
			name:      "qualified tag ref",
			revision:  "refs/tags/same",
			page:      1,
			pageSize:  10,
			wantIDs:   []string{fixture.second, fixture.first, fixture.initial},
			wantTotal: 3,
		},
		{
			name:      "tag",
			revision:  "v1.0",
			page:      1,
			pageSize:  10,
			wantIDs:   []string{fixture.first, fixture.initial},
			wantTotal: 2,
		},
		{
			name:      "short commit SHA",
			revision:  fixture.first[:7],
			page:      1,
			pageSize:  10,
			wantIDs:   []string{fixture.first, fixture.initial},
			wantTotal: 2,
		},
		{
			name:     "unknown revision",
			revision: "missing",
			page:     1,
			pageSize: 10,
			wantErr:  "failed to resolve revision",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			commits, total, err := fixture.git.ListCommits(ctx, "model", repoTestProject, repoTestName, test.revision, test.page, test.pageSize)
			if test.wantErr != "" {
				require.ErrorContains(t, err, test.wantErr)
				require.Nil(t, commits)
				require.Zero(t, total)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.wantTotal, total)
			ids := make([]string, 0, len(commits))
			for _, commit := range commits {
				ids = append(ids, commit.ID)
			}
			require.Equal(t, test.wantIDs, ids)
		})
	}
}

func TestGitRepoGetBlobResolvesRevisions(t *testing.T) {
	ctx := context.Background()
	fixture := newRefTestRepo(t)

	for _, test := range []struct {
		name       string
		revision   string
		path       string
		want       git.TreeEntry
		wantCommit string
		wantErr    string
	}{
		{
			name:     "branch wins over tag with the same name",
			revision: "same",
			path:     "README.md",
			want: git.TreeEntry{
				Name: "README.md", Path: "README.md", Type: git.FileTypeFile, Size: 2,
				URL: "/test-project/test-model/resolve/same/README.md",
			},
			wantCommit: fixture.first,
		},
		{
			name:     "tag",
			revision: "v1.0",
			path:     "README.md",
			want: git.TreeEntry{
				Name: "README.md", Path: "README.md", Type: git.FileTypeFile, Size: 2,
				URL: "/test-project/test-model/resolve/v1.0/README.md",
			},
			wantCommit: fixture.first,
		},
		{
			name:     "default branch",
			revision: "",
			path:     "README.md",
			want: git.TreeEntry{
				Name: "README.md", Path: "README.md", Type: git.FileTypeFile, Size: 6,
				URL: "/test-project/test-model/resolve/main/README.md",
