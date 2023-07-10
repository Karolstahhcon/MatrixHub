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

package backend

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/matrixhub-ai/hfd/pkg/authenticate"
	"github.com/matrixhub-ai/hfd/pkg/permission"
	"github.com/matrixhub-ai/hfd/pkg/receive"
	"github.com/matrixhub-ai/hfd/pkg/repository"

	"github.com/matrixhub-ai/matrixhub/internal/infra/log"
)

type repoInformation struct {
	RepoType string

	FullName  string
	Namespace string
	Name      string
}

var repoTypePrefixes = []struct {
	prefix   string
	repoType string
}{
	{"datasets/", "datasets"},
	{"spaces/", "spaces"},
}

// getRepoInformation returns the repository information extracted from the request, including repo type, storage path, namespace, and name.
func getRepoInformation(r *http.Request) repoInformation {
	vars := mux.Vars(r)
	name := vars["repo"]

	repoType := "models"
	namespacedName := name
	for _, p := range repoTypePrefixes {
		if strings.HasPrefix(name, p.prefix) {
			repoType = p.repoType
			namespacedName = strings.TrimPrefix(name, p.prefix)
			break
		}
	}

	namespace, repoName, ok := strings.Cut(namespacedName, "/")
	if !ok {
		return repoInformation{}
	}

	return repoInformation{
		RepoType:  repoType,
		FullName:  name,
		Namespace: namespace,
		Name:      repoName,
	}
}

// gitProtocolEnv returns a GIT_PROTOCOL environment variable derived from the
// request's Git-Protocol header if the value is present and valid, or nil otherwise.
func gitProtocolEnv(r *http.Request) []string {
