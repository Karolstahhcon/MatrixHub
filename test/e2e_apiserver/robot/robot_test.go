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

package robot_test

import (
	"context"
	"strings"

	"github.com/antihax/optional"
	v1alpha1robot "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/robot"
	"github.com/matrixhub-ai/matrixhub/test/tools"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// robotPrefix is prepended to every robot name by the server.
const robotPrefix = "robot$"

// generateRobotName returns a unique robot name (without prefix) for each test.
func generateRobotName(prefix string) string {
	return tools.GenerateTestUsername(prefix)
}

var _ = Describe("Robot", Label("robot"), func() {
	var (
		ctx       context.Context
		robotsApi *v1alpha1robot.RobotsApiService
	)

	BeforeEach(func() {
		ctx = context.Background()
		robotsApi = tools.GetV1alpha1RobotsApi()
	})

	// ═══════════════════════════════════════════════════════════
	// 1. CreateRobotAccount API
	// ═══════════════════════════════════════════════════════════
	Context("CreateRobotAccount API", func() {
		It("should create a robot account and return a non-empty token", Label("R00001", "smoke"), func() {
			name := generateRobotName("rb-create")
			resp, _, err := robotsApi.RobotsCreateRobotAccount(ctx, v1alpha1robot.V1alpha1CreateRobotAccountRequest{
				Name:        name,
				Description: "e2e test robot",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Token).NotTo(BeEmpty())

			// Cleanup: find the robot and delete it.
			listResp, _, err := robotsApi.RobotsListRobotAccounts(ctx, &v1alpha1robot.RobotsApiRobotsListRobotAccountsOpts{
				Search: optional.NewString(robotPrefix + name),
			})
			Expect(err).NotTo(HaveOccurred())
			for _, r := range listResp.Items {
				if r.Name == robotPrefix+name {
					_, _, _ = robotsApi.RobotsDeleteRobotAccount(ctx, r.Id)
				}
			}

			GinkgoWriter.Printf("Created robot token prefix: %s...\n", resp.Token[:12])
		})

		It("should fail to create a robot account with an empty name", Label("R00002"), func() {
			_, _, err := robotsApi.RobotsCreateRobotAccount(ctx, v1alpha1robot.V1alpha1CreateRobotAccountRequest{
				Name: "",
			})
			Expect(err).To(HaveOccurred(), "empty name must fail validation")
