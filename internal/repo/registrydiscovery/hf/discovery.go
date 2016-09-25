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

package hf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/matrixhub-ai/matrixhub/internal/domain/registry"
	"github.com/matrixhub-ai/matrixhub/internal/domain/registrydiscovery"
)

// Discovery implements registrydiscovery.Discovery for HuggingFace-compatible
// registries. It has no built-in host: the API base always comes from the
// registry being synced, so a registry pointed at a mirror or an air-gapped
// proxy is discovered against that host and never against a default one.
type Discovery struct {
	client *http.Client
}

// Option configures a Discovery.
type Option func(*Discovery)

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(disc *Discovery) {
		disc.client = &http.Client{Timeout: d}
	}
}

// New creates a new HF Discovery with sensible defaults.
func New(opts ...Option) *Discovery {
	d := &Discovery{
		client: &http.Client{Timeout: 60 * time.Second},
	}
	for _, o := range opts {
		o(d)
	}
	return d
}

// endpointBase returns the API base for a registry. A registry without a URL is
// an error rather than a fallback to some default host: the sync path fails on
// such a registry too, and silently substituting another host would send the
// registry's credential to somewhere the operator never configured.
func endpointBase(reg *registry.Registry) (string, error) {
	if reg == nil {
