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

package model

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registry"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
)

// IModelService defines the service interface for model operations.
//
//go:generate go tool mockgen -source=model_service.go -destination=mocks/model_service_mock.go -package=mocks
type IModelService interface {
	// Model CRUD operations
	CreateModel(ctx context.Context, project, name string) (*Model, error)
	GetModel(ctx context.Context, project, name string) (*Model, error)
	ListModels(ctx context.Context, filter *Filter) ([]*Model, int64, error)
	DeleteModel(ctx context.Context, project, name string) error

	// EnsureModel ensures that a model exists in the database for the
	// given namespace and name.
	EnsureModel(ctx context.Context, project, name string) (*Model, error)

	// Label operations
	ListModelTaskLabels(ctx context.Context) ([]*Label, error)
	ListModelFrameLabels(ctx context.Context) ([]*Label, error)

	// Git operations
	ListModelRevisions(ctx context.Context, project, name string) (*git.Revisions, error)
	ListModelCommits(ctx context.Context, project, name, revision string, page, pageSize int) ([]*git.Commit, int64, error)
	GetModelCommit(ctx context.Context, project, name, commitID string) (*git.Commit, error)
	CreateModelCommit(ctx context.Context, project, name, revision string, commit *git.Commit, ops []git.CommitOperation) (string, error)
	GetModelTree(ctx context.Context, project, name, revision, path string) ([]*git.TreeEntry, error)
	GetModelBlob(ctx context.Context, project, name, revision, path string) (*git.TreeEntry, error)

	// Settings
	UpdateModelSetting(ctx context.Context, project, name string, update *SettingUpdate) error

	// Metadata sync
	SyncMetadata(ctx context.Context, project, name string) error

	// CheckOrSyncFromRemote ensures the model exists and synchronizes proxy models
	// from their remote registry when needed.
	CheckOrSyncFromRemote(ctx context.Context, project, name string) error
}

// ModelService implements the model service operations.
type ModelService struct {
	modelRepo    IModelRepo
	labelRepo    ILabelRepo
	gitRepo      git.IGitRepo
	projectRepo  project.IProjectRepo
	registryRepo registry.IRegistryRepo
}

// NewModelService creates a new ModelService instance.
func NewModelService(modelRepo IModelRepo, labelRepo ILabelRepo, gitRepo git.IGitRepo, projectRepo project.IProjectRepo, registryRepo registry.IRegistryRepo) IModelService {
	return &ModelService{
		modelRepo:    modelRepo,
		labelRepo:    labelRepo,
		gitRepo:      gitRepo,
		projectRepo:  projectRepo,
		registryRepo: registryRepo,
	}
}

// CreateModel creates a new model in the system.
func (s *ModelService) CreateModel(ctx context.Context, project, name string) (*Model, error) {
	if err := validateModelPath(project, name); err != nil {
		return nil, err
	}

