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

package syncpolicy

import (
	"context"
	"time"

	"github.com/matrixhub-ai/matrixhub/internal/domain/job"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	"github.com/matrixhub-ai/matrixhub/internal/domain/syncjob"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
)

//go:generate go tool mockgen -source=sync_policy_service.go -destination=mocks/sync_policy_service_mock.go -package=mocks
type ISyncPolicyService interface {
	GetSyncPolicy(ctx context.Context, id int) (*SyncPolicy, error)
	CreateSyncPolicy(ctx context.Context, param *SyncPolicy) error
	UpdateSyncPolicy(ctx context.Context, param *SyncPolicy) error
	DeleteSyncPolicy(ctx context.Context, id int) error
	ListSyncPolicies(ctx context.Context, page, pageSize int, search string) ([]*SyncPolicy, int64, error)
	GetSyncTask(ctx context.Context, id int) (*SyncTask, error)
	CreateSyncTask(ctx context.Context, param *SyncTask) (*SyncTask, error)
	UpdateSyncTask(ctx context.Context, param *SyncTask) error
	ListSyncTasksByPolicyID(ctx context.Context, policyID int, page, pageSize int, status SyncTaskStatus) ([]*SyncTask, int64, error)
	CreateSyncTaskAndSyncJobs(ctx context.Context, policy *SyncPolicy) error
	CreateSyncTaskAsync(ctx context.Context, policy *SyncPolicy) (*SyncTask, error)

	ClaimDueSyncPolicies(ctx context.Context, nowMs int64) ([]job.DueJob, error)
	// CreatePendingSyncTask inserts a sync_tasks row only; sync_task_processor runs git work later.
	CreatePendingSyncTask(ctx context.Context, policyID int, triggerType int) error

	// ClaimPendingSyncTasks selects pending tasks and CAS-claims them for execution.
	ClaimPendingSyncTasks(ctx context.Context, nowMs int64) ([]job.DueJob, error)
	// ExecuteSyncTask generates jobs from a claimed task and creates pending sync_jobs.
	ExecuteSyncTask(ctx context.Context, taskID int) error

	ReportTaskStatus(ctx context.Context, taskID int) error
}

type SyncPolicyService struct {
	syncPolicyRepo ISyncPolicyRepo
	syncTaskRepo   ISyncTaskRepo
	syncJobService syncjob.ISyncJobService
	projectRepo    project.IProjectRepo
	jobGenerator   SyncJobGenerator
	onTaskCreated  func(taskID int) bool
}

func NewSyncPolicyService(sprepo ISyncPolicyRepo, strepo ISyncTaskRepo, sjservice syncjob.ISyncJobService, prepo project.IProjectRepo, jobGenerator SyncJobGenerator) ISyncPolicyService {
	return &SyncPolicyService{
		syncPolicyRepo: sprepo,
		syncTaskRepo:   strepo,
		syncJobService: sjservice,
		projectRepo:    prepo,
		jobGenerator:   jobGenerator,
	}
}

func (sps *SyncPolicyService) GetSyncPolicy(ctx context.Context, id int) (*SyncPolicy, error) {
	return sps.syncPolicyRepo.GetSyncPolicy(ctx, id)
}

// SetOnTaskCreated registers a best-effort wakeup called after a task is persisted.
