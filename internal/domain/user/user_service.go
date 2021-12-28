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

package user

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/matrixhub-ai/matrixhub/internal/domain/auth"
)

var (
	InvalidUsernameOrPassword = errors.New("invalid username or password")
)

//go:generate go tool mockgen -source=user_service.go -destination=mocks/user_service_mock.go -package=mocks
type IUserService interface {
	LoginUser(ctx context.Context, username, password string, rememberMe bool) error
	LogoutUser(ctx context.Context) error
	GetCurrentUser(ctx context.Context) (*User, error)
}

type UserService struct {
	sessionRepo ISessionRepo
	userRepo    IUserRepo
}

func (us UserService) GetCurrentUser(ctx context.Context) (*User, error) {
	username := us.sessionRepo.Manager().GetString(ctx, UsernameCtxKey)
	if username == "" {
