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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	hfdssh "github.com/matrixhub-ai/hfd/pkg/ssh"
	"github.com/spf13/viper"

	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/db"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
)

type Config struct {
	Debug         bool             `yaml:"debug"`
	Log           log.Config       `yaml:"log"`
	APIServer     *APIServerConfig `yaml:"apiServer" validate:"required"`
	UI            UIConfig         `yaml:"ui"`
	MigrationPath string           `yaml:"migrationPath" validate:"required"`

	DataDir string `yaml:"dataDir" validate:"required"`

	Database db.Config     `yaml:"database" validate:"required"`
	Session  SessionConfig `yaml:"session"`
	Presets  PresetsConfig `yaml:"presets"`

	// JobServer runs delayed sync (and future kinds). If nil, jobserver is disabled.
	JobServer *JobServerConfig `yaml:"jobServer"`
}

type PresetsConfig struct {
	Enabled *bool `yaml:"enabled"`
}

func (c PresetsConfig) IsEnabled() bool {
	return c.Enabled == nil || *c.Enabled
}

// JobServerConfig is the top-level jobserver configuration (YAML key `jobServer`).
type JobServerConfig struct {
	Enabled       bool             `yaml:"enabled"`
	ShutdownGrace time.Duration    `yaml:"shutdownGrace"`
	SyncPolicy    SyncPolicyConfig `yaml:"syncPolicy"`
	SyncTask      SyncTaskConfig   `yaml:"syncTask"`
	SyncJob       SyncJobConfig    `yaml:"syncJob"`
	LogDir        string           `yaml:"logDir"`
}

// SyncPolicyConfig holds per-processor tuning for the sync-policy delayed-job poller.
type SyncPolicyConfig struct {
	PollInterval    time.Duration `yaml:"pollInterval"`
	MaxConcurrent   int           `yaml:"maxConcurrent"`
	TaskMaxDuration time.Duration `yaml:"taskMaxDuration"`
}

// SyncTaskConfig holds tuning for the sync-task processor.
type SyncTaskConfig struct {
	PollInterval    time.Duration `yaml:"pollInterval"`
	MaxConcurrent   int           `yaml:"maxConcurrent"`
	TaskMaxDuration time.Duration `yaml:"taskMaxDuration"`
}

// SyncJobConfig holds tuning for the sync-job processor.
type SyncJobConfig struct {
	PollInterval    time.Duration `yaml:"pollInterval"`
	MaxConcurrent   int           `yaml:"maxConcurrent"`
	TaskMaxDuration time.Duration `yaml:"taskMaxDuration"`
