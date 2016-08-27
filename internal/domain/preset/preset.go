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

package preset

import "context"

type RegistrySpec struct {
	ID          int
	Name        string
	Description string
	URL         string
}

type ProjectSpec struct {
	Name         string
	Organization string
	Registry     string
}

type ModelSpec struct {
	Project        string
	Name           string
	Recommended    bool
	Size           int64
	ParameterCount int64
	Labels         []LabelSpec
}

type LabelSpec struct {
	Category string
	Name     string
}

type Manifest struct {
	Registries []RegistrySpec
	Projects   []ProjectSpec
	Models     []ModelSpec
}

// DefaultManifest returns independent seed data; upgrades do not backfill it.
// Sizes include all upstream files; parameter counts come from safetensors.total.
func DefaultManifest() *Manifest {
	return &Manifest{
		Registries: []RegistrySpec{
			{ID: 1, Name: "huggingface", Description: "Official Hugging Face Hub", URL: "https://huggingface.co"},
			{ID: 2, Name: "hf-mirror", Description: "Hugging Face community mirror", URL: "https://hf-mirror.com"},
			{ID: 3, Name: "hf.m.daocloud", Description: "DaoCloud Hugging Face mirror", URL: "https://hf.m.daocloud.io"},
		},
		Projects: []ProjectSpec{
			{Name: "deepseek-ai", Organization: "deepseek-ai", Registry: "hf.m.daocloud"},
			{Name: "Qwen", Organization: "Qwen", Registry: "hf.m.daocloud"},
			{Name: "MiniMaxAI", Organization: "MiniMaxAI", Registry: "hf.m.daocloud"},
			{Name: "zai-org", Organization: "zai-org", Registry: "hf.m.daocloud"},
		},
		Models: []ModelSpec{
			{
				Project: "deepseek-ai", Name: "DeepSeek-V4-Pro", Recommended: true,
				Size: 864739856856, ParameterCount: 1598839674782,
				Labels: []LabelSpec{
					{Category: "library", Name: "transformers"},
					{Category: "license", Name: "mit"},
					{Category: "other", Name: "deepseek_v4"},
					{Category: "other", Name: "fp8"},
				},
			},
			{
				Project: "deepseek-ai", Name: "DeepSeek-V4-Flash", Recommended: true,
				Size: 159630041626, ParameterCount: 290944616402,
				Labels: []LabelSpec{
					{Category: "library", Name: "transformers"},
					{Category: "license", Name: "mit"},
					{Category: "other", Name: "deepseek_v4"},
					{Category: "other", Name: "fp8"},
				},
			},
			{
				Project: "deepseek-ai", Name: "DeepSeek-V4-Flash-Base", Recommended: true,
				Size: 294692678949, ParameterCount: 292021347282,
				Labels: []LabelSpec{
					{Category: "other", Name: "deepseek_v4"},
					{Category: "other", Name: "fp8"},
				},
			},
			{
				Project: "deepseek-ai", Name: "DeepSeek-OCR", Recommended: true,
				Size: 6684381812, ParameterCount: 3336106240,
				Labels: []LabelSpec{
					{Category: "library", Name: "transformers"},
