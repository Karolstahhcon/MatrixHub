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

package preset_test

import (
	"context"
	"errors"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	gitmocks "github.com/matrixhub-ai/matrixhub/internal/domain/git/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
	modelmocks "github.com/matrixhub-ai/matrixhub/internal/domain/model/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/preset"
	presetmocks "github.com/matrixhub-ai/matrixhub/internal/domain/preset/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	projectmocks "github.com/matrixhub-ai/matrixhub/internal/domain/project/mocks"
)

func TestApplyIsBestEffort(t *testing.T) {
	for _, test := range []struct {
		name    string
		enabled bool
		applied bool
		err     error
	}{
		{name: "disabled"},
		{name: "seeded", enabled: true, applied: true},
		{name: "existing instance", enabled: true},
		{name: "database failure", enabled: true, err: errors.New("database unavailable")},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := presetmocks.NewMockIPresetRepo(ctrl)
			manifest := preset.DefaultManifest()
			if test.enabled {
				repo.EXPECT().SeedIfEmpty(gomock.Any(), manifest).Return(test.applied, test.err)
			}
			service := preset.NewService(test.enabled, manifest, repo, nil, nil, nil)
			require.NotPanics(t, func() { service.Apply(context.Background()) })
		})
	}
}

func TestReconcileDisabled(t *testing.T) {
	service := preset.NewService(false, preset.DefaultManifest(), nil, nil, nil, nil)
	require.NotPanics(t, func() { service.Reconcile(context.Background()) })
}

func TestReconcileRepairsOnlyExistingProxyModelsWithoutRemoteAccess(t *testing.T) {
	now := time.Now()
	for _, test := range []struct {
		name        string
		mod         model.Model
		deleted     bool
		local       bool
		repoExists  bool
		createError error
	}{
		{name: "deleted record", deleted: true},
		{name: "local project", local: true},
		{name: "missing repository"},
		{name: "existing repository", repoExists: true},
		{name: "synchronized model still repairs repository", mod: model.Model{SyncedAt: &now}},
		{name: "README is unchanged", mod: model.Model{ReadmeContent: "custom"}},
		{name: "size is unchanged", mod: model.Model{Size: 1}},
		{name: "labels are unchanged", mod: model.Model{Labels: []model.Label{{Name: "custom"}}}},
		{name: "parameter count is unchanged", mod: model.Model{ParameterCount: 1}},
		{name: "repository repair failure", createError: errors.New("disk unavailable")},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			models := modelmocks.NewMockIModelRepo(ctrl)
			projects := projectmocks.NewMockIProjectRepo(ctrl)
			git := gitmocks.NewMockIGitRepo(ctrl)
			manifest := &preset.Manifest{Models: []preset.ModelSpec{{Project: "project", Name: "model"}}}

			if test.deleted {
				models.EXPECT().GetByProjectAndName(gomock.Any(), "project", "model").
					Return(nil, errors.New("record not found"))
			} else {
				models.EXPECT().GetByProjectAndName(gomock.Any(), "project", "model").Return(&test.mod, nil)
				prj := &project.Project{}
				if !test.local {
					id := 3
					prj.RegistryID = &id
					git.EXPECT().RepositoryExists(gomock.Any(), "models", "project", "model").Return(test.repoExists, nil)
					if !test.repoExists {
						git.EXPECT().CreateRepository(gomock.Any(), "models", "project", "model").Return(test.createError)
					}
				}
