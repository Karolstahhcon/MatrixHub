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
	"fmt"
	"strings"
	"time"

	"github.com/matrixhub-ai/matrixhub/internal/domain/registry"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registrydiscovery"
	"github.com/matrixhub-ai/matrixhub/internal/domain/syncjob"
)

// SyncJobGenerator generates sync tasks and jobs from a sync policy.
// The abstraction decouples job generation logic from the database layer,
// making it testable and extensible for future policy types.
//
//go:generate go tool mockgen -source=generator.go -destination=mocks/generator_mock.go -package=mocks
type SyncJobGenerator interface {
	Generate(ctx context.Context, policy *SyncPolicy) (*SyncTask, []*syncjob.SyncJob, error)
}

// LocalResourceRepository lists resources already stored in MatrixHub.
type LocalResourceRepository interface {
	ListAllPaths(ctx context.Context) ([]string, error)
}

// LocalResourceSource associates a local resource repository with its type.
type LocalResourceSource struct {
	ResourceType string
	Repo         LocalResourceRepository
}

type syncJobGenerator struct {
	registryRepo  registry.IRegistryRepo
	discoveries   map[string]registrydiscovery.Discovery
	localResource map[string]LocalResourceRepository
}

// NewSyncJobGenerator creates a new SyncJobGenerator instance.
func NewSyncJobGenerator(registryRepo registry.IRegistryRepo, discoveries map[string]registrydiscovery.Discovery, localSources ...LocalResourceSource) SyncJobGenerator {
	localResource := make(map[string]LocalResourceRepository, len(localSources))
	for _, source := range localSources {
		if source.Repo != nil {
			localResource[source.ResourceType] = source.Repo
		}
	}

	return &syncJobGenerator{
		registryRepo:  registryRepo,
		discoveries:   discoveries,
		localResource: localResource,
	}
}

func (g *syncJobGenerator) Generate(ctx context.Context, policy *SyncPolicy) (*SyncTask, []*syncjob.SyncJob, error) {
	task := g.buildTask(policy)

	var jobs []*syncjob.SyncJob
	var err error

	switch {
	case policy.IsPullBase() && policy.HasWildcardResourceName():
		jobs, err = g.buildJobsFromDiscovery(ctx, policy)
	case policy.IsPushBase() && policy.HasWildcardResourceName():
		jobs, err = g.buildJobsFromLocalDiscovery(ctx, policy)
	default:
		jobs = g.buildJobsFromStatic(policy)
	}
	if err != nil {
		return nil, nil, err
	}

	task.TotalItems = len(jobs)
	return task, jobs, nil
}

func (g *syncJobGenerator) buildJobsFromLocalDiscovery(ctx context.Context, policy *SyncPolicy) ([]*syncjob.SyncJob, error) {
	if len(g.localResource) == 0 {
		return nil, fmt.Errorf("local resource discovery is not configured")
	}

	resourceTypes := g.parseResourceTypes(policy.ResourceTypes)
	var jobs []*syncjob.SyncJob
	for _, rt := range resourceTypes {
		repo, ok := g.localResource[rt]
		if !ok {
			return nil, fmt.Errorf("no local resource repository for type: %s", rt)
		}

		paths, err := repo.ListAllPaths(ctx)
		if err != nil {
			return nil, fmt.Errorf("list local resources (type=%s): %w", rt, err)
		}
		for _, path := range paths {
			project, name := splitResourcePath(path)
			if !matchesWildcardResource(project, name, policy.LocalProjectName, policy.LocalResourceName) {
				continue
			}
			jobs = append(jobs, g.buildPushJobFromPath(policy, project, name, rt))
		}
	}
	return jobs, nil
}

func (g *syncJobGenerator) buildPushJobFromPath(policy *SyncPolicy, project, name, resourceType string) *syncjob.SyncJob {
	return &syncjob.SyncJob{
		RemoteRegistryID:   policy.RegistryID,
		RemoteProjectName:  policy.RemoteProjectName,
		RemoteResourceName: name,
		ProjectName:        project,
		ResourceName:       name,
		ResourceType:       resourceType,
		SyncType:           "push",
		Status:             syncjob.SyncJobStatusRunning,
		CompletePercents:   0,
	}
}

func splitResourcePath(path string) (project, name string) {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return "", path
	}
