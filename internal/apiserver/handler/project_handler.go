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
	"errors"
	"strings"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	projectv1alpha1 "github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/auth"
	"github.com/matrixhub-ai/matrixhub/internal/domain/authz"
	"github.com/matrixhub-ai/matrixhub/internal/domain/project"
	"github.com/matrixhub-ai/matrixhub/internal/domain/robot"
	"github.com/matrixhub-ai/matrixhub/internal/domain/role"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

type ProjectHandler struct {
	projectRepo  project.IProjectRepo
	authzService authz.IAuthzService
}

func NewProjectHandler(repo project.IProjectRepo, authzService authz.IAuthzService) IHandler {
	return &ProjectHandler{
		projectRepo:  repo,
		authzService: authzService,
	}
}

func (h *ProjectHandler) RegisterToServer(opt *ServerOptions) {
	projectv1alpha1.RegisterProjectsServer(opt.GRPCServer, h)
	if err := projectv1alpha1.RegisterProjectsHandlerFromEndpoint(context.Background(), opt.GatewayMux, opt.GRPCAddr, opt.GRPCDialOpt); err != nil {
		log.Errorf("register handler error: %s", err.Error())
	}
}

func (h *ProjectHandler) CreateProject(ctx context.Context, req *projectv1alpha1.CreateProjectRequest) (*projectv1alpha1.CreateProjectResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !hasAtLeastTwoDistinctChars(req.GetName()) {
		return nil, status.Error(codes.InvalidArgument, "project name must contain at least 2 distinct characters")
	}
	if strings.TrimSpace(req.GetOrganization()) != "" && req.GetRegistryId() == nil {
		return nil, status.Error(codes.InvalidArgument, "registry is required when organization is specified")
	}

	existingProject, err := h.projectRepo.GetProjectByName(ctx, req.GetName())
	if err == nil && existingProject != nil {
		return nil, status.Error(codes.AlreadyExists, "project with this name already exists")
	}

	p := &project.Project{
		Name:         req.GetName(),
		Type:         convertProtoTypeToDomain(req.GetType()),
		Organization: req.GetOrganization(),
	}

	if req.RegistryId != nil {
		registryID := int(req.RegistryId.GetValue())
		p.RegistryID = &registryID
	}

	if _, err := h.projectRepo.CreateProject(ctx, p); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &projectv1alpha1.CreateProjectResponse{}, nil
}

func (h *ProjectHandler) GetProject(ctx context.Context, req *projectv1alpha1.GetProjectRequest) (*projectv1alpha1.GetProjectResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	p, err := h.projectRepo.GetProjectByName(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if allowed, err := h.authzService.VerifyProjectPermission(ctx, p.ID, role.ProjectGet); err != nil || !allowed {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	resp := &projectv1alpha1.GetProjectResponse{
		Name:         p.Name,
		Type:         convertDomainTypeToProto(p.Type),
		RegistryUrl:  p.RegistryURL,
		Organization: p.Organization,
		ModelCount:   uint32(p.ModelCount),
		DatasetCount: uint32(p.DatasetCount),
		UpdatedAt:    timestamppb.New(p.UpdatedAt),
	}

	return resp, nil
