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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	stdpath "path"
	"strings"
	"time"

	hfdlfs "github.com/matrixhub-ai/hfd/pkg/lfs"
	"github.com/matrixhub-ai/hfd/pkg/mirror"
	"github.com/matrixhub-ai/hfd/pkg/repository"
	"github.com/matrixhub-ai/hfd/pkg/storage"

	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
)

type gitRepo struct {
	storage *storage.Storage
	mirror  *mirror.Mirror
}

const maxSafetensorsHeaderBytes uint64 = 64 * 1024 * 1024

// safetensorsHeaderReadBudget bounds the time one metadata extraction may spend
// reading safetensors headers, across all files rather than per file. A read
// served by the mirror tee cache blocks until the upstream download reaches the
// header bytes, so extraction must not wait on it indefinitely: it runs inline on
// the proxy sync path that git and Hugging Face clients wait for. A shared budget
// also keeps a sharded model from multiplying the wait by its shard count.
const safetensorsHeaderReadBudget = 10 * time.Second

type safetensorsIndexFiles struct {
	Metadata struct {
		TotalSize int64 `json:"total_size"`
	} `json:"metadata"`
	WeightMap map[string]string `json:"weight_map"`
}

// NewGitDB creates a new GitRepo instance
func NewGitDB(storage *storage.Storage, mirror *mirror.Mirror) git.IGitRepo {
	return &gitRepo{
		storage: storage,
		mirror:  mirror,
	}
}

func repoPrefix(repoType string) string {
	switch repoType {
	case "model", "models":
		return ""
	case "dataset", "datasets":
		return "datasets/"
	case "space", "spaces":
		return "spaces/"
	default:
		return ""
	}
}

func (g *gitRepo) gitPath(repoType string, project, name string) string {
	repoName := repoPrefix(repoType) + project + "/" + name
	repoPath := g.storage.ResolvePath(repoName)
	return repoPath

}

func (g *gitRepo) buildURL(repoType, project, name, revision, path string) string {
	if revision == "" {
		revision = "main"
	} else if strings.Contains(revision, "/") {
		// If revision is a ref like refs/heads/main or refs/tags/v1, keep only the last part
		revision = stdpath.Base(revision)
	}

	return fmt.Sprintf("/%s/%s/resolve/%s/%s", repoPrefix(repoType)+project, name, revision, path)
}

// resolveRef disambiguates a short revision name by trying refs/heads/ before refs/tags/.
// Already-qualified refs (refs/heads/..., refs/tags/...) and 40-char SHAs pass through unchanged.
func resolveRef(repo *repository.Repository, rev string) string {
	if rev == "" || strings.HasPrefix(rev, "refs/") || isCommitSHA(rev) {
		return rev
	}
	if _, err := repo.ResolveRevision("refs/heads/" + rev); err == nil {
		return "refs/heads/" + rev
	}
	if _, err := repo.ResolveRevision("refs/tags/" + rev); err == nil {
		return "refs/tags/" + rev
	}
	return rev
}

// convertCommitOperations converts domain CommitOperation slice to repository CommitOperation slice.
func convertCommitOperations(ops []git.CommitOperation) []repository.CommitOperation {
	result := make([]repository.CommitOperation, len(ops))
	for i, op := range ops {
		var typ repository.CommitOperationType
		switch op.Type {
		case git.CommitOperationAdd:
			typ = repository.CommitOperationAdd
		case git.CommitOperationDelete:
			typ = repository.CommitOperationDelete
		}
		result[i] = repository.CommitOperation{
			Type:    typ,
			Path:    op.Path,
			Content: op.Content,
		}
	}
	return result
}

func isCommitSHA(s string) bool {
	if len(s) != 40 {
		return false
	}
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}

// CreateRepository initializes a Git repository
func (g *gitRepo) CreateRepository(ctx context.Context, repoType, project, name string) error {
	gitPath := g.gitPath(repoType, project, name)
	if repository.IsRepository(gitPath) {
		return fmt.Errorf("repository already exists at %s", gitPath)
	}

	defaultBranch := "main"
	repo, err := repository.Init(ctx, gitPath, defaultBranch)
	if err != nil {
		return err
	}

	// TODO(@scyda): use actual user info from context when available
	user := "HuggingFace"
	email := "hf@users.noreply.huggingface.co"

	// Create initial commit with default .gitattributes
	_, err = repo.CreateCommit(ctx, defaultBranch, "Initial commit", user, email, []repository.CommitOperation{
		{
			Type:    repository.CommitOperationAdd,
			Path:    repository.GitattributesFileName,
			Content: repository.GitattributesText,
		},
	}, "")
	if err != nil {
		_ = repo.Remove()
		return err
	}

	return nil
}

// RepositoryExists checks whether a Git repository exists on disk.
func (g *gitRepo) RepositoryExists(ctx context.Context, repoType, project, name string) (bool, error) {
	return repository.IsRepository(g.gitPath(repoType, project, name)), nil
}

// DeleteRepository removes the Git repository
func (g *gitRepo) DeleteRepository(ctx context.Context, repoType, project, name string) error {
	gitPath := g.gitPath(repoType, project, name)
	if !repository.IsRepository(gitPath) {
		return fmt.Errorf("repository does not exist at %s", gitPath)
	}
	repo, err := repository.Open(gitPath)
	if err != nil {
		return err
	}
	return repo.Remove()
