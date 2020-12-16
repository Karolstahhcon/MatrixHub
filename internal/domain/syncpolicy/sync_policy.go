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

package syncpolicy

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

// ==================== Custom Types ====================

type SyncPolicyType int

const (
	SyncPolicyTypeUnspecified SyncPolicyType = iota
	SyncPolicyTypePull
	SyncPolicyTypePush
)

func (t SyncPolicyType) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t *SyncPolicyType) Scan(value any) error {
	if value == nil {
		*t = SyncPolicyTypeUnspecified
		return nil
	}
	switch v := value.(type) {
	case int64:
		*t = SyncPolicyType(v)
	case int:
		*t = SyncPolicyType(v)
	case int32:
		*t = SyncPolicyType(v)
	case uint8:
		*t = SyncPolicyType(v)
	default:
		return fmt.Errorf("cannot scan %T into SyncPolicyType", value)
	}
	return nil
}

type TriggerType int

const (
	TriggerTypeUnspecified TriggerType = iota
	TriggerTypeManual
	TriggerTypeScheduled
)

func (t TriggerType) Value() (driver.Value, error) {
	return int64(t), nil
}

func (t *TriggerType) Scan(value any) error {
	if value == nil {
		*t = TriggerTypeUnspecified
		return nil
	}
	switch v := value.(type) {
	case int64:
		*t = TriggerType(v)
	case int:
		*t = TriggerType(v)
	case int32:
		*t = TriggerType(v)
	case uint8:
		*t = TriggerType(v)
	default:
		return fmt.Errorf("cannot scan %T into TriggerType", value)
	}
	return nil
}

// ==================== Entity ====================

type SyncPolicy struct {
	ID                 int            `gorm:"primarykey"`
	Name               string         `gorm:"column:name"`
	Description        string         `gorm:"column:description"`
	PolicyType         SyncPolicyType `gorm:"column:policy_type"`
	TriggerType        TriggerType    `gorm:"column:trigger_type"`
	RegistryID         int            `gorm:"column:registry_id"`
	LocalResourceName  string         `gorm:"column:local_resource_name"`
	LocalProjectName   string         `gorm:"column:local_project_name"`
	RemoteResourceName string         `gorm:"column:remote_resource_name"`
	RemoteProjectName  string         `gorm:"column:remote_project_name"`
	ResourceTypes      string         `gorm:"column:resource_types"`
	Bandwidth          string         `gorm:"column:bandwidth"`
	Cron               string         `gorm:"column:cron"`        // cron expression when TriggerType is scheduled
	LastRunAt          int64          `gorm:"column:last_run_at"` // ms since epoch; last claim / run bookkeeping
