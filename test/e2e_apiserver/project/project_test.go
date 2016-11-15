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

package project_test

import (
	"context"
	"fmt"
	"strings"

	v1alpha1project "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/project"
	v1alpha1user "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/user"
	"github.com/matrixhub-ai/matrixhub/test/tools"

	"github.com/antihax/optional"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project", Label("project"), func() {
	var (
		ctx         context.Context
		projectsApi *v1alpha1project.ProjectsApiService
	)

	BeforeEach(func() {
		ctx = context.Background()
		projectsApi = tools.GetV1alpha1ProjectsApi()
	})

	Context("CreateProject API", func() {
		It("should create a project successfully", Label("L00001", "smoke"), func() {
			projectName := tools.GenerateTestProjectName("project")
			projectType := v1alpha1project.PRIVATE_V1alpha1ProjectType
			GinkgoWriter.Printf("Creating project: %v\n", projectName)

			resp, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name:  projectName,
				Type_: &projectType,
			})
			Expect(err).NotTo(HaveOccurred())

			GinkgoWriter.Printf("Create response: project=%v\n", resp)

			_, _, _ = projectsApi.ProjectsDeleteProject(ctx, projectName)
		})

		It("should fail to create a project with empty name", Label("L00002"), func() {
			_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name: "",
			})
			Expect(err).To(HaveOccurred())

			GinkgoWriter.Printf("Error creating project with empty name: %v\n", err)
		})

		It("should create project with valid name - alphanumeric and hyphen", Label("L00018"), func() {
			validNames := []string{
				tools.GenerateTestProjectName("test-project"),
				tools.GenerateTestProjectName("t1"),
				tools.GenerateTestProjectName("1test2"),
				tools.GenerateTestProjectName("test"),
				tools.GenerateTestProjectName("12"),
				tools.GenerateTestProjectName("Test"),
				tools.GenerateTestProjectName("MyProject"),
				tools.GenerateTestProjectName("ABC-123"),
			}
			projectType := v1alpha1project.PRIVATE_V1alpha1ProjectType

			for _, name := range validNames {
				_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
					Name:  name,
					Type_: &projectType,
				})
				Expect(err).NotTo(HaveOccurred(), "project %s should be created successfully", name)
				_, _, _ = projectsApi.ProjectsDeleteProject(ctx, name)
			}
		})

		It("should fail to create project with invalid name", Label("L00019"), func() {
			invalidNames := []string{
				"t",        // too short (less than 2 chars)
				"-test",    // starts with hyphen
				"test 01",  // contains space
				"test%123", // contains special char
				"test*123", // contains special char
				"test~01",  // contains special char
				"1test-",   // ends with hyphen
				"aa",       // fewer than 2 distinct characters
				"11",       // fewer than 2 distinct characters
				"AAAA",     // fewer than 2 distinct characters
			}

			for _, name := range invalidNames {
				_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
					Name: name,
				})
				Expect(err).To(HaveOccurred(), "project %s should fail to create", name)
				GinkgoWriter.Printf("Invalid name '%s': error=%v\n", name, err)
			}
		})

		It("should create public project successfully", Label("L00020"), func() {
			projectName := tools.GenerateTestProjectName("public-project")
			projectType := v1alpha1project.PUBLIC_V1alpha1ProjectType

			_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name:  projectName,
				Type_: &projectType,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify project is created
			getResp, _, err := projectsApi.ProjectsGetProject(ctx, projectName)
			Expect(err).NotTo(HaveOccurred())
			Expect(*getResp.Type_).To(Equal(v1alpha1project.PUBLIC_V1alpha1ProjectType))

			_, _, _ = projectsApi.ProjectsDeleteProject(ctx, projectName)
		})

		It("should create private project successfully", Label("L00021"), func() {
			projectName := tools.GenerateTestProjectName("private-project")
			projectType := v1alpha1project.PRIVATE_V1alpha1ProjectType

			_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name:  projectName,
				Type_: &projectType,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify project is created
			getResp, _, err := projectsApi.ProjectsGetProject(ctx, projectName)
			Expect(err).NotTo(HaveOccurred())
			Expect(*getResp.Type_).To(Equal(v1alpha1project.PRIVATE_V1alpha1ProjectType))

			_, _, _ = projectsApi.ProjectsDeleteProject(ctx, projectName)
		})

		It("should fail to create duplicate project name", Label("L00022"), func() {
			projectName := tools.GenerateTestProjectName("duplicate-test")
			projectType := v1alpha1project.PRIVATE_V1alpha1ProjectType

			// Create first project
			_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name:  projectName,
				Type_: &projectType,
			})
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_, _, _ = projectsApi.ProjectsDeleteProject(ctx, projectName)
			}()

			// Try to create duplicate
			_, _, err = projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name:  projectName,
				Type_: &projectType,
			})
			Expect(err).To(HaveOccurred(), "duplicate project should fail")

			GinkgoWriter.Printf("Duplicate project error: %v\n", err)
		})

		It("should treat project names as case-sensitive", Label("L00046"), func() {
			lowerName := tools.GenerateTestProjectName("case-sensitive")
			upperName := strings.ToUpper(lowerName[:1]) + lowerName[1:]
			projectType := v1alpha1project.PRIVATE_V1alpha1ProjectType

			_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name:  lowerName,
				Type_: &projectType,
			})
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_, _, _ = projectsApi.ProjectsDeleteProject(ctx, lowerName)
			}()

			_, _, err = projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name:  upperName,
				Type_: &projectType,
			})
			Expect(err).NotTo(HaveOccurred(), "name differing only in case should be a distinct project")
			defer func() {
				_, _, _ = projectsApi.ProjectsDeleteProject(ctx, upperName)
			}()
		})

		It("should fail to create project with unspecified type", Label("L00045"), func() {
			projectName := tools.GenerateTestProjectName("unspecified-type")
			unspecifiedType := v1alpha1project.UNSPECIFIED_V1alpha1ProjectType

			_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name:  projectName,
				Type_: &unspecifiedType,
			})
			Expect(err).To(HaveOccurred(), "project with unspecified type should fail")

			GinkgoWriter.Printf("Unspecified type error: %v\n", err)
		})
	})

	Context("GetProject API", func() {
		It("should get an existing project", Label("L00003"), func() {
			projectName := tools.GenerateTestProjectName("project")
			projectType := v1alpha1project.PRIVATE_V1alpha1ProjectType
			GinkgoWriter.Printf("Get project test: %v\n", projectName)

			_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name:  projectName,
				Type_: &projectType,
			})
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_, _, _ = projectsApi.ProjectsDeleteProject(ctx, projectName)
			}()

			// Get project
			getResp, _, err := projectsApi.ProjectsGetProject(ctx, projectName)
			Expect(err).NotTo(HaveOccurred())
			Expect(getResp.Name).To(Equal(projectName))

			GinkgoWriter.Printf("Get response: name=%v\n", getResp.Name)
		})

		It("should fail to get a non-existent project", Label("L00004"), func() {
			_, _, err := projectsApi.ProjectsGetProject(ctx, "non-existent-project-xyz")
			Expect(err).To(HaveOccurred())

			GinkgoWriter.Printf("Error getting non-existent project: %v\n", err)
		})
	})

	Context("ListProjects API", func() {
		It("should list projects successfully", Label("L00005"), func() {
			projectName := tools.GenerateTestProjectName("project")
			projectType := v1alpha1project.PRIVATE_V1alpha1ProjectType
			GinkgoWriter.Printf("List projects test: %v\n", projectName)

			// Create a project first
			_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name:  projectName,
				Type_: &projectType,
			})
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				_, _, _ = projectsApi.ProjectsDeleteProject(ctx, projectName)
			}()

			// List projects
			listResp, _, err := projectsApi.ProjectsListProjects(ctx, &v1alpha1project.ProjectsApiProjectsListProjectsOpts{
				Page:     optional.NewInt32(1),
				PageSize: optional.NewInt32(10),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(listResp.Projects).NotTo(BeNil())

			// Find our project in the list
			var found bool
			for _, p := range listResp.Projects {
				if p.Name == projectName {
					found = true
