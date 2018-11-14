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

package matrixhub

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	v1alpha1 "github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registry"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registrydiscovery"
)

const defaultPageSize = 100

// Discovery implements registrydiscovery.Discovery for a remote MatrixHub instance.
type Discovery struct {
	client *http.Client
}

// New creates a new MatrixHub Discovery.
func New() *Discovery {
	return &Discovery{
		client: &http.Client{},
	}
}

// ListRepositories discovers remote repositories on a MatrixHub registry.
func (d *Discovery) ListRepositories(ctx context.Context, reg *registry.Registry, filter registrydiscovery.Filter) ([]registrydiscovery.RemoteRepository, error) {
