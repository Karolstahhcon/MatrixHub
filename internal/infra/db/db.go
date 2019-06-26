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
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	gmysql "gorm.io/driver/mysql"
	gpostgres "gorm.io/driver/postgres"
	gsqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const databasePingTimeout = 10 * time.Second

func New(config Config) (*gorm.DB, error) {
	config.applyDefaults()
	if err := config.validate(); err != nil {
		return nil, err
	}

	var sqlitePath string
	var sqliteMemory bool
	var sqliteDSN string
	if config.Driver == DriverSQLite {
		var err error
		_, fileBacked, _, err := sqliteFilePath(config.DSN)
		if err != nil {
			return nil, err
		}
		sqliteMemory = !fileBacked
		sqlitePath, _, err = prepareSQLiteFile(config.DSN)
		if err != nil {
			return nil, err
		}
		sqliteDSN, err = sqliteDSNWithForeignKeys(config.DSN)
		if err != nil {
			return nil, err
		}
	}

	var dialector gorm.Dialector
	switch config.Driver {
	case DriverMySQL:
		dialector = gmysql.Open(config.DSN)
	case DriverPostgres:
		dialector = gpostgres.Open(config.DSN)
	case DriverSQLite:
		dialector = gsqlite.Open(sqliteDSN)
	}

	database, err := gorm.Open(dialector, &gorm.Config{
