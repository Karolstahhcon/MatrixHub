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
