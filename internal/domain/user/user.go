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

	"github.com/matrixhub-ai/matrixhub/internal/infra/crypto"
)

type User struct {
	ID             int `gorm:"primary_key"`
	Username       string
	Password       string
	Email          string
	IsAdmin        bool `gorm:"-"` // Platform admin flag (not read from DB directly, determined via members_roles_projects table)
	SessionVersion int
	CreatedAt      time.Time
