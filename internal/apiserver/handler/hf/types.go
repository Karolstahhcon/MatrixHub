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
