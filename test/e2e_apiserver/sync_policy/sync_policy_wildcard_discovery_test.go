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

package sync_policy_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1alpha1dataset "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/dataset"
	v1alpha1model "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/model"
	v1alpha1 "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/sync_policy"
	"github.com/matrixhub-ai/matrixhub/test/tools"
)

// A wildcard pull policy resolves the concrete resource list through registry
// discovery. Discovery must query the registry's configured URL — the same host
// the resulting sync jobs pull from. When it queries a hardcoded upstream
// instead, an operator who points a registry at a mirror or an air-gapped proxy
// gets an empty or wrong match set while the pull half of the pipeline still
// targets the mirror.
//
// MatrixHub serves a HuggingFace-compatible /api/models, so this suite points a
// HUGGINGFACE-typed registry back at the instance under test and uses it as the
// mirror. That keeps the assertion hermetic: the repositories that must be
