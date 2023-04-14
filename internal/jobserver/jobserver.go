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

package jobserver

import (
	"context"
	"sync"
	"time"

	"github.com/matrixhub-ai/matrixhub/internal/domain/syncjob"
	"github.com/matrixhub-ai/matrixhub/internal/domain/syncpolicy"
	"github.com/matrixhub-ai/matrixhub/internal/infra/config"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
	"github.com/matrixhub-ai/matrixhub/internal/jobserver/canceller"
	"github.com/matrixhub-ai/matrixhub/internal/jobserver/logstore"
	"github.com/matrixhub-ai/matrixhub/internal/jobserver/processor"
)

// JobServer owns one goroutine per processor (syncPolicy, ...). Processors do not share state with each other.
type JobServer struct {
	cfg             *config.JobServerConfig
	processors      []processor.Adapter
	syncTaskTrigger processor.Adapter
	syncJobTrigger  processor.Adapter
}

// New builds a JobServer from config and domain services (cfg is copied; defaults applied without mutating the caller's struct).
func New(cfg *config.JobServerConfig, syncSvc syncpolicy.ISyncPolicyService, syncJobSvc syncjob.ISyncJobService, ls logstore.LogStore, canc canceller.Canceller) *JobServer {
