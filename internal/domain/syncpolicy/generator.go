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

