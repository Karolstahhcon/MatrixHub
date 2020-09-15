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
	"time"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

type CurrentUserHandler struct {
	userRepo        user.IUserRepo
	accessTokenRepo user.IAccessTokenRepo
	sshKeyRepo      user.ISSHKeyRepo
	sshKeyService   user.ISSHKeyService
}

func (cu *CurrentUserHandler) CreateSSHKey(ctx context.Context, request *v1alpha1.CreateSSHKeyRequest) (*v1alpha1.CreateSSHKeyResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	expireAt, err := cu.parseExpireAt(request.GetExpireAt())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	fp, err := utils.ParseFingerprint(request.GetPublicKey())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	sk := user.SSHKey{
		Name:        request.GetName(),
		UserId:      user.GetCurrentUserId(ctx),
		PublicKey:   request.GetPublicKey(),
		Fingerprint: fp,
		ExpireAt:    expireAt,
	}
	err = cu.sshKeyService.CreateSSHKey(ctx, sk)
	if err != nil {
		if errors.Is(err, user.ErrSSHKeyAlreadyExists) || errors.Is(err, user.ErrSSHKeyExpired) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
