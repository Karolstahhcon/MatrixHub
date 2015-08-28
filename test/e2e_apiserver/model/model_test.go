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

package model_test

import (
	"context"
	"strings"

	v1alpha1model "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/model"
	v1alpha1project "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/project"
	"github.com/matrixhub-ai/matrixhub/test/tools"

	"github.com/antihax/optional"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func boolPtr(b bool) *bool { return &b }

func createGitModelFixture(
	ctx context.Context,
	modelsAPI *v1alpha1model.ModelsApiService,
	projectsAPI *v1alpha1project.ProjectsApiService,
) (string, string) {
	projectName := tools.GenerateTestProjectName("model-git")
	projectType := v1alpha1project.PRIVATE_V1alpha1ProjectType
	_, _, err := projectsAPI.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
		Name:  projectName,
		Type_: &projectType,
	})
	Expect(err).NotTo(HaveOccurred())

	modelName := tools.GenerateTestModelName("git-model")
	_, _, err = modelsAPI.ModelsCreateModel(ctx, v1alpha1model.V1alpha1CreateModelRequest{
		Project: projectName,
		Name:    modelName,
	})
	Expect(err).NotTo(HaveOccurred())

	return projectName, modelName
}

func cleanupGitModelFixture(
	ctx context.Context,
	modelsAPI *v1alpha1model.ModelsApiService,
	projectsAPI *v1alpha1project.ProjectsApiService,
	projectName, modelName string,
) {
	if modelName != "" {
		_, _, _ = modelsAPI.ModelsDeleteModel(ctx, projectName, modelName)
	}
	if projectName != "" {
		_, _, _ = projectsAPI.ProjectsDeleteProject(ctx, projectName)
	}
}

var _ = Describe("Model", Label("model"), func() {
	var (
		ctx         context.Context
		modelsApi   *v1alpha1model.ModelsApiService
		projectsApi *v1alpha1project.ProjectsApiService
	)

	BeforeEach(func() {
		ctx = context.Background()
		modelsApi = tools.GetV1alpha1ModelsApi()
		projectsApi = tools.GetV1alpha1ProjectsApi()
	})

	// ═══════════════════════════════════════════════════════════
	// 1. Label APIs
	// ═══════════════════════════════════════════════════════════
	Context("ListModelTaskLabels API", func() {
		It("should list task labels successfully", Label("M00001"), func() {
			resp, _, err := modelsApi.ModelsListModelTaskLabels(ctx)
			Expect(err).NotTo(HaveOccurred())

			GinkgoWriter.Printf("Task labels: found %d labels\n", len(resp.Items))
		})
	})

	Context("ListModelFrameLabels API", func() {
		It("should list framework/library labels successfully", Label("M00002"), func() {
			resp, _, err := modelsApi.ModelsListModelFrameLabels(ctx)
			Expect(err).NotTo(HaveOccurred())

			GinkgoWriter.Printf("Frame labels: found %d labels\n", len(resp.Items))
		})
	})

	// ═══════════════════════════════════════════════════════════
	// 2. CRUD APIs (CreateModel / GetModel / ListModels / DeleteModel)
	// ═══════════════════════════════════════════════════════════
	Context("CRUD APIs", func() {
		var projectName string

		BeforeEach(func() {
			projectName = tools.GenerateTestProjectName("model-crud")
			projectType := v1alpha1project.PRIVATE_V1alpha1ProjectType
			_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
				Name:  projectName,
				Type_: &projectType,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_, _, _ = projectsApi.ProjectsDeleteProject(ctx, projectName)
		})

		// --- CreateModel ---
		It("should create a model successfully", Label("M00003", "smoke"), func() {
			modelName := tools.GenerateTestModelName("model")

			_, _, err := modelsApi.ModelsCreateModel(ctx, v1alpha1model.V1alpha1CreateModelRequest{
				Project: projectName,
				Name:    modelName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify model is created via GetModel
			getResp, _, err := modelsApi.ModelsGetModel(ctx, projectName, modelName)
			Expect(err).NotTo(HaveOccurred())
			Expect(getResp.Name).To(Equal(modelName))
			Expect(getResp.Project).To(Equal(projectName))

			GinkgoWriter.Printf("Created model: name=%v, project=%v\n", getResp.Name, getResp.Project)
		})

		It("should fail to create duplicate model", Label("M00004"), func() {
			modelName := tools.GenerateTestModelName("dup-model")

			_, _, err := modelsApi.ModelsCreateModel(ctx, v1alpha1model.V1alpha1CreateModelRequest{
				Project: projectName,
				Name:    modelName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Try to create duplicate
			_, _, err = modelsApi.ModelsCreateModel(ctx, v1alpha1model.V1alpha1CreateModelRequest{
				Project: projectName,
				Name:    modelName,
			})
			Expect(err).To(HaveOccurred(), "duplicate model should fail")

			GinkgoWriter.Printf("Duplicate model error: %v\n", err)
		})

		It("should fail to create a model with empty project", Label("M00005"), func() {
			modelName := tools.GenerateTestModelName("model")
			_, _, err := modelsApi.ModelsCreateModel(ctx, v1alpha1model.V1alpha1CreateModelRequest{
				Project: "",
				Name:    modelName,
			})
			Expect(err).To(HaveOccurred())
		})

		It("should fail to create a model with empty name", Label("M00006"), func() {
			_, _, err := modelsApi.ModelsCreateModel(ctx, v1alpha1model.V1alpha1CreateModelRequest{
				Project: projectName,
				Name:    "",
			})
			Expect(err).To(HaveOccurred())
		})

		It("should fail to create a model in non-existent project", Label("M00007"), func() {
			modelName := tools.GenerateTestModelName("model")
			_, _, err := modelsApi.ModelsCreateModel(ctx, v1alpha1model.V1alpha1CreateModelRequest{
				Project: "non-existent-project-xyz",
				Name:    modelName,
			})
			Expect(err).To(HaveOccurred())
		})

		It("should reject invalid model names", Label("M00058"), func() {
			invalidNames := []string{
				"-leading-dash",
				"contains space",
				"contains/slash",
				"models",
				strings.Repeat("a", 252),
			}

			for _, modelName := range invalidNames {
				_, _, err := modelsApi.ModelsCreateModel(ctx, v1alpha1model.V1alpha1CreateModelRequest{
					Project: projectName,
					Name:    modelName,
				})
				Expect(err).To(HaveOccurred(), "model name %q should fail validation", modelName)
			}
		})

		// --- GetModel ---
		It("should get an existing model with all expected fields", Label("M00008"), func() {
			modelName := tools.GenerateTestModelName("get-model")

			_, _, err := modelsApi.ModelsCreateModel(ctx, v1alpha1model.V1alpha1CreateModelRequest{
				Project: projectName,
				Name:    modelName,
			})
			Expect(err).NotTo(HaveOccurred())

			getResp, _, err := modelsApi.ModelsGetModel(ctx, projectName, modelName)
			Expect(err).NotTo(HaveOccurred())
			Expect(getResp.Name).To(Equal(modelName))
			Expect(getResp.Project).To(Equal(projectName))
			Expect(getResp.Id).NotTo(BeZero())
			Expect(getResp.CreatedAt).NotTo(BeEmpty())

			GinkgoWriter.Printf("GetModel: id=%v, name=%v, project=%v\n", getResp.Id, getResp.Name, getResp.Project)
		})

		It("should return default fields for a newly created model", Label("M00059"), func() {
			modelName := tools.GenerateTestModelName("defaults")
			_, _, err := modelsApi.ModelsCreateModel(ctx, v1alpha1model.V1alpha1CreateModelRequest{
				Project: projectName,
				Name:    modelName,
			})
			Expect(err).NotTo(HaveOccurred())

			getResp, _, err := modelsApi.ModelsGetModel(ctx, projectName, modelName)
			Expect(err).NotTo(HaveOccurred())
			Expect(getResp.DefaultBranch).To(Equal("main"))
			Expect(getResp.CloneUrls).NotTo(BeNil())
			Expect(getResp.Labels).To(BeEmpty())
			Expect(getResp.ReadmeContent).To(BeEmpty())
			Expect(getResp.Size).To(Equal("0"))
			Expect(getResp.ParameterCount).To(Equal("0"))
			Expect(getResp.Popular).To(BeFalse())
		})

		It("should fail to get a non-existent model", Label("M00009"), func() {
			_, _, err := modelsApi.ModelsGetModel(ctx, projectName, "non-existent-model-xyz")
			Expect(err).To(HaveOccurred())
		})

		It("should fail to get a model from non-existent project", Label("M00010"), func() {
			_, _, err := modelsApi.ModelsGetModel(ctx, "nonexistent", "nonexistent-model")
			Expect(err).To(HaveOccurred())
