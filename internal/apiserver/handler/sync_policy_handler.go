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
	"io"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1alpha1 "github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registry"
	"github.com/matrixhub-ai/matrixhub/internal/domain/syncjob"
	"github.com/matrixhub-ai/matrixhub/internal/domain/syncpolicy"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
	"github.com/matrixhub-ai/matrixhub/internal/jobserver/canceller"
	"github.com/matrixhub-ai/matrixhub/internal/jobserver/logstore"
)

type SyncPolicyHandler struct {
	syncPolicyService syncpolicy.ISyncPolicyService
	syncJobService    syncjob.ISyncJobService
	registryRepo      registry.IRegistryRepo
	logStore          logstore.LogStore
	canceller         canceller.Canceller
}

func NewSyncPolicyHandler(syncPolicyService syncpolicy.ISyncPolicyService, syncJobService syncjob.ISyncJobService, registryRepo registry.IRegistryRepo, logStore logstore.LogStore, canc canceller.Canceller) IHandler {
	return &SyncPolicyHandler{
		syncPolicyService: syncPolicyService,
		syncJobService:    syncJobService,
		registryRepo:      registryRepo,
		logStore:          logStore,
		canceller:         canc,
	}
}

func (h *SyncPolicyHandler) RegisterToServer(options *ServerOptions) {
	// Register GRPC Handler
	v1alpha1.RegisterSyncPolicyServer(options.GRPCServer, h)
	if err := v1alpha1.RegisterSyncPolicyHandlerFromEndpoint(context.Background(), options.GatewayMux, options.GRPCAddr, options.GRPCDialOpt); err != nil {
		log.Errorf("register sync policy handler error: %s", err.Error())
	}
}

// ListSyncPolicies lists all sync policies with pagination and search
func (h *SyncPolicyHandler) ListSyncPolicies(ctx context.Context, request *v1alpha1.ListSyncPoliciesRequest) (*v1alpha1.ListSyncPoliciesResponse, error) {
	// Validate request
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	page := int(request.Page)
	pageSize := int(request.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// Call service
	policies, total, err := h.syncPolicyService.ListSyncPolicies(ctx, page, pageSize, request.Search)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list sync policies")
	}

	// Convert to proto
	items := make([]*v1alpha1.SyncPolicyItem, len(policies))
	for i, p := range policies {
		items[i] = h.syncPolicyToProto(ctx, p)
	}

	return &v1alpha1.ListSyncPoliciesResponse{
		SyncPolicies: items,
		Pagination: &v1alpha1.Pagination{
			Total:    utils.ClampInt32(total),
			Page:     int32(page),
			PageSize: int32(pageSize),
			Pages:    utils.CalculatePages(total, int32(pageSize)),
		},
	}, nil
}

// GetSyncPolicy gets a sync policy by ID
func (h *SyncPolicyHandler) GetSyncPolicy(ctx context.Context, request *v1alpha1.GetSyncPolicyRequest) (*v1alpha1.GetSyncPolicyResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	policy, err := h.syncPolicyService.GetSyncPolicy(ctx, int(request.SyncPolicyId))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, "sync policy not found")
		}
		return nil, status.Error(codes.Internal, "failed to get sync policy")
	}

	return &v1alpha1.GetSyncPolicyResponse{
		SyncPolicy: h.syncPolicyToProto(ctx, policy),
	}, nil
}

// CreateSyncPolicy creates a new sync policy
func (h *SyncPolicyHandler) CreateSyncPolicy(ctx context.Context, request *v1alpha1.CreateSyncPolicyRequest) (*v1alpha1.CreateSyncPolicyResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	policy := &syncpolicy.SyncPolicy{}
	h.applyCommonFields(policy, request.Name, request.Description, request.Bandwidth,
		request.TriggerTypeSchedule, syncpolicy.TriggerType(request.TriggerType),
		request.IsOverwrite, false)

	switch request.PolicyType {
	case v1alpha1.SyncPolicyType_SYNC_POLICY_TYPE_PULL_BASE:
		pullPolicy := request.GetPullBasePolicy()
		if pullPolicy == nil {
			return nil, status.Error(codes.InvalidArgument, "pull_base_policy is required")
		}
		h.applyPullBasePolicy(policy, pullPolicy)
	case v1alpha1.SyncPolicyType_SYNC_POLICY_TYPE_PUSH_BASE:
		pushPolicy := request.GetPushBasePolicy()
		if pushPolicy == nil {
			return nil, status.Error(codes.InvalidArgument, "push_base_policy is required")
		}
		h.applyPushBasePolicy(policy, pushPolicy)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid policy_type")
	}

	if err := h.syncPolicyService.CreateSyncPolicy(ctx, policy); err != nil {
		return nil, status.Error(codes.Internal, "failed to create sync policy")
	}

	return &v1alpha1.CreateSyncPolicyResponse{
		SyncPolicy: h.syncPolicyToProto(ctx, policy),
	}, nil
}

// UpdateSyncPolicy updates a sync policy
func (h *SyncPolicyHandler) UpdateSyncPolicy(ctx context.Context, request *v1alpha1.UpdateSyncPolicyRequest) (*v1alpha1.UpdateSyncPolicyResponse, error) {
	if err := request.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	existingPolicy, err := h.syncPolicyService.GetSyncPolicy(ctx, int(request.SyncPolicyId))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, "sync policy not found")
		}
		return nil, status.Error(codes.Internal, "failed to get sync policy")
	}

	h.applyCommonFields(existingPolicy, request.Name, request.Description, request.Bandwidth,
		request.TriggerTypeSchedule, syncpolicy.TriggerType(request.TriggerType),
		request.IsOverwrite, request.IsDisabled)
