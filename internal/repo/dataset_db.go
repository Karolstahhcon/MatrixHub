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

	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/internal/domain/dataset"
	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
)

type datasetDB struct {
	db *gorm.DB
}

// NewDatasetDB creates a new DatasetRepo instance
func NewDatasetDB(db *gorm.DB) dataset.IDatasetRepo {
	return &datasetDB{db: db}
}

// ListAllPaths returns all valid dataset paths (project/name format) from database.
// Datasets without a valid project reference are excluded; those are data integrity
// issues, not orphaned repositories.
func (d *datasetDB) ListAllPaths(ctx context.Context) ([]string, error) {
	var rows []struct {
		ProjectName string
		Name        string
	}
	if err := d.db.WithContext(ctx).
		Table("datasets d").
		Select("p.name AS project_name, d.name AS name").
		Joins("LEFT JOIN projects p ON d.project_id = p.id").
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

// List retrieves datasets with filtering and pagination
func (d *datasetDB) List(ctx context.Context, filter *model.Filter) ([]*dataset.Dataset, int64, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	// Build base query with JOIN to projects
	query := d.db.WithContext(ctx).
		Table("datasets d").
		Select(`d.id, d.name, d.project_id, d.default_branch,
				d.num_rows, d.size, d.readme_content,
				d.is_popular,
				d.created_at, d.updated_at,
				p.name as project_name`).
		Joins("INNER JOIN projects p ON d.project_id = p.id")

	// Apply filters
	if filter.Project != "" {
		query = query.Where("p.name = ?", filter.Project)
	}

	if filter.Search != "" {
		pattern := "%" + filter.Search + "%"
		if d.db.Name() == "sqlite" {
			query = query.Where("instr(p.name, ?) > 0 OR d.name LIKE ?", filter.Search, pattern)
		} else {
			query = query.Where("p.name LIKE ? OR d.name LIKE ?", pattern, pattern)
		}
	}

	if len(filter.Label) > 0 {
		for _, label := range filter.Label {
			query = query.Where(`EXISTS (
				SELECT 1 FROM datasets_labels dl
				INNER JOIN labels l ON dl.label_id = l.id
				WHERE dl.dataset_id = d.id AND l.name = ?
			)`, label)
		}
	}

	// Apply popular filter
