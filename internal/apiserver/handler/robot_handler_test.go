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

package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/auth"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	projectmocks "github.com/matrixhub-ai/matrixhub/internal/domain/project/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/robot"
	robotmocks "github.com/matrixhub-ai/matrixhub/internal/domain/robot/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

func robotUserContext() context.Context {
	return auth.WithIdentity(context.Background(), user.NewUserIdentity(42, "alice"))
}

func TestRobotHandlerCreateRobotAccount(t *testing.T) {
	ctx := robotUserContext()
	request := &v1alpha1.CreateRobotAccountRequest{
		Name:                "ci",
		Description:         "build robot",
		ExpireDays:          7,
		PlatformPermissions: []string{role.UserCreate.String()},
		ProjectPermissions:  []string{role.ModelPull.String()},
		Projects:            []string{"demo", "demo"},
		ProjectScope:        v1alpha1.RobotAccountProjectScope_ROBOT_ACCOUNT_PROJECT_SCOPE_SELECTED,
	}

	t.Run("creates a selected-project robot", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		robotRepo := robotmocks.NewMockIRobotRepo(ctrl)
		projectRepo := projectmocks.NewMockIProjectRepo(ctrl)
		robotRepo.EXPECT().GetRobotByName(ctx, "robot$ci").Return(nil, gorm.ErrRecordNotFound)
		projectRepo.EXPECT().ListProjectInfoByNames(ctx, []string{"demo"}).Return([]*project.Project{{ID: 1, Name: "demo"}}, nil)
		robotRepo.EXPECT().CreateRobot(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, rb *robot.Robot) error {
			require.Equal(t, "robot$ci", rb.Name)
			require.Equal(t, 42, rb.CreateBy)
			require.Equal(t, robot.ProjectScopeSelected, rb.ProjectScope)
