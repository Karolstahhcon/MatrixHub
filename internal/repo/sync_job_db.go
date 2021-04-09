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

	"github.com/matrixhub-ai/matrixhub/internal/domain/syncjob"
)

type syncJobDB struct {
	db *gorm.DB
}

// NewSyncJobDB creates a new syncJobDB instance
func NewSyncJobDB(db *gorm.DB) syncjob.ISyncJobRepo {
	return &syncJobDB{db: db}
}

// CreateSyncJob creates a new sync job
func (r *syncJobDB) CreateSyncJob(ctx context.Context, job *syncjob.SyncJob) error {
	return r.db.WithContext(ctx).Create(job).Error
}

// CreateSyncJobs creates a batch atomically so processors never observe a partially-created task.
func (r *syncJobDB) CreateSyncJobs(ctx context.Context, jobs []*syncjob.SyncJob) error {
	if len(jobs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(&jobs).Error
	})
}

// GetSyncJob gets a sync job by ID
func (r *syncJobDB) GetSyncJob(ctx context.Context, id int) (*syncjob.SyncJob, error) {
	var job syncjob.SyncJob
	err := r.db.WithContext(ctx).First(&job, id).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// UpdateSyncJob updates a sync job
func (r *syncJobDB) UpdateSyncJob(ctx context.Context, job *syncjob.SyncJob) error {
