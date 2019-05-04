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
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matrixhub-ai/hfd/pkg/authenticate"
	"github.com/matrixhub-ai/hfd/pkg/repository"
	"github.com/matrixhub-ai/hfd/pkg/storage"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	domaingit "github.com/matrixhub-ai/matrixhub/internal/domain/git"
	gitmocks "github.com/matrixhub-ai/matrixhub/internal/domain/git/mocks"
	modelmocks "github.com/matrixhub-ai/matrixhub/internal/domain/model/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/robot"
	"github.com/matrixhub-ai/matrixhub/internal/infra/authcodec"
)

func TestHandleCommitRoutesByRepositoryType(t *testing.T) {
	repoTypes := []struct {
		repoType    string
		storageName string
	}{
		{repoType: "models", storageName: "test-user/shared-name"},
		{repoType: "datasets", storageName: "datasets/test-user/shared-name"},
		{repoType: "spaces", storageName: "spaces/test-user/shared-name"},
	}
	outcomes := []struct {
		name       string
		commitHash string
		err        error
		wantStatus int
	}{
		{name: "success", commitHash: "commit-sha", wantStatus: http.StatusOK},
		{name: "commit error", err: errors.New("commit failed"), wantStatus: http.StatusInternalServerError},
	}

	for _, rt := range repoTypes {
		for _, outcome := range outcomes {
			t.Run(rt.repoType+"/"+outcome.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				modelService := modelmocks.NewMockIModelService(ctrl)
				gitRepo := gitmocks.NewMockIGitRepo(ctrl)
				commit := &domaingit.Commit{
					Message:     "Add readme\n\nTest description",
					AuthorName:  "HuggingFace",
					AuthorEmail: "hf@users.noreply.huggingface.co",
				}
				ops := []domaingit.CommitOperation{{
					Type:    domaingit.CommitOperationAdd,
					Path:    "README.md",
					Content: []byte("readme content\n"),
				}}
				if rt.repoType == "models" {
