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

package hf

import (
	"encoding/json"

	"github.com/matrixhub-ai/hfd/pkg/repository"
)

// whoamiResponse represents the response for the /api/whoami-v2 endpoint.
type whoamiResponse struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Fullname      string   `json:"fullname"`
	Email         string   `json:"email,omitempty"`
	EmailVerified bool     `json:"emailVerified"`
	IsPro         bool     `json:"isPro"`
	CanPay        bool     `json:"canPay"`
	AvatarURL     string   `json:"avatarUrl,omitempty"`
	Orgs          []any    `json:"orgs"`
	Auth          authInfo `json:"auth"`
}

// authInfo represents the auth section of the whoami response.
type authInfo struct {
	AccessToken accessToken `json:"accessToken"`
}

// accessToken represents the access token info in the whoami response.
type accessToken struct {
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`
}

// treeEntry is the API response type for a tree entry, with JSON annotations.
type treeEntry struct {
	OID        string               `json:"oid"`
	Path       string               `json:"path"`
	Type       repository.EntryType `json:"type"`
	Size       int64                `json:"size"`
	LFS        *lfsPointer          `json:"lfs,omitempty"`
	LastCommit *treeLastCommit      `json:"lastCommit,omitempty"`
}

// lfsPointer is the API response type for an LFS pointer, with JSON annotations.
type lfsPointer struct {
	OID         string `json:"oid"`
	Size        int64  `json:"size"`
	PointerSize int64  `json:"pointerSize"`
}

// treeLastCommit is the API response type for the last commit of a tree entry, with JSON annotations.
type treeLastCommit struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}

// treeSize represents the response for the Get folder size API.
type treeSize struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

// repoListItem represents a repository item in the list response for models, datasets, and spaces.
type repoListItem struct {
	RepoID        string   `json:"id"`
	Likes         int      `json:"likes"`
	TrendingScore int      `json:"trendingScore"`
	Private       bool     `json:"private"`
	Downloads     int      `json:"downloads"`
	Tags          []string `json:"tags,omitempty"`
	PipelineTag   string   `json:"pipeline_tag,omitempty"`
	LibraryName   string   `json:"library_name,omitempty"`
	ModelID       string   `json:"modelId,omitempty"`
}

// repoInfo represents the info response for HuggingFace API
type repoInfo struct {
	ID           string    `json:"id"`
	ModelID      string    `json:"modelId,omitempty"`
	SHA          string    `json:"sha"`
	Private      bool      `json:"private"`
	Disabled     bool      `json:"disabled"`
	Gated        bool      `json:"gated"`
	Downloads    int       `json:"downloads"`
	Likes        int       `json:"likes"`
