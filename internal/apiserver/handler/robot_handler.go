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
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	"github.com/matrixhub-ai/matrixhub/internal/domain/robot"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

type RobotHandler struct {
	robotRepo   robot.IRobotRepo
	projectRepo project.IProjectRepo
}

func (r *RobotHandler) CreateRobotAccount(ctx context.Context, request *v1alpha1.CreateRobotAccountRequest) (*v1alpha1.CreateRobotAccountResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	robotName := robot.RobotPrefix + request.Name
	_, err := r.robotRepo.GetRobotByName(ctx, robotName)
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "robot name already exists")
	}
	platformPermissions := role.StringsToPermissions(request.PlatformPermissions)
	projectPermissions := role.StringsToPermissions(request.ProjectPermissions)
	if err = role.PlatformPermissions.CheckPermissions(platformPermissions); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err = role.ProjectPermissions.CheckPermissions(projectPermissions); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, hash, err := utils.GenerateRobotToken()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	scope := robot.ProjectScopeSelected
	if request.ProjectScope == v1alpha1.RobotAccountProjectScope_ROBOT_ACCOUNT_PROJECT_SCOPE_ALL {
		scope = robot.ProjectScopeAll
	}
	projects, err := r.checkProjects(ctx, request.ProjectScope, request.Projects)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	rb := &robot.Robot{
		Name:                robotName,
		Description:         request.Description,
		Duration:            int(request.ExpireDays),
		PlatformPermissions: platformPermissions,
		ProjectPermissions:  projectPermissions,
		Enabled:             true,
		TokenHash:           hash,
		Projects:            projects,
		ProjectScope:        scope,
		CreateBy:            user.GetCurrentUserId(ctx),
	}
	if request.ExpireDays > 0 {
		expireAt := time.Now().AddDate(0, 0, int(request.ExpireDays))
		rb.ExpireAt = &expireAt
	}
	if err = r.robotRepo.CreateRobot(ctx, rb); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("robot %s already exists", rb.Name))
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1alpha1.CreateRobotAccountResponse{
		Token: token,
	}, nil
}

func (r *RobotHandler) ListRobotAccounts(ctx context.Context, request *v1alpha1.ListRobotAccountsRequest) (*v1alpha1.ListRobotAccountsResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	page := utils.NewPage(request.Page, request.PageSize)
	robots, total, err := r.robotRepo.ListSystemRobots(ctx, int(page.Page), int(page.PageSize), request.Search)
	if err != nil {
