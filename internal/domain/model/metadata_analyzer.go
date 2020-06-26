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

package model

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math"
	"strings"

	"github.com/matrixhub-ai/hfd/pkg/hf"

	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
)

// ClassifiedTag represents model metadata tags after applying model-domain rules.
type ClassifiedTag struct {
	Name     string
	Category string
}

// RepoMetadata is the structured model metadata extracted from raw repo files.
type RepoMetadata struct {
	ReadmeContent  string
	Size           int64
	ParameterCount int64
	Tags           []ClassifiedTag
}

type safetensorsIndex struct {
	Metadata struct {
		TotalSize int64 `json:"total_size"`
	} `json:"metadata"`
}

type safetensorsTensor struct {
	Shape []int64 `json:"shape"`
}

// AnalyzeRepoMetadata converts raw repo files into model-domain metadata.
// The git repo layer only loads bytes from tracked files; all model-specific
// interpretation (label classification, parameter_count inference) stays here.
func AnalyzeRepoMetadata(files *git.RepoMetadataFiles) (*RepoMetadata, error) {
	metadata := &RepoMetadata{Size: files.Size}
	var tags []ClassifiedTag

	if len(files.ReadmeContent) > 0 {
		metadata.ReadmeContent = string(files.ReadmeContent)
		if readme, err := hf.ParseReadme(bytes.NewReader(files.ReadmeContent)); err == nil {
			card := readme.Card
			if card.PipelineTag != "" {
				tags = append(tags, ClassifiedTag{Name: card.PipelineTag, Category: "task"})
			}
			if card.LibraryName != "" {
				tags = append(tags, ClassifiedTag{Name: card.LibraryName, Category: "library"})
			}
			for _, lang := range card.Language {
				tags = append(tags, ClassifiedTag{Name: lang, Category: "language"})
			}
			for _, lic := range card.License {
				tags = append(tags, ClassifiedTag{Name: lic, Category: "license"})
			}
			for _, tag := range card.Tags {
				tags = append(tags, ClassifiedTag{Name: tag, Category: "other"})
			}
		}
	}

	// parameter_count inference only uses config.json and safetensors metadata.
	// README.md is intentionally excluded here because it is descriptive text,
	// not structured machine metadata.
	//
	// Default to fp16/bf16-sized parameters when config.json does not expose
	// enough dtype/quantization information. This is a pragmatic default for
	// modern Hugging Face model repos.
	parameterBytes := int64(2)
	if len(files.ConfigJSON) > 0 {
		if cfg, err := hf.ParseConfigData(bytes.NewReader(files.ConfigJSON)); err == nil {
			if cfg.ModelType != "" {
				tags = append(tags, ClassifiedTag{Name: cfg.ModelType, Category: "other"})
			}
			if cfg.QuantizationConfig != nil && cfg.QuantizationConfig.QuantMethod != "" {
				tags = append(tags, ClassifiedTag{Name: cfg.QuantizationConfig.QuantMethod, Category: "other"})
			}
		}
		if inferred := inferParameterBytes(files.ConfigJSON); inferred > 0 {
			parameterBytes = inferred
		}
	}

	if len(files.SafetensorsIndexJSON) > 0 {
		// For most sharded HF models, model.safetensors.index.json metadata.total_size
		// gives total tensor bytes. Dividing by bytes-per-parameter gives a cheap,
		// usually-good parameter_count estimate without scanning all shards.
		var index safetensorsIndex
		if err := json.Unmarshal(files.SafetensorsIndexJSON, &index); err == nil && index.Metadata.TotalSize > 0 {
			metadata.ParameterCount = estimateParameterCount(index.Metadata.TotalSize, parameterBytes)
		}
	}
	if metadata.ParameterCount == 0 {
		// Header scans are exact but only cover files whose weights were readable.
		// Files that fell back to the size recorded in their LFS pointer are
		// estimated from that size instead, using the same arithmetic as the index
		// path above. The two sets are disjoint — the git layer records each file
		// as one or the other — so counting both is what keeps a partly-fetched
		// sharded model from reporting only the shards it could read.
		exact := countSafetensorsParameters(files.SafetensorsFiles)
		estimated := estimateParameterCount(sumSafetensorsSizes(files.SafetensorsSizes), parameterBytes)
		metadata.ParameterCount = addParameterCounts(exact, estimated)
	}

	metadata.Tags = deduplicateClassifiedTags(tags)
	return metadata, nil
}

func deduplicateClassifiedTags(tags []ClassifiedTag) []ClassifiedTag {
	type tagKey struct {
		Name     string
		Category string
	}

	seen := make(map[tagKey]struct{})
	result := make([]ClassifiedTag, 0, len(tags))

	for _, tag := range tags {
		if tag.Name == "" {
