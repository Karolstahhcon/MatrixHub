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

package user_test

import (
	"context"

	"github.com/antihax/optional"
	v1alpha1user "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/user"
	"github.com/matrixhub-ai/matrixhub/test/tools"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", Label("user"), func() {
	var (
		ctx      context.Context
		usersApi *v1alpha1user.UsersApiService
	)

	BeforeEach(func() {
		ctx = context.Background()
		usersApi = tools.GetV1alpha1UsersApi()
	})

	// ═══════════════════════════════════════════════════════════
	// 1. CreateUser API
	// ═══════════════════════════════════════════════════════════
	Context("CreateUser API", func() {
		It("should create a regular user successfully", Label("U00001", "smoke"), func() {
			username := tools.GenerateTestUsername("u-create")
			_, _, err := usersApi.UsersCreateUser(ctx, v1alpha1user.V1alpha1CreateUserRequest{
				Username: username,
				Password: "Test@123456",
			})
			Expect(err).NotTo(HaveOccurred())

			// Cleanup
			id, err := tools.GetUserIDByUsername(username)
			Expect(err).NotTo(HaveOccurred())
			_ = tools.DeleteUser(id)

			GinkgoWriter.Printf("Created user: %s (id=%d)\n", username, id)
		})

		It("should create an admin user successfully", Label("U00002"), func() {
			username := tools.GenerateTestUsername("u-admin")
			_, _, err := usersApi.UsersCreateUser(ctx, v1alpha1user.V1alpha1CreateUserRequest{
				Username: username,
				Password: "Test@123456",
				IsAdmin:  true,
			})
			Expect(err).NotTo(HaveOccurred())

			id, err := tools.GetUserIDByUsername(username)
			Expect(err).NotTo(HaveOccurred())

			// Verify isAdmin flag is set.
			listResp, _, err := usersApi.UsersListUsers(ctx, &v1alpha1user.UsersApiUsersListUsersOpts{
				Search: optional.NewString(username),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(listResp.Users).NotTo(BeEmpty())
			Expect(listResp.Users[0].IsAdmin).To(BeTrue())

			_ = tools.DeleteUser(id)
		})

		It("should fail to create a user with an empty username", Label("U00003"), func() {
			_, _, err := usersApi.UsersCreateUser(ctx, v1alpha1user.V1alpha1CreateUserRequest{
				Username: "",
				Password: "Test@123456",
			})
			Expect(err).To(HaveOccurred(), "empty username must fail validation")
		})

		It("should fail to create a user with an empty password", Label("U00004"), func() {
			_, _, err := usersApi.UsersCreateUser(ctx, v1alpha1user.V1alpha1CreateUserRequest{
				Username: tools.GenerateTestUsername("u-nopwd"),
				Password: "",
			})
			Expect(err).To(HaveOccurred(), "empty password must fail validation")
		})

		It("should fail to create a user with a duplicate username", Label("U00005"), func() {
			username := tools.GenerateTestUsername("u-dup")
			_, _, err := usersApi.UsersCreateUser(ctx, v1alpha1user.V1alpha1CreateUserRequest{
				Username: username,
				Password: "Test@123456",
			})
			Expect(err).NotTo(HaveOccurred())

			_, _, err = usersApi.UsersCreateUser(ctx, v1alpha1user.V1alpha1CreateUserRequest{
				Username: username,
				Password: "Test@123456",
			})
			Expect(err).To(HaveOccurred(), "duplicate username must be rejected")

			id, err := tools.GetUserIDByUsername(username)
			Expect(err).NotTo(HaveOccurred())
			_ = tools.DeleteUser(id)
		})
	})

	// ═══════════════════════════════════════════════════════════
	// 2. GetUser API
	// ═══════════════════════════════════════════════════════════
	Context("GetUser API", func() {
		var (
			username string
			userID   int64
		)

		BeforeEach(func() {
			username = tools.GenerateTestUsername("u-get")
			_, _, err := usersApi.UsersCreateUser(ctx, v1alpha1user.V1alpha1CreateUserRequest{
				Username: username,
				Password: "Test@123456",
			})
			Expect(err).NotTo(HaveOccurred())
			userID, err = tools.GetUserIDByUsername(username)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_ = tools.DeleteUser(userID)
		})

		It("should get a user by ID and return expected fields", Label("U00006"), func() {
			resp, _, err := usersApi.UsersGetUser(ctx, userID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Id).To(Equal(userID))
			Expect(resp.Username).To(Equal(username))
			Expect(resp.CreatedAt).NotTo(BeZero())

			GinkgoWriter.Printf("GetUser: id=%d, username=%s\n", resp.Id, resp.Username)
		})

		It("should fail to get a non-existent user", Label("U00007"), func() {
			_, _, err := usersApi.UsersGetUser(ctx, 999999999)
			Expect(err).To(HaveOccurred(), "non-existent user ID must return an error")
		})
	})

	// ═══════════════════════════════════════════════════════════
	// 3. ListUsers API
	// ═══════════════════════════════════════════════════════════
	Context("ListUsers API", func() {
		var (
			username string
			userID   int64
		)

		BeforeEach(func() {
			username = tools.GenerateTestUsername("u-list")
			_, _, err := usersApi.UsersCreateUser(ctx, v1alpha1user.V1alpha1CreateUserRequest{
				Username: username,
				Password: "Test@123456",
			})
			Expect(err).NotTo(HaveOccurred())
			userID, err = tools.GetUserIDByUsername(username)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_ = tools.DeleteUser(userID)
		})

		It("should return a paginated user list", Label("U00008"), func() {
			resp, _, err := usersApi.UsersListUsers(ctx, &v1alpha1user.UsersApiUsersListUsersOpts{
