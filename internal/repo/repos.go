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

package repo

import (
	"github.com/matrixhub-ai/hfd/pkg/mirror"
	gitstorage "github.com/matrixhub-ai/hfd/pkg/storage"
	"gorm.io/gorm"

	"github.com/matrixhub-ai/matrixhub/internal/domain/authz"
	"github.com/matrixhub-ai/matrixhub/internal/domain/dataset"
	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
	"github.com/matrixhub-ai/matrixhub/internal/domain/preset"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registry"
	"github.com/matrixhub-ai/matrixhub/internal/domain/robot"
	"github.com/matrixhub-ai/matrixhub/internal/domain/syncjob"
	"github.com/matrixhub-ai/matrixhub/internal/domain/syncpolicy"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/config"
	"github.com/matrixhub-ai/matrixhub/internal/infra/db"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
)

type Repos struct {
	DB          *gorm.DB
	GitStorage  *gitstorage.Storage
	GitMirror   *mirror.Mirror
	Project     project.IProjectRepo
	User        user.IUserRepo
	Registry    registry.IRegistryRepo
	Model       model.IModelRepo
	Label       model.ILabelRepo
	Git         git.IGitRepo
	Dataset     dataset.IDatasetRepo
	Session     user.ISessionRepo
	AccessToken user.IAccessTokenRepo
	SSHKey      user.ISSHKeyRepo
	SyncPolicy  syncpolicy.ISyncPolicyRepo
	SyncTask    syncpolicy.ISyncTaskRepo
	SyncJob     syncjob.ISyncJobRepo
	Authz       authz.IAuthzProjectRepo
	Robot       robot.IRobotRepo
