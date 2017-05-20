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

package tools

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/antihax/optional"
	v1alpha1login "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/login"
	v1alpha1user "github.com/matrixhub-ai/matrixhub/test/client/v1alpha1/user"
)

var (
	authInitOnce sync.Once
	authInitErr  error

	// Admin credentials (from environment or default)
	adminUsername string
	adminPassword string
	adminCookie   string

	// API clients for admin
	adminLoginApi *v1alpha1login.LoginApiService
	adminUsersApi *v1alpha1user.UsersApiService
)

// Environment variable names for admin credentials
const (
	EnvAdminUsername = "MATRIXHUB_ADMIN_USERNAME"
	EnvAdminPassword = "MATRIXHUB_ADMIN_PASSWORD"
)

// Default admin credentials
const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "changeme"
)

// InitAuth initializes authentication with admin login
func InitAuth() error {
	authInitOnce.Do(func() {
		// Get admin credentials from environment or use defaults
		adminUsername = os.Getenv(EnvAdminUsername)
		if adminUsername == "" {
			adminUsername = DefaultAdminUsername
		}
		adminPassword = os.Getenv(EnvAdminPassword)
		if adminPassword == "" {
			adminPassword = DefaultAdminPassword
		}

		log.Printf("Initializing auth with admin user: %s\n", adminUsername)

		// Login as admin to get cookie
		cookie, err := LoginUser(adminUsername, adminPassword)
		if err != nil {
			authInitErr = fmt.Errorf("failed to login as admin: %w", err)
			return
		}

		adminCookie = cookie
		log.Println("Admin login successful")

		// Initialize admin API clients with cookie
		initAdminClients()
	})

	return authInitErr
}

// initAdminClients initializes API clients with admin cookie
func initAdminClients() {
	baseURL := GetBaseURL()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402
			Proxy:           http.ProxyFromEnvironment,
		},
	}

	defaultHeaders := map[string]string{
		"Cookie":       adminCookie,
		"Content-Type": "application/json",
	}

	// Initialize Login API client
	loginCfg := &v1alpha1login.Configuration{
		BasePath:      baseURL,
		DefaultHeader: defaultHeaders,
		HTTPClient:    httpClient,
	}
	adminLoginApi = v1alpha1login.NewAPIClient(loginCfg).LoginApi

	// Initialize Users API client
	userCfg := &v1alpha1user.Configuration{
		BasePath:      baseURL,
		DefaultHeader: defaultHeaders,
		HTTPClient:    httpClient,
	}
	adminUsersApi = v1alpha1user.NewAPIClient(userCfg).UsersApi
