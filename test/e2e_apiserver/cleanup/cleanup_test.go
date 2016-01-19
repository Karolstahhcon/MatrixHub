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

package cleanup_test

import (
	"context"
	"crypto/sha256"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	v1alpha1cleanup "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/cleanup"
	v1alpha1project "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/project"
	"github.com/matrixhub-ai/matrixhub/test/tools"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const hfCommandTimeout = 2 * time.Minute

var _ = Describe("Cleanup", Label("cleanup"), func() {
	var (
		ctx        context.Context
		cleanupApi *v1alpha1cleanup.CleanupApiService
	)

	BeforeEach(func() {
		ctx = context.Background()
		cleanupApi = tools.GetV1alpha1CleanupApi()
	})

	It("should report storage stats", Label("CL00001", "smoke"), func() {
		stats, resp, err := cleanupApi.CleanupGetStorageStats(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		for _, size := range []string{stats.TotalSizeBytes, stats.RepositoriesSizeBytes, stats.LfsSizeBytes, stats.OrphanedSizeBytes} {
			n, err := strconv.ParseInt(size, 10, 64)
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(BeNumerically(">=", 0))
		}
	})

	It("should preview and dry-run cleanup", Label("CL00002", "smoke"), func() {
		_, resp, err := cleanupApi.CleanupPreviewCleanup(ctx, v1alpha1cleanup.V1alpha1PreviewCleanupRequest{
			IncludeOrphanedRepos: true,
			IncludeOrphanedLfs:   true,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		res, resp, err := cleanupApi.CleanupExecuteCleanup(ctx, v1alpha1cleanup.V1alpha1ExecuteCleanupRequest{
			CleanOrphanedRepos: true,
			CleanOrphanedLfs:   true,
			DryRun:             true,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(res.Errors).To(BeEmpty())
	})
