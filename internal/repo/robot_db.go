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

	"github.com/matrixhub-ai/matrixhub/internal/domain/robot"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

type robotRepo struct {
	db *gorm.DB
}

func (r *robotRepo) GetRobotByName(ctx context.Context, name string) (*robot.Robot, error) {
	var rb robot.Robot
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&rb).Error
	if err != nil {
		return nil, err
	}
	return &rb, nil
}

func (r *robotRepo) UpdateRobot(ctx context.Context, robot *robot.Robot) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(robot).
			Where("id = ?", robot.ID).
			Updates(map[string]interface{}{
				"description":          robot.Description,
				"enabled":              robot.Enabled,
				"expire_at":            robot.ExpireAt,
				"duration":             robot.Duration,
				"token_hash":           robot.TokenHash,
				"project_scope":        robot.ProjectScope,
				"platform_permissions": robot.PlatformPermissions,
				"project_permissions":  robot.ProjectPermissions,
			}).Error; err != nil {
			return err
