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
