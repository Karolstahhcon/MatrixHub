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

package current_user_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	v1alpha1current_user "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/current_user"
	"github.com/matrixhub-ai/matrixhub/test/tools"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// generateSSHPublicKey generates a fresh ed25519 key-pair on every call and
// returns the public key in OpenSSH authorized_keys format.  Each call
// produces a unique fingerprint, so tests never collide on the global
// uix_ssh_keys_fingerprint unique index.
func generateSSHPublicKey() string {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	pub, err := ssh.NewPublicKey(priv.Public())
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pub)))
}

var _ = Describe("CurrentUser", Label("current-user"), func() {
	var (
		ctx            context.Context
		currentUserApi *v1alpha1current_user.CurrentUserApiService
		username       string
		password       string
		userID         int32
	)

	BeforeEach(func() {
		ctx = context.Background()
		password = "Test@123456"
		username = tools.GenerateTestUsername("cu")
		var cookie string
		var err error
		userID, cookie, err = tools.CreateUserAndLoginWithID(username, password, false)
		Expect(err).NotTo(HaveOccurred())
		currentUserApi = tools.CreateCurrentUserClientWithCookie(cookie)
	})

	AfterEach(func() {
		_ = tools.DeleteUser(int64(userID))
	})

	// ═══════════════════════════════════════════════════════════
	// 1. GetCurrentUser API
	// ═══════════════════════════════════════════════════════════
	Context("GetCurrentUser API", func() {
		It("should return the correct identity for the logged-in user", Label("CU00001", "smoke"), func() {
			resp, _, err := currentUserApi.CurrentUserGetCurrentUser(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Username).To(Equal(username))
			Expect(resp.Id).NotTo(BeZero())
			Expect(resp.IsAdmin).To(BeFalse())

			GinkgoWriter.Printf("CurrentUser: id=%d, username=%s, isAdmin=%v\n", resp.Id, resp.Username, resp.IsAdmin)
		})

		It("should reflect isAdmin=true for an admin user", Label("CU00002"), func() {
			adminCookie, err := tools.LoginUser(tools.DefaultAdminUsername, tools.DefaultAdminPassword)
			Expect(err).NotTo(HaveOccurred())

			adminApi := tools.CreateCurrentUserClientWithCookie(adminCookie)
			resp, _, err := adminApi.CurrentUserGetCurrentUser(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Username).To(Equal(tools.DefaultAdminUsername))
			Expect(resp.IsAdmin).To(BeTrue())
		})
	})

	// ═══════════════════════════════════════════════════════════
	// 2. Access Token API
	// ═══════════════════════════════════════════════════════════
	Context("Access Token API", func() {
		It("should return an empty list when no tokens exist", Label("CU00003"), func() {
			resp, _, err := currentUserApi.CurrentUserListAccessTokens(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Items).To(BeEmpty())
		})

		It("should create an access token and return a non-empty token string", Label("CU00004"), func() {
			resp, _, err := currentUserApi.CurrentUserCreateAccessToken(ctx, v1alpha1current_user.V1alpha1CreateAccessTokenRequest{
				Name: "e2e-token",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Token).NotTo(BeEmpty())

			GinkgoWriter.Printf("Created access token prefix: %s...\n", resp.Token[:12])
		})

		It("should fail to create an access token with an empty name", Label("CU00005"), func() {
			_, _, err := currentUserApi.CurrentUserCreateAccessToken(ctx, v1alpha1current_user.V1alpha1CreateAccessTokenRequest{
				Name: "",
			})
			Expect(err).To(HaveOccurred(), "empty name must fail validation")
		})

		It("should list access tokens including the newly created one", Label("CU00006"), func() {
			_, _, err := currentUserApi.CurrentUserCreateAccessToken(ctx, v1alpha1current_user.V1alpha1CreateAccessTokenRequest{
				Name: "e2e-list-token",
			})
			Expect(err).NotTo(HaveOccurred())

			listResp, _, err := currentUserApi.CurrentUserListAccessTokens(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(listResp.Items).To(HaveLen(1))
			Expect(listResp.Items[0].Name).To(Equal("e2e-list-token"))

