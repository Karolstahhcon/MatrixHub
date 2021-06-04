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

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/matrixhub-ai/matrixhub/internal/domain/auth"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
)

// methodPermissions maps GRPC methods to required permissions
var methodPermissions = map[string]role.Permission{
	// User management
	"/matrixhub.v1alpha1.Users/GetUser":           role.UserGet,
	"/matrixhub.v1alpha1.Users/CreateUser":        role.UserCreate,
	"/matrixhub.v1alpha1.Users/SetUserSysAdmin":   role.UserAuthorize,
	"/matrixhub.v1alpha1.Users/DeleteUser":        role.UserDelete,
	"/matrixhub.v1alpha1.Users/ResetUserPassword": role.UserResetPassword,

	// Registry management
	"/matrixhub.v1alpha1.Registries/ListRegistries": role.RegistryGet,
	"/matrixhub.v1alpha1.Registries/GetRegistry":    role.RegistryGet,
	"/matrixhub.v1alpha1.Registries/PingRegistry":   role.RegistryGet,
	"/matrixhub.v1alpha1.Registries/CreateRegistry": role.RegistryCreate,
	"/matrixhub.v1alpha1.Registries/UpdateRegistry": role.RegistryUpdate,
	"/matrixhub.v1alpha1.Registries/DeleteRegistry": role.RegistryDelete,

	// Sync policy management
	"/matrixhub.v1alpha1.SyncPolicy/ListSyncPolicies": role.SyncGet,
	"/matrixhub.v1alpha1.SyncPolicy/GetSyncPolicy":    role.SyncGet,
