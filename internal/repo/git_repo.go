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
}

// ListRevisions returns all branches and tags for a model
func (g *gitRepo) ListRevisions(ctx context.Context, repoType, project, name string) (*git.Revisions, error) {
	gitPath := g.gitPath(repoType, project, name)
	if !repository.IsRepository(gitPath) {
		return nil, fmt.Errorf("repository does not exist at %s", gitPath)
	}
	repo, err := repository.Open(gitPath)
	if err != nil {
		return nil, err
	}

	revisions := &git.Revisions{
		Branches: []*git.Revision{},
		Tags:     []*git.Revision{},
	}
	branches, err := repo.Branches()
	if err != nil {
		return nil, err
	}
	for _, branch := range branches {
		revisions.Branches = append(revisions.Branches, &git.Revision{
			Name: branch,
		})
	}

	tags, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		revisions.Tags = append(revisions.Tags, &git.Revision{
			Name: tag,
		})
	}

	return revisions, nil
}

// ListCommits returns the commit history for a model
func (g *gitRepo) ListCommits(ctx context.Context, repoType, project, name, revision string, page, pageSize int) ([]*git.Commit, int64, error) {
	gitPath := g.gitPath(repoType, project, name)
	if !repository.IsRepository(gitPath) {
		return nil, 0, fmt.Errorf("repository does not exist at %s", gitPath)
	}
	repo, err := repository.Open(gitPath)
	if err != nil {
		return nil, 0, err
	}

	revision = resolveRef(repo, revision)
	commits, err := repo.Commits(revision, nil)
	if err != nil {
		return nil, 0, err
	}

	total := int64(len(commits))
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(commits) {
		return []*git.Commit{}, total, nil
	}
	if end > len(commits) {
		end = len(commits)
	}
	commits = commits[start:end]

	var gitCommits []*git.Commit
	for _, c := range commits {
		gitCommits = append(gitCommits, &git.Commit{
			ID:             c.Hash().String(),
			Message:        c.Message(),
			AuthorName:     c.Author().Name(),
			AuthorEmail:    c.Author().Email(),
			AuthorDate:     c.Author().When(),
			CommitterName:  c.Committer().Name(),
			CommitterEmail: c.Committer().Email(),
			CommitterDate:  c.Committer().When(),
			CreatedAt:      c.Committer().When(),
		})
	}

	return gitCommits, total, nil
}

// GetCommit returns a specific commit by ID
func (g *gitRepo) GetCommit(ctx context.Context, repoType, project, name, commitID string) (*git.Commit, error) {
	gitPath := g.gitPath(repoType, project, name)
	if !repository.IsRepository(gitPath) {
		return nil, fmt.Errorf("repository does not exist at %s", gitPath)
	}
	repo, err := repository.Open(gitPath)
	if err != nil {
		return nil, err
	}

	commits, err := repo.Commits(commitID, &repository.CommitsOptions{
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("commit not found: %s", commitID)
	}
	c := commits[0]

	diff, _ := c.Diff()

	return &git.Commit{
		ID:             c.Hash().String(),
		Message:        c.Message(),
		Diff:           diff,
		AuthorName:     c.Author().Name(),
		AuthorEmail:    c.Author().Email(),
		AuthorDate:     c.Author().When(),
		CommitterName:  c.Committer().Name(),
		CommitterEmail: c.Committer().Email(),
		CommitterDate:  c.Committer().When(),
		CreatedAt:      c.Committer().When(),
	}, nil
}

func (g *gitRepo) CreateCommit(ctx context.Context, repoType, project, name, revision string, commit *git.Commit, ops []git.CommitOperation) (string, error) {
	gitPath := g.gitPath(repoType, project, name)
	if !repository.IsRepository(gitPath) {
		return "", fmt.Errorf("repository does not exist at %s", gitPath)
	}
	repo, err := repository.Open(gitPath)
	if err != nil {
		return "", err
	}
	commitHash, err := repo.CreateCommit(ctx, revision, commit.Message, commit.AuthorName, commit.AuthorEmail, convertCommitOperations(ops), commit.ParentCommit)
	if err != nil {
		return "", err
	}
	return commitHash, nil
}

// GetTree returns the file tree at a specific revision and path
func (g *gitRepo) GetTree(ctx context.Context, repoType, project, name, revision, path string) ([]*git.TreeEntry, error) {
	gitPath := g.gitPath(repoType, project, name)
	if !repository.IsRepository(gitPath) {
		return nil, fmt.Errorf("repository does not exist at %s", gitPath)
	}
	repo, err := repository.Open(gitPath)
	if err != nil {
		return nil, err
	}

	entries, err := repo.Tree(resolveRef(repo, revision), path, &repository.TreeOptions{
		Recursive: false,
	})
	if err != nil {
		return nil, err
	}

	var treeEntries []*git.TreeEntry
	for _, e := range entries {

		if e.Type() == repository.EntryTypeFile {
			blob, err := e.Blob()
			if err != nil {
				return nil, err
			}
			size := blob.Size()
			lfsPointer, _ := blob.LFSPointer()
			if lfsPointer != nil {
				size = lfsPointer.Size()
			}

			lastCommit := e.LastCommit()

			treeEntries = append(treeEntries, &git.TreeEntry{
				Name:  blob.Name(),
				Path:  e.Path(),
				Hash:  blob.Hash().String(),
				Type:  git.FileTypeFile,
				Size:  size,
				IsLFS: lfsPointer != nil,
				URL:   g.buildURL(repoType, project, name, revision, e.Path()),
				Commit: &git.Commit{
					ID:             lastCommit.Hash().String(),
					Message:        lastCommit.Message(),
					AuthorName:     lastCommit.Author().Name(),
					AuthorEmail:    lastCommit.Author().Email(),
					AuthorDate:     lastCommit.Author().When(),
					CommitterName:  lastCommit.Committer().Name(),
					CommitterEmail: lastCommit.Committer().Email(),
					CommitterDate:  lastCommit.Committer().When(),
					CreatedAt:      lastCommit.Committer().When(),
