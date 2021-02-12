// // Copyright The MatrixHub Authors.
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// //     http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.
package model_test

//
// import (
// 	"context"
// 	"crypto/tls"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"strings"
// 	"time"
//
// 	v1alpha1model "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/model"
// 	v1alpha1project "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/project"
// 	"github.com/matrixhub-ai/matrixhub/test/tools"
//
// 	"github.com/antihax/optional"
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// )
//
// // Test fixture: README.md with YAML front matter containing metadata tags
// const testReadmeContent = `---
// library_name: transformers
// license: apache-2.0
// pipeline_tag: text-generation
// language:
// - en
// - zh
// tags:
// - pytorch
// ---
//
// # Test Model
//
// This is a test model for metadata sync e2e testing.
// `
//
// // Test fixture: config.json with model architecture info
// const testConfigJSON = `{
//   "architectures": ["LlamaForCausalLM"],
//   "model_type": "llama",
//   "text_config": {
//     "dtype": "bfloat16"
//   },
//   "torch_dtype": "bfloat16"
// }`
//
// // Test fixture: model.safetensors.index.json with total_size for parameter count estimation
// // total_size=26000000000, dtype=bfloat16(2 bytes) => parameter_count=13000000000 (~13B)
// const testSafetensorsIndexJSON = `{
//   "metadata": {
//     "total_size": 26000000000
//   },
//   "weight_map": {
//     "model.embed_tokens.weight": "model-00001-of-00002.safetensors",
//     "model.layers.0.self_attn.q_proj.weight": "model-00001-of-00002.safetensors",
//     "lm_head.weight": "model-00002-of-00002.safetensors"
//   }
// }`
//
// // commitOperation is the NDJSON line format for the HF commit API.
// type commitOperation struct {
// 	Key   string      `json:"key"`
// 	Value interface{} `json:"value"`
// }
//
// type commitHeader struct {
// 	Summary     string `json:"summary"`
// 	Description string `json:"description"`
// }
//
// type commitFile struct {
// 	Path     string `json:"path"`
// 	Content  string `json:"content"`
// 	Encoding string `json:"encoding"`
// }
//
// type commitResponse struct {
// 	CommitURL     string `json:"commitUrl"`
// 	CommitOid     string `json:"commitOid"`
// 	CommitMessage string `json:"commitMessage"`
// }
//
// // buildNDJSON constructs the NDJSON payload for the HF commit API.
// func buildNDJSON(summary string, files map[string]string) string {
// 	var lines []string
//
// 	// Header line
// 	header := commitOperation{
// 		Key:   "header",
// 		Value: commitHeader{Summary: summary, Description: "e2e metadata sync test"},
// 	}
// 	b, _ := json.Marshal(header)
// 	lines = append(lines, string(b))
//
// 	// File lines
// 	for path, content := range files {
// 		fileOp := commitOperation{
// 			Key:   "file",
// 			Value: commitFile{Path: path, Content: content, Encoding: "utf-8"},
// 		}
// 		b, _ := json.Marshal(fileOp)
// 		lines = append(lines, string(b))
// 	}
//
// 	return strings.Join(lines, "\n")
// }
//
// // hfCommit sends an NDJSON commit to the HF-compatible commit API.
// func hfCommit(baseURL, project, model, revision, payload string) (*commitResponse, error) {
// 	url := fmt.Sprintf("%s/api/models/%s/%s/commit/%s", baseURL, project, model, revision)
//
// 	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(payload))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/x-ndjson")
//
// 	// Use admin cookie for auth
// 	cookie := tools.GetAdminCookie()
// 	if cookie != "" {
// 		req.Header.Set("Cookie", cookie)
// 	}
//
// 	client := &http.Client{
// 		Transport: &http.Transport{
// 			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402
// 		},
// 	}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("commit request failed: %w", err)
// 	}
// 	defer resp.Body.Close()
//
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response: %w", err)
// 	}
//
// 	if resp.StatusCode >= 300 {
// 		return nil, fmt.Errorf("commit failed: HTTP %d: %s", resp.StatusCode, string(body))
// 	}
//
// 	var commitResp commitResponse
// 	if err := json.Unmarshal(body, &commitResp); err != nil {
// 		return nil, fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
// 	}
//
// 	return &commitResp, nil
// }
//
// // extractLabelNames extracts label name list from model labels.
// func extractLabelNames(labels []v1alpha1model.MatrixhubV1alpha1Label) []string {
// 	names := make([]string, len(labels))
// 	for i, l := range labels {
// 		names[i] = l.Name
// 	}
// 	return names
// }
//
// var _ = Describe("Model Metadata Sync", Label("model", "metadata-sync"), func() {
// 	var (
// 		ctx         context.Context
// 		modelsApi   *v1alpha1model.ModelsApiService
// 		projectsApi *v1alpha1project.ProjectsApiService
// 		projectName string
// 		modelName   string
// 		baseURL     string
// 	)
//
// 	BeforeEach(func() {
// 		ctx = context.Background()
// 		modelsApi = tools.GetV1alpha1ModelsApi()
// 		projectsApi = tools.GetV1alpha1ProjectsApi()
// 		baseURL = tools.GetBaseURL()
//
// 		// Create project
// 		projectName = tools.GenerateTestProjectName("sync")
// 		_, _, err := projectsApi.ProjectsCreateProject(ctx, v1alpha1project.V1alpha1CreateProjectRequest{
// 			Name: projectName,
// 		})
// 		Expect(err).NotTo(HaveOccurred())
//
// 		// Create model
// 		modelName = tools.GenerateTestModelName("sync")
// 		_, _, err = modelsApi.ModelsCreateModel(ctx, v1alpha1model.V1alpha1CreateModelRequest{
// 			Project: projectName,
// 			Name:    modelName,
// 		})
// 		Expect(err).NotTo(HaveOccurred())
// 	})
//
// 	AfterEach(func() {
// 		_, _, _ = modelsApi.ModelsDeleteModel(ctx, projectName, modelName)
// 		_, _, _ = projectsApi.ProjectsDeleteProject(ctx, projectName)
// 	})
//
// 	// ═══════════════════════════════════════════════════════════
// 	// HF Commit API basic
// 	// ═══════════════════════════════════════════════════════════
// 	Context("HF Commit API", func() {
// 		It("should commit files via NDJSON API", Label("M00039"), func() {
// 			payload := buildNDJSON("add README.md", map[string]string{
// 				"README.md": "# Hello\nThis is a test.",
// 			})
//
// 			commitResp, err := hfCommit(baseURL, projectName, modelName, "main", payload)
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(commitResp.CommitOid).NotTo(BeEmpty())
//
// 			GinkgoWriter.Printf("Commit OID: %s\n", commitResp.CommitOid)
// 		})
// 	})
//
// 	// ═══════════════════════════════════════════════════════════
// 	// Full metadata sync + all model API verification
// 	// ═══════════════════════════════════════════════════════════
// 	Context("Metadata sync end-to-end with all model APIs", func() {
// 		It("should sync all metadata and reflect in every model API", Label("M00040"), func() {
// 			By("recording initial model state")
// 			initialModel, _, err := modelsApi.ModelsGetModel(ctx, projectName, modelName)
// 			Expect(err).NotTo(HaveOccurred())
// 			initialReadme := initialModel.ReadmeContent
// 			Expect(len(initialModel.Labels)).To(Equal(0), "new model should have no labels")
//
// 			By("committing README.md + config.json + safetensors index via HF commit API")
// 			payload := buildNDJSON("add model metadata files", map[string]string{
// 				"README.md":                    testReadmeContent,
// 				"config.json":                  testConfigJSON,
// 				"model.safetensors.index.json": testSafetensorsIndexJSON,
// 			})
//
// 			commitResp, err := hfCommit(baseURL, projectName, modelName, "main", payload)
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(commitResp.CommitOid).NotTo(BeEmpty())
// 			commitOid := commitResp.CommitOid
// 			GinkgoWriter.Printf("Commit OID: %s\n", commitOid)
//
// 			By("waiting for postReceiveHook to complete metadata sync")
// 			time.Sleep(3 * time.Second)
//
// 			// ── GetModel: verify metadata fields ───────────────
// 			By("verifying GetModel returns synced metadata")
// 			updatedModel, _, err := modelsApi.ModelsGetModel(ctx, projectName, modelName)
// 			Expect(err).NotTo(HaveOccurred())
//
// 			// readmeContent should be updated
// 			Expect(updatedModel.ReadmeContent).NotTo(Equal(initialReadme))
// 			Expect(updatedModel.ReadmeContent).To(ContainSubstring("# Test Model"))
// 			GinkgoWriter.Printf("GetModel: readmeContent length=%d\n", len(updatedModel.ReadmeContent))
//
// 			// parameterCount: 26000000000 / 2 = 13000000000
// 			Expect(updatedModel.ParameterCount).To(Equal("13000000000"))
// 			GinkgoWriter.Printf("GetModel: parameterCount=%s\n", updatedModel.ParameterCount)
//
// 			// labels should be populated from front matter + config
// 			Expect(len(updatedModel.Labels)).To(BeNumerically(">", 0))
// 			labelNames := extractLabelNames(updatedModel.Labels)
// 			GinkgoWriter.Printf("GetModel: labels=%v\n", labelNames)
//
// 			Expect(labelNames).To(ContainElement("text-generation")) // pipeline_tag -> task
// 			Expect(labelNames).To(ContainElement("transformers"))    // library_name -> library
// 			Expect(labelNames).To(ContainElement("apache-2.0"))      // license
// 			Expect(labelNames).To(ContainElement("llama"))           // model_type from config.json
//
// 			// ── ListModels: verify synced model appears with labels ─
// 			By("verifying ListModels returns the model with labels")
// 			listResp, _, err := modelsApi.ModelsListModels(ctx, &v1alpha1model.ModelsApiModelsListModelsOpts{
// 				Project:  optional.NewString(projectName),
// 				Page:     optional.NewInt32(1),
// 				PageSize: optional.NewInt32(10),
// 			})
// 			Expect(err).NotTo(HaveOccurred())
