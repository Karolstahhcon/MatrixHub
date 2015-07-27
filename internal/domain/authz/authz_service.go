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

package authz

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/internal/domain/auth"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	"github.com/matrixhub-ai/matrixhub/internal/domain/robot"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
)

//go:generate go tool mockgen -source=authz_service.go -destination=mocks/authz_service_mock.go -package=mocks

// IAuthzService permission verification service interface
type IAuthzService interface {
	// GetUserPermissions gets user's permission list in a project
	GetUserPermissions(ctx context.Context, userID int, projectID int) ([]role.Permission, error)

	// VerifyPlatformPermission verifies platform-level permission (gets user/robot info from ctx)
	VerifyPlatformPermission(ctx context.Context, perm role.Permission) (bool, error)

	VerifyProjectPermission(ctx context.Context, projectID int, perm role.Permission) (bool, error)

	// VerifyProjectPermissionByName resolves project name to ID, then verifies permission
	VerifyProjectPermissionByName(ctx context.Context, projectName string, perm role.Permission) (bool, error)

	// GetUserAccessibleProjectIDs gets all project IDs accessible to a user
	GetUserAccessibleProjectIDs(ctx context.Context, userID int) ([]int, error)
}

// AuthzService permission verification service
type AuthzService struct {
	authzRepo   IAuthzProjectRepo
	projectRepo project.IProjectRepo
	robotRepo   robot.IRobotRepo
}

// NewAuthzService creates permission verification service
func NewAuthzService(authzRepo IAuthzProjectRepo, projectRepo project.IProjectRepo, robotRepo robot.IRobotRepo) IAuthzService {
	return &AuthzService{
		authzRepo:   authzRepo,
		projectRepo: projectRepo,
		robotRepo:   robotRepo,
	}
}

// GetUserPermissions gets user's permission list in a project
func (s *AuthzService) GetUserPermissions(ctx context.Context, userID int, projectID int) ([]role.Permission, error) {
	platformPerms, err := s.authzRepo.GetUserPlatformPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

