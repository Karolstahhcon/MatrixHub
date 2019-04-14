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
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	modelv1alpha1 "github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/authz"
	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

type ModelHandler struct {
	ms           model.IModelService
	authzService authz.IAuthzService
}

func NewModelHandler(modelService model.IModelService, authzService authz.IAuthzService) IHandler {
	handler := &ModelHandler{
		ms:           modelService,
		authzService: authzService,
	}
	return handler
}

// formatSize converts bytes to string without unit conversion
func formatSize(bytes int64) string {
	return fmt.Sprintf("%d", bytes)
}

// formatParameterCount converts parameter count to string without unit conversion
func formatParameterCount(count int64) string {
	return fmt.Sprintf("%d", count)
}

func labelCategoryToProto(category string) modelv1alpha1.Category {
	switch category {
	case "library":
		return modelv1alpha1.Category_LIBRARY
	case "license":
		return modelv1alpha1.Category_LICENSE
	case "language":
		return modelv1alpha1.Category_LANGUAGE
	case "other":
		return modelv1alpha1.Category_OTHER
	case "task":
		return modelv1alpha1.Category_TASK
	default:
		return modelv1alpha1.Category_OTHER
	}
}

// modelToProto converts domain Model to proto Model
func modelToProto(m *model.Model) *modelv1alpha1.Model {
	labels := make([]*modelv1alpha1.Label, len(m.Labels))
	for i, l := range m.Labels {
		labels[i] = &modelv1alpha1.Label{
			Id:        int32(l.ID),
			Name:      l.Name,
			Category:  labelCategoryToProto(l.Category),
			CreatedAt: l.CreatedAt.Format(time.RFC3339),
			UpdatedAt: l.UpdatedAt.Format(time.RFC3339),
		}
	}
	return &modelv1alpha1.Model{
		Id:            int32(m.ID),
		Name:          m.Name,
		Nickname:      "", // Not implemented yet
		DefaultBranch: m.DefaultBranch,
		SyncedAt:      formatOptionalTime(m.SyncedAt),
		CreatedAt:     m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     m.UpdatedAt.Format(time.RFC3339),
		CloneUrls: &modelv1alpha1.CloneUrls{
			SshUrl:  "",
			HttpUrl: "",
		},
		Labels:         labels,
		Project:        m.ProjectName,
		ReadmeContent:  m.ReadmeContent,
		Size:           formatSize(m.Size),
		ParameterCount: formatParameterCount(m.ParameterCount),
		Popular:        m.IsPopular,
	}
}

func formatOptionalTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format(time.RFC3339)
}

// revisionToProto converts domain git.Revision to proto Revision
func revisionToProto(r *git.Revision) *modelv1alpha1.Revision {
	return &modelv1alpha1.Revision{
		Name: r.Name,
	}
}

// revisionsToProto converts domain git.Revisions to proto Revisions
func revisionsToProto(revisions *git.Revisions) *modelv1alpha1.Revisions {
	branches := make([]*modelv1alpha1.Revision, len(revisions.Branches))
	for i, b := range revisions.Branches {
		branches[i] = revisionToProto(b)
	}

	tags := make([]*modelv1alpha1.Revision, len(revisions.Tags))
	for i, t := range revisions.Tags {
		tags[i] = revisionToProto(t)
	}

