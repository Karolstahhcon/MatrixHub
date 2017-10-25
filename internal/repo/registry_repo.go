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
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/internal/domain/registry"
)

type RegistryRepo struct {
	db *gorm.DB
}

func NewRegistryRepo(db *gorm.DB) *RegistryRepo {
	return &RegistryRepo{db}
}

func (r *RegistryRepo) ListRegistries(ctx context.Context, page, pageSize int, search string) (rs []*registry.Registry, total int64, err error) {
	baseQuery := r.db.WithContext(ctx).Model(&registry.Registry{})
	if search != "" {
		baseQuery = baseQuery.Where("name LIKE ?", "%"+search+"%")
	}
	if err = baseQuery.Count(&total).Error; err != nil {
