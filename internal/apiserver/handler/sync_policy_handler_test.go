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

package handler

import (
	"testing"

	v1alpha1 "github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/syncjob"
	syncjobmocks "github.com/matrixhub-ai/matrixhub/internal/domain/syncjob/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/syncpolicy"
	syncpolicymocks "github.com/matrixhub-ai/matrixhub/internal/domain/syncpolicy/mocks"
	"go.uber.org/mock/gomock"
)

// syncJobToProto and syncTaskToProto convert the domain status by a plain
// numeric cast. A domain state with no counterpart in the proto enum therefore
// travels the wire as a bare number — protojson only emits the symbolic name
// for values the enum declares — and every generated client fails to decode the
// whole response. Sync jobs are created in the pending state, so the gap is
// reachable through ListSyncJobs on any freshly created task.
func TestStatusEnumsAreRepresentableInProto(t *testing.T) {
	t.Run("sync job", func(t *testing.T) {
