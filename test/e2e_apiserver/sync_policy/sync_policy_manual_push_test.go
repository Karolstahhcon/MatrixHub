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

var _ = Describe("SyncPolicy Manual Push", Label("sync-policy", "git"), func() {
	It("should execute manual push sync end-to-end", Label("SP0020"), func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		api := tools.GetV1alpha1SyncPolicyApi()

		// 1. Get preset local git model
		gitProject := tools.GetGitModelProject()
		gitModel := tools.GetGitModelName()

		// Verify local model exists, otherwise skip
		modelsApi := tools.GetV1alpha1ModelsApi()
		_, _, err := modelsApi.ModelsGetModel(ctx, gitProject, gitModel)
		if err != nil {
			Skip("GIT_MODEL not available: " + gitProject + "/" + gitModel)
		}

		// 2. Create registry (target remote)
		registry, err := tools.CreateHuggingFaceRegistryFixture(ctx, "e2e-push", "https://hf-mirror.com")
		Expect(err).NotTo(HaveOccurred())
		defer registry.Cleanup(ctx)
		Expect(registry.ID).To(BeNumerically(">", 0))

		// 3. Create Manual + Push Base Policy
		name := fmt.Sprintf("e2e-manual-push-%d", time.Now().UnixNano())
		req := v1alpha1.V1alpha1CreateSyncPolicyRequest{
			Name:        name,
			Description: "e2e manual push test",
			PolicyType:  ptrPolicyType(v1alpha1.PUSH_BASE_V1alpha1SyncPolicyType),
			TriggerType: ptrTriggerType(v1alpha1.MANUAL_V1alpha1TriggerType),
			PushBasePolicy: &v1alpha1.V1alpha1PushBasePolicy{
				ResourceName:      fmt.Sprintf("%s/%s", gitProject, gitModel),
				ResourceTypes:     []v1alpha1.V1alpha1ResourceType{v1alpha1.MODEL_V1alpha1ResourceType},
				TargetRegistryId:  registry.ID,
				TargetProjectName: fmt.Sprintf("e2e-push-target-%d", time.Now().UnixNano()),
			},
			IsOverwrite: false,
		}

		resp, _, err := api.SyncPolicyCreateSyncPolicy(ctx, req)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.SyncPolicy).NotTo(BeNil())
		pid := resp.SyncPolicy.Id
		Expect(pid).To(BeNumerically(">", 0))
		defer func() {
