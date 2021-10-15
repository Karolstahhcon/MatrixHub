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

	if _, exists, err := s.findModel(ctx, project, name); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.New("model already exists")
	}

	repositoryExists, err := s.gitRepo.RepositoryExists(ctx, "models", project, name)
	if err != nil {
		return nil, err
	}
	if !repositoryExists {
		if err := s.gitRepo.CreateRepository(ctx, "models", project, name); err != nil {
			return nil, err
		}
	}

	return s.createModelRecord(ctx, project, name)
}

// GetModel retrieves a model by project and name.
func (s *ModelService) GetModel(ctx context.Context, project, name string) (*Model, error) {
	if project == "" {
		return nil, errors.New("invalid project")
	}
	if name == "" {
		return nil, errors.New("invalid input")
	}

	mod, err := s.modelRepo.GetByProjectAndName(ctx, project, name)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

// ListModels returns a paginated list of models with optional filtering.
func (s *ModelService) ListModels(ctx context.Context, filter *Filter) ([]*Model, int64, error) {
	if filter == nil {
		filter = &Filter{
			Page:     1,
			PageSize: 20,
		}
	}

	return s.modelRepo.List(ctx, filter)
}

// DeleteModel deletes a model by project and name.
func (s *ModelService) DeleteModel(ctx context.Context, project, name string) error {
	if project == "" {
		return errors.New("invalid project")
	}
	if name == "" {
		return errors.New("invalid input")
	}

	// First delete the Git repository, then delete the model record in the database.
	if err := s.gitRepo.DeleteRepository(ctx, "models", project, name); err != nil {
		return err
	}

	return s.modelRepo.Delete(ctx, project, name)
}

// ListModelTaskLabels returns all task labels for models.
func (s *ModelService) ListModelTaskLabels(ctx context.Context) ([]*Label, error) {
	return s.labelRepo.ListByCategoryAndScope(ctx, "task", "model")
}

// ListModelFrameLabels returns all framework labels for models.
func (s *ModelService) ListModelFrameLabels(ctx context.Context) ([]*Label, error) {
	return s.labelRepo.ListByCategoryAndScope(ctx, "library", "model")
}

// ListModelRevisions returns all branches and tags for a model.
func (s *ModelService) ListModelRevisions(ctx context.Context, project, name string) (*git.Revisions, error) {
	if project == "" {
		return nil, errors.New("invalid project")
	}
	if name == "" {
		return nil, errors.New("invalid input")
	}

	// Check if model exists in database first
	_, err := s.modelRepo.GetByProjectAndName(ctx, project, name)
	if err != nil {
		return nil, err
	}

	return s.gitRepo.ListRevisions(ctx, "models", project, name)
}

