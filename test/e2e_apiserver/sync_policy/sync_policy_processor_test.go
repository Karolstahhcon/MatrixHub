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

	"github.com/antihax/optional"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1alpha1 "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/sync_policy"
	"github.com/matrixhub-ai/matrixhub/test/tools"
)

var _ = Describe("SyncPolicy Processor", func() {
	It("creates scheduled policy and observes processor-created sync task", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
		defer cancel()

		api := tools.GetV1alpha1SyncPolicyApi()

		name := fmt.Sprintf("e2e-syncpolicy-processor-%d", time.Now().UnixNano())
		policyType := v1alpha1.V1alpha1SyncPolicyType("SYNC_POLICY_TYPE_PULL_BASE")
		triggerType := v1alpha1.V1alpha1TriggerType("TRIGGER_TYPE_SCHEDULED")

		project, err := tools.CreateProjectFixture(ctx, "e2e-syncpolicy-processor-project")
		Expect(err).NotTo(HaveOccurred())
		defer project.Cleanup(ctx)

		registry, err := tools.CreateHuggingFaceRegistryFixture(ctx, "e2e-syncpolicy-processor-registry", "https://hf-mirror.com")
		Expect(err).NotTo(HaveOccurred())
		defer registry.Cleanup(ctx)
