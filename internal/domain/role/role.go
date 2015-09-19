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

package role

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/samber/lo"
)

// RoleType represents role type
type RoleType int

const (
	PlatformRoleAdmin RoleType = iota + 1
	ProjectRoleAdmin
	ProjectRoleEditor
	ProjectRoleViewer
)

var ProjectWritePermission = []Permission{ModelPush, DatasetPush}

// Role represents a role in the system
type Role struct {
	ID          int            `gorm:"primary_key"`
	Name        string         `gorm:"column:name"`
	Permissions PermissionList `gorm:"column:permissions;type:text"`
	Scope       string         `gorm:"column:scope"`
}

func (Role) TableName() string {
	return "roles"
}

// Permission represents permission type, format is Resource.Action
type Permission string

func (p Permission) String() string {
	return string(p)
}

func (p Permission) IsWildcard() bool {
	return strings.Contains(string(p), "*")
}

// ToRegex converts wildcard permission to regex
func (p Permission) ToRegex() *regexp.Regexp {
	if !p.IsWildcard() {
