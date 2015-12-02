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

package hf_cli_test

import (
	"context"
	"crypto/sha256"
	"os"
	"path/filepath"
	"time"

	v1alpha1project "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/project"
	"github.com/matrixhub-ai/matrixhub/test/tools"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const hfCommandTimeout = 2 * time.Minute

var _ = Describe("HF CLI", Label("hf-cli"), func() {
	var (
		ctx     context.Context
		fixture *tools.ProjectUserFixture
		root    string
		token   string
		env     map[string]string
	)

	BeforeEach(func() {
		ctx = context.Background()
		root = GinkgoT().TempDir()

		var err error
		fixture, err = tools.CreateProjectUserFixture(ctx, "hf-cli", v1alpha1project.EDITOR_V1alpha1ProjectRoleType)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			if err := fixture.Cleanup(context.Background()); err != nil {
				GinkgoWriter.Printf("protocol fixture cleanup failed: %v\n", err)
			}
		})

		token, err = fixture.CreateAccessToken(ctx)
		Expect(err).NotTo(HaveOccurred())
		env = tools.HFCLIEnvironment(root, token)
	})

	It("should authenticate, upload, and download a model with the real hf CLI", Label("HF00001", "smoke"), func() {
		model := tools.GenerateTestModelName("hf-model")
		repoID := fixture.Project.Name + "/" + model

		whoami := runHF(root, env, "auth", "whoami")
		Expect(whoami.Err).NotTo(HaveOccurred(), whoami.FailureMessage())
		Expect(whoami.Stdout).To(ContainSubstring(fixture.Username))

		create := runHF(root, env, "repos", "create", repoID, "--repo-type", "model")
		Expect(create.Err).NotTo(HaveOccurred(), create.FailureMessage())
