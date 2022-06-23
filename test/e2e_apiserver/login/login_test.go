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

package login_test

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"

	v1alpha1login "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/login"
	"github.com/matrixhub-ai/matrixhub/test/tools"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// newLoginApi creates a fresh Login API client with no session cookie.
func newLoginApi() *v1alpha1login.LoginApiService {
	baseURL := tools.GetBaseURL()
	cfg := &v1alpha1login.Configuration{
		BasePath: baseURL,
		DefaultHeader: map[string]string{
			"Content-Type": "application/json",
		},
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402
				Proxy:           http.ProxyFromEnvironment,
			},
		},
	}
	return v1alpha1login.NewAPIClient(cfg).LoginApi
}

// newLoginApiWithCookie creates a Login API client pre-loaded with the given session cookie.
func newLoginApiWithCookie(cookie string) *v1alpha1login.LoginApiService {
	baseURL := tools.GetBaseURL()
	cfg := &v1alpha1login.Configuration{
		BasePath: baseURL,
		DefaultHeader: map[string]string{
			"Content-Type": "application/json",
			"Cookie":       cookie,
		},
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402
				Proxy:           http.ProxyFromEnvironment,
			},
		},
	}
	return v1alpha1login.NewAPIClient(cfg).LoginApi
}

// extractSetCookie parses the name=value part from a Set-Cookie header value.
func extractSetCookie(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	raw := resp.Header.Get("Set-Cookie")
	if raw == "" {
		return ""
	}
	if idx := strings.Index(raw, ";"); idx > 0 {
		return strings.TrimSpace(raw[:idx])
	}
	return strings.TrimSpace(raw)
}

var _ = Describe("Login", Label("login"), func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = context.Background()
	})

	// ═══════════════════════════════════════════════════════════
	// 1. Login API
	// ═══════════════════════════════════════════════════════════
