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

package middleware

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/matrixhub-ai/matrixhub/internal/domain/auth"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
)

func TestAuthzInterceptorCleanupPermissions(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		permission role.Permission
	}{
		{
			name:       "preview cleanup requires cleanup get",
			method:     "/matrixhub.v1alpha1.Cleanup/PreviewCleanup",
			permission: role.CleanupGet,
		},
		{
			name:       "storage stats requires cleanup get",
			method:     "/matrixhub.v1alpha1.Cleanup/GetStorageStats",
			permission: role.CleanupGet,
		},
		{
			name:       "execute cleanup requires cleanup execute",
			method:     "/matrixhub.v1alpha1.Cleanup/ExecuteCleanup",
			permission: role.CleanupExecute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := auth.WithIdentity(context.Background(), user.NewUserIdentity(1, "admin"))
			var gotPerm role.Permission
			interceptor := AuthzInterceptor(func(ctx context.Context, perm role.Permission) (bool, error) {
				gotPerm = perm
