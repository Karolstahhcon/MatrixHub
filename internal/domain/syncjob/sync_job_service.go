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

package syncjob

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/matrixhub-ai/matrixhub/internal/domain/dataset"
	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
	"github.com/matrixhub-ai/matrixhub/internal/domain/job"
	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registry"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
)

//go:generate go tool mockgen -source=sync_job_service.go -destination=mocks/sync_job_service_mock.go -package=mocks
type ISyncJobService interface {
	GetSyncJob(ctx context.Context, id int) (*SyncJob, error)
	CreateSyncJob(ctx context.Context, param *SyncJob) error
	CreateSyncJobs(ctx context.Context, params []*SyncJob) error
	UpdateSyncJob(ctx context.Context, param *SyncJob) error
	ExecuteSyncJob(ctx context.Context, param *SyncJob) error
	ListSyncJobsByTaskID(ctx context.Context, taskID int, page, pageSize int, status SyncJobStatus, resourceType string) ([]*SyncJob, int64, error)

	ClaimPendingSyncJobs(ctx context.Context, nowMs int64) ([]job.DueJob, error)
	ExecuteSyncJobWithLog(ctx context.Context, jobID int) error

	SetOnJobDone(fn func(ctx context.Context, taskID int) error)
}

// resourceTypeDataset is the sync job resource type for datasets; every other
// value is handled as a model.
const resourceTypeDataset = "dataset"

// MetadataSyncer refreshes the database metadata (README, size, labels) of a
// resource from its on-disk Git repository.
type MetadataSyncer interface {
	SyncMetadata(ctx context.Context, project, name string) error
}

type SyncJobService struct {
	syncJobRepo   ISyncJobRepo
	registryRepo  registry.IRegistryRepo
	projectRepo   project.IProjectRepo
	modelRepo     model.IModelRepo
	datasetRepo   dataset.IDatasetRepo
	gitRepo       git.IGitRepo
	modelMetadata MetadataSyncer
	datasetMeta   MetadataSyncer
	logStore      LogStore
	onJobDone     func(ctx context.Context, taskID int) error
	onJobCreated  func(jobID int) bool
}

type LogStore interface {
	Writer(jobID int) (io.WriteCloser, error)
	Reader(jobID int) (io.ReadCloser, error)
}

func NewSyncJobService(srepo ISyncJobRepo, rrepo registry.IRegistryRepo, prepo project.IProjectRepo, mrepo model.IModelRepo, dsrepo dataset.IDatasetRepo, grepo git.IGitRepo, modelMeta, datasetMeta MetadataSyncer, logStore LogStore) ISyncJobService {
	return &SyncJobService{
		syncJobRepo:   srepo,
		registryRepo:  rrepo,
		projectRepo:   prepo,
		modelRepo:     mrepo,
		datasetRepo:   dsrepo,
		gitRepo:       grepo,
		modelMetadata: modelMeta,
		datasetMeta:   datasetMeta,
		logStore:      logStore,
	}
}

// SetOnJobDone sets the callback invoked after a job finishes (used for task status aggregation).
func (sjs *SyncJobService) SetOnJobDone(fn func(ctx context.Context, taskID int) error) {
	sjs.onJobDone = fn
}

// SetOnJobCreated registers a best-effort wakeup called after a job is persisted.
func (sjs *SyncJobService) SetOnJobCreated(fn func(jobID int) bool) {
	sjs.onJobCreated = fn
}

func (sjs *SyncJobService) GetSyncJob(ctx context.Context, id int) (*SyncJob, error) {
	return sjs.syncJobRepo.GetSyncJob(ctx, id)
}

func (sjs *SyncJobService) CreateSyncJob(ctx context.Context, syncJob *SyncJob) error {
	return sjs.CreateSyncJobs(ctx, []*SyncJob{syncJob})
}

// CreateSyncJobs persists the full batch before triggering any job, so task
// status aggregation cannot observe a partially-created set of jobs.
func (sjs *SyncJobService) CreateSyncJobs(ctx context.Context, syncJobs []*SyncJob) error {
	if err := sjs.syncJobRepo.CreateSyncJobs(ctx, syncJobs); err != nil {
		return err
	}
	for _, syncJob := range syncJobs {
		if sjs.onJobCreated != nil && !sjs.onJobCreated(syncJob.ID) {
			log.Warnw("sync job trigger was not queued; periodic polling will recover it", "jobID", syncJob.ID)
		}
	}
	return nil
}

func (sjs *SyncJobService) UpdateSyncJob(ctx context.Context, syncJob *SyncJob) error {
	return sjs.syncJobRepo.UpdateSyncJob(ctx, syncJob)
}

func (sjs *SyncJobService) ListSyncJobsByTaskID(ctx context.Context, taskID int, page, pageSize int, status SyncJobStatus, resourceType string) ([]*SyncJob, int64, error) {
	return sjs.syncJobRepo.ListSyncJobsByTaskID(ctx, taskID, page, pageSize, status, resourceType)
}

func (sjs *SyncJobService) ExecuteSyncJob(ctx context.Context, syncJob *SyncJob) error {
	defer func() {
		syncJob.CompletedTimestamp = time.Now().Unix()
		if err := sjs.syncJobRepo.UpdateSyncJob(ctx, syncJob); err != nil {
			fmt.Printf("update sync job failed: %v\n", err)
		}
	}()

	return sjs.executeGitJob(ctx, syncJob, nil)
}

const claimJobBatchLimit = 32

// ClaimPendingSyncJobs selects pending jobs and CAS-claims them for execution.
func (sjs *SyncJobService) ClaimPendingSyncJobs(ctx context.Context, _ int64) ([]job.DueJob, error) {
	candidates, err := sjs.syncJobRepo.SelectPendingJobs(ctx, claimJobBatchLimit)
	if err != nil {
		return nil, err
	}
	var out []job.DueJob
	for _, j := range candidates {
		claimed, err := sjs.syncJobRepo.UpdateJobStatusCAS(ctx, j.ID, SyncJobStatusPending, SyncJobStatusRunning)
