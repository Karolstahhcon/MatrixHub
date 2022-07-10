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
	"errors"

	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

type ProjectDBRepo struct {
	db *gorm.DB
}

var _ project.IProjectRepo = (*ProjectDBRepo)(nil)

func NewProjectDBRepo(db *gorm.DB) *ProjectDBRepo {
	return &ProjectDBRepo{db: db}
}

func (r *ProjectDBRepo) CreateProject(ctx context.Context, param *project.Project) (*project.Project, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(param).Error; err != nil {
			return err
		}

		creatorID := user.GetCurrentUserId(ctx)
		if creatorID != 0 {
			member := &project.ProjectMember{
				ProjectID:  &param.ID,
				MemberID:   creatorID,
				MemberType: project.MemberTypeUser,
				RoleID:     role.ProjectRoleAdmin,
			}
			if err := tx.Create(member).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return param, nil
}

func (r *ProjectDBRepo) GetProjectByID(ctx context.Context, id int) (*project.Project, error) {
	var p project.Project
	err := r.db.WithContext(ctx).
		Select(`projects.*,
			(SELECT COUNT(*) FROM models WHERE models.project_id = projects.id) as model_count,
			(SELECT COUNT(*) FROM datasets WHERE datasets.project_id = projects.id) as dataset_count,
			registries.url as registry_url`).
		Joins("LEFT JOIN registries ON registries.id = projects.registry_id").
		Where("projects.id = ?", id).
		First(&p).Error
	if err != nil {
		return nil, err
	}
