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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	registryv1alpha1 "github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registry"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
	pageutils "github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

type RegistryHandler struct {
	registryRepo registry.IRegistryRepo
}

func NewRegistryHandler(repo registry.IRegistryRepo) IHandler {
	return &RegistryHandler{
		registryRepo: repo,
	}
}

func (rh *RegistryHandler) RegisterToServer(options *ServerOptions) {
	// Register GRPC Handler
	registryv1alpha1.RegisterRegistriesServer(options.GRPCServer, rh)
	if err := registryv1alpha1.RegisterRegistriesHandlerFromEndpoint(context.Background(), options.GatewayMux, options.GRPCAddr, options.GRPCDialOpt); err != nil {
		log.Errorf("register handler error: %s", err.Error())
	}
}

func (rh *RegistryHandler) ListRegistries(ctx context.Context, request *registryv1alpha1.ListRegistriesRequest) (*registryv1alpha1.ListRegistriesResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	page := pageutils.NewPage(request.Page, request.PageSize)
	domainRegistries, total, err := rh.registryRepo.ListRegistries(ctx, int(page.Page), int(page.PageSize), request.Search)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pageutils.SetPageTotal(page, int32(total))

	var list []*registryv1alpha1.Registry
	for _, registry := range domainRegistries {
		list = append(list, convertDomainRegistryToAPIRegistry(registry))
	}

	return &registryv1alpha1.ListRegistriesResponse{
		Registries: list,
		Pagination: page,
	}, nil
}

func (rh *RegistryHandler) GetRegistry(ctx context.Context, request *registryv1alpha1.GetRegistryRequest) (*registryv1alpha1.GetRegistryResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	domainRegistry, err := rh.registryRepo.GetRegistry(ctx, int(request.Id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &registryv1alpha1.GetRegistryResponse{
		Registry: convertDomainRegistryToAPIRegistry(domainRegistry),
	}, nil
}

func (rh *RegistryHandler) CreateRegistry(ctx context.Context, request *registryv1alpha1.CreateRegistryRequest) (*registryv1alpha1.CreateRegistryResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	domainRegistry := registry.Registry{
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type.String(),
		URL:         request.Url,
		Insecure:    request.Insecure,
	}
	if b := request.GetBasic(); b != nil {
		domainRegistry.SetCredential(registry.NewBasicCredential(b.Username, b.Password))
	}

	created, err := rh.registryRepo.CreateRegistry(ctx, domainRegistry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
