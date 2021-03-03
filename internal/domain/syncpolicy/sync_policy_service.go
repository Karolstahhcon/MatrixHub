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
func (sps *SyncPolicyService) SetOnTaskCreated(fn func(taskID int) bool) {
	sps.onTaskCreated = fn
}

func (sps *SyncPolicyService) notifyTaskCreated(taskID int) {
	if sps.onTaskCreated != nil && !sps.onTaskCreated(taskID) {
		log.Warnw("sync task trigger was not queued; periodic polling will recover it", "taskID", taskID)
	}
}

func (sps *SyncPolicyService) CreateSyncPolicy(ctx context.Context, param *SyncPolicy) error {
	if err := param.ApplyScheduleNextRun(time.Now()); err != nil {
		return err
	}
	return sps.syncPolicyRepo.CreateSyncPolicy(ctx, param)
}

func (sps *SyncPolicyService) UpdateSyncPolicy(ctx context.Context, param *SyncPolicy) error {
	if err := param.ApplyScheduleNextRun(time.Now()); err != nil {
		return err
	}
	return sps.syncPolicyRepo.UpdateSyncPolicy(ctx, param)
}

func (sps *SyncPolicyService) DeleteSyncPolicy(ctx context.Context, id int) error {
	return sps.syncPolicyRepo.DeleteSyncPolicy(ctx, id)
}

func (sps *SyncPolicyService) ListSyncPolicies(ctx context.Context, page, pageSize int, search string) ([]*SyncPolicy, int64, error) {
	return sps.syncPolicyRepo.ListSyncPolicies(ctx, page, pageSize, search)
}

func (sps *SyncPolicyService) GetSyncTask(ctx context.Context, id int) (*SyncTask, error) {
	return sps.syncTaskRepo.GetSyncTask(ctx, id)
}

func (sps *SyncPolicyService) CreateSyncTask(ctx context.Context, syncTask *SyncTask) (*SyncTask, error) {
	task, err := sps.syncTaskRepo.CreateSyncTask(ctx, syncTask)
	if err == nil {
		sps.notifyTaskCreated(task.ID)
	}
	return task, err
}

func (sps *SyncPolicyService) UpdateSyncTask(ctx context.Context, syncTask *SyncTask) error {
	return sps.syncTaskRepo.UpdateSyncTask(ctx, syncTask)
}

func (sps *SyncPolicyService) ListSyncTasksByPolicyID(ctx context.Context, policyID int, page, pageSize int, status SyncTaskStatus) ([]*SyncTask, int64, error) {
	return sps.syncTaskRepo.ListSyncTasksByPolicyID(ctx, policyID, page, pageSize, status)
}

func (sps *SyncPolicyService) CreateSyncTaskAndSyncJobs(ctx context.Context, policy *SyncPolicy) error {
	task, jobs, err := sps.jobGenerator.Generate(ctx, policy)
	if err != nil {
		return err
	}
	task, err = sps.syncTaskRepo.CreateSyncTask(ctx, task)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		job.SyncTaskID = task.ID
	}
	if err := sps.syncJobService.CreateSyncJobs(ctx, jobs); err != nil {
		log.Infow("CreateSyncJobs failed", "error", err)
		task.Status = SyncTaskStatusFailed
		task.CompletedTimestamp = time.Now().Unix()
		if updateErr := sps.syncTaskRepo.UpdateSyncTask(ctx, task); updateErr != nil {
			log.Errorw("failed to mark task as failed after job creation error", "taskID", task.ID, "error", updateErr)
		}
		return err
	}
	return nil
}

