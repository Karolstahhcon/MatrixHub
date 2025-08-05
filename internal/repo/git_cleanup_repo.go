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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/matrixhub-ai/hfd/pkg/repository"
	"golang.org/x/sync/errgroup"

	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
)

type lfsObjectInfo struct {
	size int64
	path string
}

var lfsOIDPattern = regexp.MustCompile("^[0-9a-f]{64}$")

// FindOrphanedRepos finds orphaned Git repositories on disk.
func (g *gitRepo) FindOrphanedRepos(ctx context.Context, validModelPaths, validDatasetPaths []string) ([]*git.OrphanedRepo, error) {
	validPaths := make(map[string]bool)
	for _, p := range validModelPaths {
		validPaths[p+".git"] = true
	}
	for _, p := range validDatasetPaths {
		validPaths["datasets/"+p+".git"] = true
	}

	reposDir := g.storage.RepositoriesDir()
	orphaned := []*git.OrphanedRepo{}

	err := filepath.Walk(reposDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if !repository.IsRepository(path) {
			return nil
		}

		relPath := strings.TrimPrefix(path, reposDir+"/")
		if validPaths[relPath] {
			return filepath.SkipDir
		}

		parts := strings.Split(relPath, "/")
		repoType := "model"
		if strings.HasPrefix(relPath, "datasets/") {
			repoType = "dataset"
		}

		var projectName, resourceName string
		if len(parts) >= 2 {
			if repoType == "dataset" && len(parts) >= 3 {
				projectName = parts[1]
				resourceName = strings.TrimSuffix(parts[2], ".git")
			} else {
				projectName = parts[0]
				resourceName = strings.TrimSuffix(parts[1], ".git")
			}
		}

		orphaned = append(orphaned, &git.OrphanedRepo{
			Path:         relPath,
			Type:         repoType,
			ProjectName:  projectName,
			ResourceName: resourceName,
			SizeBytes:    calculateDirSize(path),
		})

		return filepath.SkipDir
	})

	return orphaned, err
}

// FindOrphanedLFS finds orphaned LFS objects on disk.
func (g *gitRepo) FindOrphanedLFS(ctx context.Context) ([]*git.OrphanedLFS, error) {
	allOIDs, err := g.scanLFSObjects(ctx)
	if err != nil {
		return nil, err
	}

	referencedOIDs, err := g.collectReferencedOIDs(ctx)
	if err != nil {
		return nil, err
	}

	orphaned := make([]*git.OrphanedLFS, 0)
	for oid, info := range allOIDs {
		if !referencedOIDs[oid] {
			orphaned = append(orphaned, &git.OrphanedLFS{
				OID:       oid,
				SizeBytes: info.size,
				Path:      info.path,
			})
		}
	}
