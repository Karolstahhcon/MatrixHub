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

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1alpha1 "github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/authz"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

type UserHandler struct {
	userRepo     user.IUserRepo
	authzService authz.IAuthzService
}

func (u *UserHandler) SetUserSysAdmin(ctx context.Context, request *userv1alpha1.SetUserSysAdminRequest) (*userv1alpha1.SetUserSysAdminResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if allowed, err := u.authzService.VerifyPlatformPermission(ctx, role.UserAuthorize); err != nil || !allowed {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	userID := int(request.Id)
	if err := u.userRepo.SetUserSysAdmin(ctx, userID, request.SysadminFlag); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userv1alpha1.SetUserSysAdminResponse{}, nil
}

func (u *UserHandler) ResetUserPassword(ctx context.Context, request *userv1alpha1.ResetUserPasswordRequest) (*userv1alpha1.ResetUserPasswordResponse, error) {
