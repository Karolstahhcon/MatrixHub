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

package preset_test

import (
	"context"
	"strconv"
	"strings"

	"github.com/antihax/optional"

	"github.com/matrixhub-ai/matrixhub/internal/domain/preset"
	v1alpha1model "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/model"
	v1alpha1project "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/project"
	v1alpha1registry "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/registry"
	"github.com/matrixhub-ai/matrixhub/test/tools"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Preset initialization", Label("preset"), func() {
	It("should expose preset registries, proxy projects, admin memberships and recommended models", Label("PR00001", "smoke"), func() {
		ctx := context.Background()
