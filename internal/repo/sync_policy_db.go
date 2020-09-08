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

type syncPolicyDB struct {
	db *gorm.DB
}

// NewSyncPolicyDB creates a new syncPolicyDB instance
func NewSyncPolicyDB(db *gorm.DB) syncpolicy.ISyncPolicyRepo {
	return &syncPolicyDB{db: db}
}

// CreateSyncPolicy creates a new sync policy
func (r *syncPolicyDB) CreateSyncPolicy(ctx context.Context, policy *syncpolicy.SyncPolicy) error {
	return r.db.WithContext(ctx).Create(policy).Error
}

// GetSyncPolicy gets a sync policy by ID
func (r *syncPolicyDB) GetSyncPolicy(ctx context.Context, id int) (*syncpolicy.SyncPolicy, error) {
	var policy syncpolicy.SyncPolicy
	err := r.db.WithContext(ctx).First(&policy, id).Error
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// UpdateSyncPolicy updates a sync policy
func (r *syncPolicyDB) UpdateSyncPolicy(ctx context.Context, policy *syncpolicy.SyncPolicy) error {
	return r.db.WithContext(ctx).Save(policy).Error
}

// DeleteSyncPolicy deletes a sync policy by ID
func (r *syncPolicyDB) DeleteSyncPolicy(ctx context.Context, id int) error {
