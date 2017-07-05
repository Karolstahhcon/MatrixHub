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

package db

import (
	"fmt"
	"net/url"
	"path/filepath"
)

const (
	DriverMySQL    = "mysql"
	DriverPostgres = "postgres"
	DriverSQLite   = "sqlite"

	MATRIXHUB_DSN_ENV = "MATRIXHUB_DATABASE_DSN"
)

type Config struct {
	Debug                  bool   `yaml:"debug"`
	Driver                 string `yaml:"driver"`
	AccessType             string `yaml:"accessType"`
	DSN                    string `yaml:"dsn"`
	MaxOpenConns           int    `yaml:"maxOpenConns"`
	MaxIdleConns           int    `yaml:"maxIdleConns"`
	ConnMaxLifetimeSeconds int    `yaml:"connMaxLifetimeSeconds"`
	ConnMaxIdleSeconds     int    `yaml:"connMaxIdleSeconds"`

	SQLPath string `yaml:"sqlPath"`
	Migrate bool   `yaml:"migrate"`
}

func (config *Config) applyDefaults() {
	if config.Driver != DriverSQLite {
		return
	}

	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 1
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 1
	}
}

func (config *Config) validate() error {
