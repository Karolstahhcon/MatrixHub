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
