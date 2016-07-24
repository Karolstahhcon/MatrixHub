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
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
)

type modelDB struct {
	db *gorm.DB
}

// NewModelDB creates a new ModelRepo instance
func NewModelDB(db *gorm.DB) model.IModelRepo {
	return &modelDB{db: db}
}

// ListAllPaths returns all valid model paths (project/name format) from database.
// Models without a valid project reference are excluded; those are data integrity
// issues, not orphaned repositories.
func (m *modelDB) ListAllPaths(ctx context.Context) ([]string, error) {
	var rows []struct {
		ProjectName string
		Name        string
	}
	if err := m.db.WithContext(ctx).
		Table("models m").
		Select("p.name AS project_name, m.name AS name").
		Joins("LEFT JOIN projects p ON m.project_id = p.id").
		Where("p.name IS NOT NULL").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(rows))
	for _, row := range rows {
		paths = append(paths, row.ProjectName+"/"+row.Name)
	}
	return paths, nil
}

// List retrieves models with filtering and pagination
func (m *modelDB) List(ctx context.Context, filter *model.Filter) ([]*model.Model, int64, error) {
	// Build base query with JOIN to projects
	query := m.db.WithContext(ctx).
		Table("models m").
		Select(`m.id, m.name, m.project_id, m.size, m.parameter_count,
				m.readme_content, m.is_popular, m.default_branch,
				m.synced_at, m.created_at, m.updated_at,
				p.name as project_name`).
		Joins("INNER JOIN projects p ON m.project_id = p.id")

	// Apply filters
	if len(filter.ProjectIDs) > 0 {
		query = query.Where("m.project_id IN ?", filter.ProjectIDs)
	}

	if filter.Project != "" {
		query = query.Where("p.name = ?", filter.Project)
	}

	if filter.Search != "" {
		pattern := "%" + filter.Search + "%"
		if m.db.Name() == "sqlite" {
			query = query.Where("instr(p.name, ?) > 0 OR m.name LIKE ?", filter.Search, pattern)
		} else {
			query = query.Where("p.name LIKE ? OR m.name LIKE ?", pattern, pattern)
		}
	}

	if len(filter.Label) > 0 {
		for _, label := range filter.Label {
			query = query.Where(`EXISTS (
				SELECT 1 FROM models_labels ml
				INNER JOIN labels l ON ml.label_id = l.id
				WHERE ml.model_id = m.id AND l.name = ?
			)`, label)
		}
	}

	// Apply popular filter
	if filter.Popular != nil && *filter.Popular {
		query = query.Where("m.is_popular = ?", 1)
