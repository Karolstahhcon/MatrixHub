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

package tools

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Environment variable names
const (
	EnvMatrixHubBaseURL = "MATRIXHUB_BASE_URL" // HTTP base URL (e.g., http://localhost:3001)
	// EnvMatrixHubSelfURL is the URL by which MatrixHub reaches itself, for
	// tests that make the deployment act as its own remote registry. It is not
	// the same as EnvMatrixHubBaseURL: that one is resolved by the test process,
	// whereas this one is resolved by the server components — under KIND the
	// test talks to a NodePort on the host while the job server has to go
	// through the in-cluster Service.
	EnvMatrixHubSelfURL = "MATRIXHUB_SELF_URL"
)

// Default values
const (
	DefaultBaseURL = "http://localhost:3001"
)

var (
	baseURL   string
	closeOnce sync.Once
)

// Close cleans up resources
func Close() {
	closeOnce.Do(func() {
