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

	v1alpha1 "github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/cleanup"
	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
)

// CleanupHandler handles cleanup API requests.
type CleanupHandler struct {
	cleanupService cleanup.ICleanupService
}

// NewCleanupHandler creates a new CleanupHandler instance.
func NewCleanupHandler(cleanupService cleanup.ICleanupService) IHandler {
	return &CleanupHandler{
		cleanupService: cleanupService,
	}
}

// RegisterToServer registers the handler to gRPC and HTTP gateway.
func (h *CleanupHandler) RegisterToServer(options *ServerOptions) {
	v1alpha1.RegisterCleanupServer(options.GRPCServer, h)
	if err := v1alpha1.RegisterCleanupHandlerFromEndpoint(context.Background(), options.GatewayMux, options.GRPCAddr, options.GRPCDialOpt); err != nil {
		log.Errorf("register cleanup handler error: %s", err.Error())
	}
