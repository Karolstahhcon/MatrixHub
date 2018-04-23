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
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := u.userRepo.UpdateUserPassword(ctx, int(request.Id), request.Password); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &userv1alpha1.ResetUserPasswordResponse{}, nil
}

func (u *UserHandler) CreateUser(ctx context.Context, request *userv1alpha1.CreateUserRequest) (*userv1alpha1.CreateUserResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if _, err := u.userRepo.GetUserByName(ctx, request.Username); err == nil {
		return nil, status.Error(codes.InvalidArgument, "user already exists")
	}
	user := &user.User{
		Username: request.Username,
		Password: request.Password,
	}
	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if request.IsAdmin {
		if err := u.userRepo.SetUserSysAdmin(ctx, user.ID, true); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &userv1alpha1.CreateUserResponse{}, nil
}

func (u *UserHandler) GetUser(ctx context.Context, request *userv1alpha1.GetUserRequest) (*userv1alpha1.GetUserResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	user, err := u.userRepo.GetUser(ctx, int(request.Id))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &userv1alpha1.GetUserResponse{
		Id:        uint32(user.ID),
		Username:  user.Username,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}, nil
}

func (u *UserHandler) DeleteUser(ctx context.Context, request *userv1alpha1.DeleteUserRequest) (*userv1alpha1.DeleteUserResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := u.userRepo.DeleteUser(ctx, int(request.Id)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &userv1alpha1.DeleteUserResponse{}, nil
