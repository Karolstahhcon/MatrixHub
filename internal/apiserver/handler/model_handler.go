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

	return &modelv1alpha1.Revisions{
		Branches: branches,
		Tags:     tags,
	}
}

// commitToProto converts domain git.Commit to proto Commit
func commitToProto(c *git.Commit) *modelv1alpha1.Commit {
	return &modelv1alpha1.Commit{
		Id:             c.ID,
		Message:        c.Message,
		AuthorName:     c.AuthorName,
		AuthorEmail:    c.AuthorEmail,
		AuthorDate:     c.AuthorDate.Format(time.RFC3339),
		CommitterName:  c.CommitterName,
		CommitterEmail: c.CommitterEmail,
		CommitterDate:  c.CommitterDate.Format(time.RFC3339),
		Diff:           c.Diff,
		CreatedAt:      c.CreatedAt.Format(time.RFC3339),
	}
}

// treeEntryToProtoFile converts domain git.TreeEntry to proto File
func treeEntryToProtoFile(entry *git.TreeEntry) *modelv1alpha1.File {
	var protoCommit *modelv1alpha1.Commit
	if entry.Commit != nil {
		protoCommit = commitToProto(entry.Commit)
	}

	return &modelv1alpha1.File{
		Name:   entry.Name,
		Type:   modelv1alpha1.FileType(entry.Type),
		Path:   entry.Path,
		Size:   entry.Size,
		Lfs:    entry.IsLFS,
		Sha256: entry.Hash,
		Commit: protoCommit,
		Url:    entry.URL,
	}
}

func (mh *ModelHandler) RegisterToServer(opt *ServerOptions) {
	modelv1alpha1.RegisterModelsServer(opt.GRPCServer, mh)
	if err := modelv1alpha1.RegisterModelsHandlerFromEndpoint(context.Background(), opt.GatewayMux, opt.GRPCAddr, opt.GRPCDialOpt); err != nil {
		log.Errorf("register model handler error: %s", err.Error())
	}
}

func (mh *ModelHandler) getAccessibleProjectIDs(ctx context.Context) ([]int, error) {
	userID := user.GetCurrentUserId(ctx)
	if userID == 0 {
		return mh.authzService.GetUserAccessibleProjectIDs(ctx, 0)
	}

	allowed, err := mh.authzService.VerifyPlatformPermission(ctx, "*.*")
	if err != nil {
		return nil, err
	}
	if allowed {
		return nil, nil
	}

	return mh.authzService.GetUserAccessibleProjectIDs(ctx, userID)
}

func (mh *ModelHandler) ListModelTaskLabels(ctx context.Context, request *modelv1alpha1.ListModelTaskLabelsRequest) (*modelv1alpha1.ListModelTaskLabelsResponse, error) {
	labels, err := mh.ms.ListModelTaskLabels(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list task labels")
	}

	items := make([]*modelv1alpha1.Label, len(labels))
	for i, l := range labels {
		items[i] = &modelv1alpha1.Label{
			Id:        int32(l.ID),
			Name:      l.Name,
			Category:  modelv1alpha1.Category_TASK,
			CreatedAt: l.CreatedAt.Format(time.RFC3339),
			UpdatedAt: l.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &modelv1alpha1.ListModelTaskLabelsResponse{
		Items: items,
	}, nil
}

func (mh *ModelHandler) ListModelFrameLabels(ctx context.Context, request *modelv1alpha1.ListModelFrameLabelsRequest) (*modelv1alpha1.ListModelFrameLabelsResponse, error) {
	labels, err := mh.ms.ListModelFrameLabels(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list library labels")
	}

	items := make([]*modelv1alpha1.Label, len(labels))
	for i, l := range labels {
		items[i] = &modelv1alpha1.Label{
			Id:        int32(l.ID),
			Name:      l.Name,
			Category:  modelv1alpha1.Category_LIBRARY,
			CreatedAt: l.CreatedAt.Format(time.RFC3339),
			UpdatedAt: l.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &modelv1alpha1.ListModelFrameLabelsResponse{
		Items: items,
	}, nil
}

func (mh *ModelHandler) ListModels(ctx context.Context, request *modelv1alpha1.ListModelsRequest) (*modelv1alpha1.ListModelsResponse, error) {
	// Validate request
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	page := utils.NewPage(request.Page, request.PageSize)

	// Build filter
	filter := &model.Filter{
		Label:    request.Labels,
		Search:   request.Search,
		Sort:     request.Sort,
		Project:  request.Project,
		Page:     page.Page,
		PageSize: page.PageSize,
		Popular:  &request.Popular,
	}

	// Inject accessible project IDs
	projectIDs, err := mh.getAccessibleProjectIDs(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get accessible projects")
	}
	if projectIDs != nil {
