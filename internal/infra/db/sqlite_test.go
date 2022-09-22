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
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewSQLiteMigratesDatabase(t *testing.T) {
	dataDir := t.TempDir()
	dsn, err := DefaultSQLiteDSN(dataDir)
	require.NoError(t, err)

	database, err := New(Config{
		Driver:  DriverSQLite,
		DSN:     dsn,
		SQLPath: migrationSQLPath(t),
		Migrate: true,
	})
	require.NoError(t, err)

	sqlDB, err := database.DB()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, sqlDB.Close())
	})

	require.Equal(t, 1, sqlDB.Stats().MaxOpenConnections)
	requirePrivateSQLiteFile(t, filepath.Join(dataDir, "matrixhub.db"))

	var foreignKeys int
	require.NoError(t, database.Raw("PRAGMA foreign_keys").Scan(&foreignKeys).Error)
	require.Equal(t, 1, foreignKeys)

	var journalMode string
	require.NoError(t, database.Raw("PRAGMA journal_mode").Scan(&journalMode).Error)
	require.Equal(t, "wal", journalMode)

	var synchronous int
	require.NoError(t, database.Raw("PRAGMA synchronous").Scan(&synchronous).Error)
	require.Equal(t, 2, synchronous)

	var sqliteVersion string
	require.NoError(t, database.Raw("SELECT sqlite_version()").Scan(&sqliteVersion).Error)
	requireSQLiteVersionAtLeast(t, sqliteVersion, 3, 51, 3)

	var migrationVersion uint
	var dirty bool
	require.NoError(t, database.Table("schema_migrations").
		Select("version, dirty").
		Row().Scan(&migrationVersion, &dirty))
	require.Equal(t, uint(1), migrationVersion)
	require.False(t, dirty)

	var userCount int64
	require.NoError(t, database.Table("users").Where("username = ?", "admin").Count(&userCount).Error)
	require.Equal(t, int64(1), userCount)

	var roleCount int64
	require.NoError(t, database.Table("roles").Count(&roleCount).Error)
	require.Equal(t, int64(4), roleCount)

	var syncedAtColumnCount int64
	require.NoError(t, database.Raw(
		"SELECT COUNT(*) FROM pragma_table_info('models') WHERE name = 'synced_at'",
	).Scan(&syncedAtColumnCount).Error)
	require.Equal(t, int64(1), syncedAtColumnCount)
