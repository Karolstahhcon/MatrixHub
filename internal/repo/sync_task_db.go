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

	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/internal/domain/syncpolicy"
)

type syncTaskDB struct {
	db *gorm.DB
}

// NewSyncTaskDB creates a new syncTaskDB instance
func NewSyncTaskDB(db *gorm.DB) syncpolicy.ISyncTaskRepo {
	return &syncTaskDB{db: db}
}

// CreateSyncTask creates a new sync task
func (r *syncTaskDB) CreateSyncTask(ctx context.Context, task *syncpolicy.SyncTask) (*syncpolicy.SyncTask, error) {
	if err := r.db.WithContext(ctx).Create(task).Error; err != nil {
		return nil, err
	}
	return task, nil
}

// GetSyncTask gets a sync task by ID
func (r *syncTaskDB) GetSyncTask(ctx context.Context, id int) (*syncpolicy.SyncTask, error) {
	var task syncpolicy.SyncTask
	err := r.db.WithContext(ctx).First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateSyncTask updates a sync task
func (r *syncTaskDB) UpdateSyncTask(ctx context.Context, task *syncpolicy.SyncTask) error {
