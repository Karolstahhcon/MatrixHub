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

package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/matrixhub-ai/hfd/pkg/authenticate"

	"github.com/matrixhub-ai/matrixhub/internal/domain/robot"
	"github.com/matrixhub-ai/matrixhub/internal/domain/user"
	"github.com/matrixhub-ai/matrixhub/internal/infra/utils"
)

func TestGitAuthRejectsInvalidCredentials(t *testing.T) {
	expired := time.Now().Add(-time.Hour)
	for _, test := range []struct {
		name   string
		token  string
		access *user.AccessToken
		robot  *robot.Robot
		err    error
		status int
	}{
		{"unrecognized", "forged", nil, nil, nil, http.StatusOK},
		{"deleted", utils.TokenPrefix + "deleted", &user.AccessToken{}, nil, nil, http.StatusUnauthorized},
		{"expired", utils.TokenPrefix + "expired", &user.AccessToken{Enabled: true, ExpireAt: &expired}, nil, nil, http.StatusUnauthorized},
		{"missing robot", utils.RobotTokenPrefix + "missing", nil, nil, nil, http.StatusUnauthorized},
