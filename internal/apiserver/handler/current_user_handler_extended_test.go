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
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/auth"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	usermocks "github.com/matrixhub-ai/matrixhub/internal/domain/user/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/infra/crypto"
)

func currentUserContext() context.Context {
	return auth.WithIdentity(context.Background(), user.NewUserIdentity(42, "alice"))
}

func TestCurrentUserHandlerGetCurrentUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := usermocks.NewMockIUserRepo(ctrl)
	ctx := currentUserContext()
	repo.EXPECT().GetUser(ctx, 42).Return(&user.User{ID: 42, Username: "alice"}, nil)
	repo.EXPECT().IsUserSysAdmin(ctx, 42).Return(true, nil)

	response, err := (&CurrentUserHandler{userRepo: repo}).GetCurrentUser(ctx, &v1alpha1.GetCurrentUserRequest{})
	require.NoError(t, err)
	require.Equal(t, uint32(42), response.Id)
	require.Equal(t, "alice", response.Username)
	require.True(t, response.IsAdmin)
}

func TestCurrentUserHandlerResetPassword(t *testing.T) {
	hash, err := crypto.HashPassword("old-password")
	require.NoError(t, err)
	ctx := currentUserContext()

	t.Run("updates a changed password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := usermocks.NewMockIUserRepo(ctrl)
		repo.EXPECT().GetUser(ctx, 42).Return(&user.User{ID: 42, Password: hash}, nil)
		repo.EXPECT().UpdateUserPassword(ctx, 42, "new-password").Return(nil)

		response, err := (&CurrentUserHandler{userRepo: repo}).ResetPassword(ctx, &v1alpha1.ResetPasswordRequest{
			OldPassword: "old-password", NewPassword: "new-password",
		})
		require.NoError(t, err)
		require.NotNil(t, response)
	})

	t.Run("rejects an incorrect old password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		repo := usermocks.NewMockIUserRepo(ctrl)
