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

package protocol_interop_test

import (
	"context"
	"os"
	"path/filepath"
	"time"

	v1alpha1project "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/project"
	"github.com/matrixhub-ai/matrixhub/test/tools"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HF and Git interoperability", Label("protocol-interop"), func() {
	It("should expose HF uploads to Git and Git pushes to HF", Label("PI00001", "smoke"), func() {
		ctx := context.Background()
		root := GinkgoT().TempDir()

		fixture, err := tools.CreateProjectUserFixture(ctx, "interop", v1alpha1project.EDITOR_V1alpha1ProjectRoleType)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			if err := fixture.Cleanup(context.Background()); err != nil {
				GinkgoWriter.Printf("protocol fixture cleanup failed: %v\n", err)
			}
