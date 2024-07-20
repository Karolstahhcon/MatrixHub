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
