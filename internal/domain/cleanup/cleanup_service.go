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

package cleanup

import (
	"context"

	"github.com/matrixhub-ai/matrixhub/internal/domain/dataset"
	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
)

// CleanupService implements the cleanup service.
type CleanupService struct {
	modelRepo   model.IModelRepo
	datasetRepo dataset.IDatasetRepo
	gitRepo     git.IGitRepo
}

// ICleanupService defines the service interface for cleanup operations.
type ICleanupService interface {
	// PreviewCleanup previews orphaned data without deleting.
	PreviewCleanup(ctx context.Context, includeRepos, includeLFS bool) (*CleanupPreview, error)
	// ExecuteCleanup executes cleanup based on options.
	ExecuteCleanup(ctx context.Context, cleanRepos, cleanLFS bool, dryRun bool) (*CleanupResult, error)
	// GetStorageStats returns storage statistics.
