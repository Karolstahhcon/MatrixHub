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
