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

package apiserver

import (
	"context"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	gitmocks "github.com/matrixhub-ai/matrixhub/internal/domain/git/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
	modelmocks "github.com/matrixhub-ai/matrixhub/internal/domain/model/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/preset"
	presetmocks "github.com/matrixhub-ai/matrixhub/internal/domain/preset/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	projectmocks "github.com/matrixhub-ai/matrixhub/internal/domain/project/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/infra/config"
	"github.com/matrixhub-ai/matrixhub/internal/repo"
)

func TestInitGitStorageDoesNotRegisterMetadataPostReceiveHook(t *testing.T) {
	server := &APIServer{
		config: &config.Config{DataDir: t.TempDir()},
	}
	server.initMirrorHooks()
	server.initGitStorage()

	postReceiveHook := reflect.ValueOf(server.gitStorage.sharedMirror).
		Elem().FieldByName("postReceiveHookFunc")
	if !postReceiveHook.IsNil() {
		t.Fatal("mirror must not register the metadata post-receive hook")
	}
}

func TestInitPresetsSkipsDisabledPresets(t *testing.T) {
	ctrl := gomock.NewController(t)
	disabled := false
	server := &APIServer{
		config: &config.Config{Presets: config.PresetsConfig{Enabled: &disabled}},
		repos: &repo.Repos{
			Preset:  presetmocks.NewMockIPresetRepo(ctrl),
			Project: projectmocks.NewMockIProjectRepo(ctrl),
			Model:   modelmocks.NewMockIModelRepo(ctrl),
			Git:     gitmocks.NewMockIGitRepo(ctrl),
		},
	}

	// No expectations are registered, so any repository access fails the test.
	server.initPresets()
}

func TestInitPresetsWaitsForRepositoryReconciliation(t *testing.T) {
