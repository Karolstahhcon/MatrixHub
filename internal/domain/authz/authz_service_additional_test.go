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

package authz_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/matrixhub-ai/matrixhub/internal/domain/auth"
	"github.com/matrixhub-ai/matrixhub/internal/domain/authz"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	"github.com/matrixhub-ai/matrixhub/internal/domain/robot"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
)

type authzRepoStub struct {
	userProject  func(context.Context, int, int) ([]role.Permission, error)
	userPlatform func(context.Context, int) ([]role.Permission, error)
	accessible   func(context.Context, int) ([]int, error)
	robotProject func(context.Context, int, int) ([]role.Permission, error)
}

func (s *authzRepoStub) GetUserProjectPermissions(ctx context.Context, userID, projectID int) ([]role.Permission, error) {
	return s.userProject(ctx, userID, projectID)
}
func (s *authzRepoStub) GetUserPlatformPermissions(ctx context.Context, userID int) ([]role.Permission, error) {
	return s.userPlatform(ctx, userID)
}
func (s *authzRepoStub) GetUserAccessibleProjectIDs(ctx context.Context, userID int) ([]int, error) {
	return s.accessible(ctx, userID)
}
func (s *authzRepoStub) GetRobotProjectPermissions(ctx context.Context, robotID, projectID int) ([]role.Permission, error) {
	return s.robotProject(ctx, robotID, projectID)
}

type completeProjectRepoStub struct {
	project.IProjectRepo
	byName func(context.Context, string) (*project.Project, error)
	byID   func(context.Context, int) (*project.Project, error)
}

func (s *completeProjectRepoStub) GetProjectByName(ctx context.Context, name string) (*project.Project, error) {
	return s.byName(ctx, name)
}
func (s *completeProjectRepoStub) GetProjectByID(ctx context.Context, id int) (*project.Project, error) {
	return s.byID(ctx, id)
}

type robotRepoStub struct {
	robot.IRobotRepo
	get func(context.Context, int) (*robot.Robot, error)
}

func (s *robotRepoStub) GetRobot(ctx context.Context, id int) (*robot.Robot, error) {
	return s.get(ctx, id)
}

type unknownIdentity struct{}

func (unknownIdentity) GetID() int       { return 1 }
func (unknownIdentity) GetName() string  { return "unknown" }
func (unknownIdentity) TypeName() string { return "unknown" }

func TestAuthzServiceUserPermissions(t *testing.T) {
	repo := &authzRepoStub{
		userPlatform: func(_ context.Context, userID int) ([]role.Permission, error) {
			require.Equal(t, 7, userID)
			return []role.Permission{role.UserGet}, nil
		},
		userProject: func(_ context.Context, userID, projectID int) ([]role.Permission, error) {
			require.Equal(t, 7, userID)
			require.Equal(t, 9, projectID)
