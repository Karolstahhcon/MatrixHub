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
	"strings"
	"time"

	"github.com/antihax/optional"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v1alpha1 "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/sync_policy"
	"github.com/matrixhub-ai/matrixhub/test/tools"
)

var _ = Describe("SyncPolicy CRUD APIs", Label("sync-policy"), func() {
	var api *v1alpha1.SyncPolicyApiService

	BeforeEach(func() {
		api = tools.GetV1alpha1SyncPolicyApi()
	})

	Context("Create and Get", func() {
		It("should create a pull-base policy and retrieve it", Label("SP0001", "smoke"), func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			project, err := tools.CreateProjectFixture(ctx, "e2e-crud-project")
			Expect(err).NotTo(HaveOccurred())
			defer project.Cleanup(ctx)

			registry, err := tools.CreateHuggingFaceRegistryFixture(ctx, "e2e-crud-registry", "https://hf-mirror.com")
			Expect(err).NotTo(HaveOccurred())
			defer registry.Cleanup(ctx)
			Expect(registry.ID).To(BeNumerically(">", 0))

			name := fmt.Sprintf("e2e-crud-pull-%d", time.Now().UnixNano())
			req := v1alpha1.V1alpha1CreateSyncPolicyRequest{
				Name:        name,
				Description: "e2e crud pull test",
				PolicyType:  ptrPolicyType(v1alpha1.PULL_BASE_V1alpha1SyncPolicyType),
				TriggerType: ptrTriggerType(v1alpha1.MANUAL_V1alpha1TriggerType),
				PullBasePolicy: &v1alpha1.V1alpha1PullBasePolicy{
					SourceRegistryId:  registry.ID,
					ResourceName:      "demo-org/demo-model",
					ResourceTypes:     []v1alpha1.V1alpha1ResourceType{v1alpha1.MODEL_V1alpha1ResourceType},
					TargetProjectName: project.Name,
				},
				IsOverwrite: false,
			}

			resp, _, err := api.SyncPolicyCreateSyncPolicy(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.SyncPolicy).NotTo(BeNil())
			pid := resp.SyncPolicy.Id
			Expect(pid).To(BeNumerically(">", 0))
			defer func() {
				_, _, _ = api.SyncPolicyDeleteSyncPolicy(ctx, pid)
			}()

			// Verify returned fields
			Expect(resp.SyncPolicy.Name).To(Equal(name))
			Expect(resp.SyncPolicy.Description).To(Equal("e2e crud pull test"))
			Expect(resp.SyncPolicy.PolicyType).NotTo(BeNil())
			Expect(*resp.SyncPolicy.PolicyType).To(Equal(v1alpha1.PULL_BASE_V1alpha1SyncPolicyType))
