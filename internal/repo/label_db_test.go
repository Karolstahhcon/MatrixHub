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

package repo

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
)

const (
	labelInsertSQL = "INSERT INTO `labels` (`name`,`category`,`scope`,`created_at`,`updated_at`) VALUES (?,?,?,?,?) " +
		"ON DUPLICATE KEY UPDATE `id`=`id`"
	labelLookupSQL = "SELECT * FROM `labels` WHERE name = ? AND category = ? AND scope = ? ORDER BY `labels`.`id` LIMIT ?"
)

func newLabelDBForTest(t *testing.T) (*labelDB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock := newMockGormDB(t)
	return NewLabelDB(db).(*labelDB), mock
}

func labelRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "name", "category", "scope", "created_at", "updated_at"})
}

func TestLabelDB_ListByCategoryAndScope(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("database unavailable")

	for _, test := range []struct {
		name    string
		rows    *sqlmock.Rows
		err     error
		want    []*model.Label
		wantErr error
	}{
		{
			name: "maps matching labels",
			rows: labelRows().
				AddRow(11, "text-generation", "task", "model", repoTestCreatedAt, repoTestUpdatedAt).
				AddRow(12, "image-classification", "task", "model", repoTestCreatedAt, repoTestUpdatedAt),
			want: []*model.Label{
				{ID: 11, Name: "text-generation", Category: "task", Scope: "model", CreatedAt: repoTestCreatedAt, UpdatedAt: repoTestUpdatedAt},
				{ID: 12, Name: "image-classification", Category: "task", Scope: "model", CreatedAt: repoTestCreatedAt, UpdatedAt: repoTestUpdatedAt},
			},
		},
		{
			name:    "query error",
			err:     dbErr,
			wantErr: dbErr,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			repo, mock := newLabelDBForTest(t)
			query := mock.ExpectQuery(exactSQL("SELECT * FROM `labels` WHERE category = ? AND scope = ?")).
				WithArgs("task", "model")
			if test.err != nil {
				query.WillReturnError(test.err)
			} else {
				query.WillReturnRows(test.rows)
