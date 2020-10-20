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
			Expect(resp.SyncPolicy.IsOverwrite).To(BeFalse())
			Expect(resp.SyncPolicy.IsDisabled).To(BeFalse())

			// Get and verify
			getResp, _, err := api.SyncPolicyGetSyncPolicy(ctx, pid)
			Expect(err).NotTo(HaveOccurred())
			Expect(getResp.SyncPolicy).NotTo(BeNil())
			Expect(getResp.SyncPolicy.Name).To(Equal(name))
			Expect(getResp.SyncPolicy.PolicyType).NotTo(BeNil())
			Expect(*getResp.SyncPolicy.PolicyType).To(Equal(v1alpha1.PULL_BASE_V1alpha1SyncPolicyType))
		})
	})

	Context("Update", func() {
		It("should update policy name and trigger type", Label("SP0002"), func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			project, err := tools.CreateProjectFixture(ctx, "e2e-update-project")
			Expect(err).NotTo(HaveOccurred())
			defer project.Cleanup(ctx)

			registry, err := tools.CreateHuggingFaceRegistryFixture(ctx, "e2e-update-registry", "https://hf-mirror.com")
			Expect(err).NotTo(HaveOccurred())
			defer registry.Cleanup(ctx)
			Expect(registry.ID).To(BeNumerically(">", 0))

			name := fmt.Sprintf("e2e-update-%d", time.Now().UnixNano())
			req := v1alpha1.V1alpha1CreateSyncPolicyRequest{
				Name:        name,
				Description: "before update",
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
			pid := resp.SyncPolicy.Id
			defer func() {
				_, _, _ = api.SyncPolicyDeleteSyncPolicy(ctx, pid)
			}()

			// Update policy
			newName := name + "-updated"
			updateReq := v1alpha1.SyncPolicyUpdateSyncPolicyBody{
				Name:        newName,
				Description: "after update",
				TriggerType: ptrTriggerType(v1alpha1.SCHEDULED_V1alpha1TriggerType),
				PullBasePolicy: &v1alpha1.V1alpha1PullBasePolicy{
					SourceRegistryId:  registry.ID,
					ResourceName:      "demo-org/demo-model",
					ResourceTypes:     []v1alpha1.V1alpha1ResourceType{v1alpha1.MODEL_V1alpha1ResourceType},
					TargetProjectName: project.Name,
				},
				TriggerTypeSchedule: &v1alpha1.V1alpha1TriggerTypeSchedule{
					Cron: "*/5 * * * *",
				},
				IsOverwrite: true,
			}

			_, _, err = api.SyncPolicyUpdateSyncPolicy(ctx, pid, updateReq)
			Expect(err).NotTo(HaveOccurred())

			// Verify update
			getResp, _, err := api.SyncPolicyGetSyncPolicy(ctx, pid)
			Expect(err).NotTo(HaveOccurred())
			Expect(getResp.SyncPolicy.Name).To(Equal(newName))
			Expect(getResp.SyncPolicy.Description).To(Equal("after update"))
			Expect(getResp.SyncPolicy.TriggerType).NotTo(BeNil())
			Expect(*getResp.SyncPolicy.TriggerType).To(Equal(v1alpha1.SCHEDULED_V1alpha1TriggerType))
			Expect(getResp.SyncPolicy.IsOverwrite).To(BeTrue())
			Expect(getResp.SyncPolicy.TriggerTypeSchedule).NotTo(BeNil())
			Expect(getResp.SyncPolicy.TriggerTypeSchedule.Cron).To(Equal("*/5 * * * *"))
		})
	})

	Context("List", func() {
		It("should list policies with pagination and search", Label("SP0003"), func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			project, err := tools.CreateProjectFixture(ctx, "e2e-list-project")
			Expect(err).NotTo(HaveOccurred())
			defer project.Cleanup(ctx)

			registry, err := tools.CreateHuggingFaceRegistryFixture(ctx, "e2e-list-registry", "https://hf-mirror.com")
			Expect(err).NotTo(HaveOccurred())
			defer registry.Cleanup(ctx)
			Expect(registry.ID).To(BeNumerically(">", 0))

			// Create two policies
			prefix := fmt.Sprintf("e2e-list-%d", time.Now().UnixNano())
			name1 := prefix + "-policy-1"
			name2 := prefix + "-policy-2"

			req1 := v1alpha1.V1alpha1CreateSyncPolicyRequest{
				Name:        name1,
				Description: "list test 1",
				PolicyType:  ptrPolicyType(v1alpha1.PULL_BASE_V1alpha1SyncPolicyType),
				TriggerType: ptrTriggerType(v1alpha1.MANUAL_V1alpha1TriggerType),
				PullBasePolicy: &v1alpha1.V1alpha1PullBasePolicy{
					SourceRegistryId:  registry.ID,
					ResourceName:      "demo-org/demo-model",
					ResourceTypes:     []v1alpha1.V1alpha1ResourceType{v1alpha1.MODEL_V1alpha1ResourceType},
					TargetProjectName: project.Name,
				},
			}
			resp1, _, err := api.SyncPolicyCreateSyncPolicy(ctx, req1)
			Expect(err).NotTo(HaveOccurred())
			pid1 := resp1.SyncPolicy.Id
			defer func() {
				_, _, _ = api.SyncPolicyDeleteSyncPolicy(ctx, pid1)
			}()

			req2 := v1alpha1.V1alpha1CreateSyncPolicyRequest{
				Name:        name2,
				Description: "list test 2",
				PolicyType:  ptrPolicyType(v1alpha1.PULL_BASE_V1alpha1SyncPolicyType),
				TriggerType: ptrTriggerType(v1alpha1.MANUAL_V1alpha1TriggerType),
				PullBasePolicy: &v1alpha1.V1alpha1PullBasePolicy{
					SourceRegistryId:  registry.ID,
					ResourceName:      "demo-org/demo-model-2",
					ResourceTypes:     []v1alpha1.V1alpha1ResourceType{v1alpha1.MODEL_V1alpha1ResourceType},
					TargetProjectName: project.Name,
				},
			}
			resp2, _, err := api.SyncPolicyCreateSyncPolicy(ctx, req2)
			Expect(err).NotTo(HaveOccurred())
			pid2 := resp2.SyncPolicy.Id
			defer func() {
				_, _, _ = api.SyncPolicyDeleteSyncPolicy(ctx, pid2)
			}()

			// List with search
			listResp, _, err := api.SyncPolicyListSyncPolicies(ctx, &v1alpha1.SyncPolicyApiSyncPolicyListSyncPoliciesOpts{
				Page:     optional.NewInt32(1),
				PageSize: optional.NewInt32(20),
				Search:   optional.NewString(prefix),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(listResp.Pagination).NotTo(BeNil())
			Expect(listResp.Pagination.Page).To(Equal(int32(1)))
			Expect(listResp.Pagination.PageSize).To(Equal(int32(20)))
			Expect(listResp.Pagination.Total).To(BeNumerically(">=", 2))

			// Verify both policies are in the list
			var found1, found2 bool
			for _, p := range listResp.SyncPolicies {
				if p.Name == name1 {
					found1 = true
				}
				if p.Name == name2 {
					found2 = true
				}
			}
