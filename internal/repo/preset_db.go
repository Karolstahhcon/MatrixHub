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

	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
	"github.com/matrixhub-ai/matrixhub/internal/domain/preset"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registry"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
)

var errPresetSeedConflict = errors.New("preset registries were inserted concurrently")

type presetDB struct {
	db *gorm.DB
}

func NewPresetDB(db *gorm.DB) preset.IPresetRepo {
	return &presetDB{db: db}
}

func (r *presetDB) SeedIfEmpty(ctx context.Context, manifest *preset.Manifest) (bool, error) {
	applied := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, table := range []string{"registries", "projects"} {
			var count int64
			if err := tx.Table(table).Count(&count).Error; err != nil {
				return err
			}
			if count != 0 {
				return nil
			}
		}

		registryIDs := make(map[string]int, len(manifest.Registries))
		for _, spec := range manifest.Registries {
			reg := &registry.Registry{
				ID: spec.ID, Name: spec.Name, Description: spec.Description,
				Type: "REGISTRY_TYPE_HUGGINGFACE", URL: spec.URL,
			}
			if err := tx.Create(reg).Error; err != nil {
				// Fixed registry IDs are the concurrency guard. Returning an
				// error here rolls back the entire losing transaction.
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					return errPresetSeedConflict
				}
				return fmt.Errorf("create preset registry %s: %w", spec.Name, err)
			}
			registryIDs[spec.Name] = reg.ID
		}

		var admin user.User
		err := tx.Where("username = ?", "admin").Take(&admin).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("find preset project administrator: %w", err)
		}
		hasAdmin := err == nil
		if !hasAdmin {
			log.Warn("preset project memberships skipped: admin user is missing")
		}

		projectIDs := make(map[string]int, len(manifest.Projects))
		for _, spec := range manifest.Projects {
			registryID, ok := registryIDs[spec.Registry]
			if !ok {
				return fmt.Errorf("preset project %s references unknown registry %s", spec.Name, spec.Registry)
			}
			prj := &project.Project{
				Name: spec.Name, Type: project.ProjectTypePublic,
				RegistryID: &registryID, Organization: spec.Organization,
			}
			if err := tx.Create(prj).Error; err != nil {
				return fmt.Errorf("create preset project %s: %w", spec.Name, err)
			}
			projectIDs[spec.Name] = prj.ID
			if hasAdmin {
				member := &project.ProjectMember{
					ProjectID: &prj.ID, MemberID: admin.ID,
					MemberType: project.MemberTypeUser, RoleID: role.ProjectRoleAdmin,
				}
				if err := tx.Create(member).Error; err != nil {
					return fmt.Errorf("create preset project membership %s: %w", spec.Name, err)
				}
			}
		}

