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

package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/matrixhub-ai/matrixhub/api/go/v1alpha1"
	"github.com/matrixhub-ai/matrixhub/internal/domain/dataset"
	datasetmocks "github.com/matrixhub-ai/matrixhub/internal/domain/dataset/mocks"
	"github.com/matrixhub-ai/matrixhub/internal/domain/git"
	"github.com/matrixhub-ai/matrixhub/internal/domain/model"
)

const (
	datasetHandlerTestProject = "demo"
	datasetHandlerTestName    = "imdb"
)

func datasetHandlerTestSetup(t *testing.T) (*DatasetHandler, *datasetmocks.MockIDatasetService) {
	t.Helper()
	svc := datasetmocks.NewMockIDatasetService(gomock.NewController(t))
	h, ok := NewDatasetHandler(svc).(*DatasetHandler)
	require.True(t, ok)
	return h, svc
}

func datasetHandlerTestDataset() *dataset.Dataset {
	return &dataset.Dataset{
		ID:            11,
		Name:          datasetHandlerTestName,
		ProjectID:     3,
		ProjectName:   datasetHandlerTestProject,
		DefaultBranch: "main",
		NumRows:       "50000",
		Size:          84_125_825,
		ReadmeContent: "# IMDB",
		IsPopular:     true,
		Labels: []model.Label{
			{ID: 1, Name: "text-classification", Category: "task", Scope: "dataset", CreatedAt: handlerTestCreatedAt, UpdatedAt: handlerTestUpdatedAt},
			{ID: 2, Name: "pandas", Category: "library", Scope: "dataset", CreatedAt: handlerTestCreatedAt, UpdatedAt: handlerTestUpdatedAt},
		},
		CreatedAt: handlerTestCreatedAt,
		UpdatedAt: handlerTestUpdatedAt,
	}
}

func datasetHandlerTestDatasetProto() *v1alpha1.Dataset {
	return &v1alpha1.Dataset{
		Id:            11,
		Project:       datasetHandlerTestProject,
		Name:          datasetHandlerTestName,
		DefaultBranch: "main",
		NumRows:       "50000",
		Labels: []*v1alpha1.Label{
			{Id: 1, Name: "text-classification", Category: v1alpha1.Category_TASK, CreatedAt: handlerTestCreatedAtText, UpdatedAt: handlerTestUpdatedAtText},
			{Id: 2, Name: "pandas", Category: v1alpha1.Category_LIBRARY, CreatedAt: handlerTestCreatedAtText, UpdatedAt: handlerTestUpdatedAtText},
		},
		Size:          "84125825",
		CloneUrls:     &v1alpha1.CloneUrls{},
		ReadmeContent: "# IMDB",
		CreatedAt:     handlerTestCreatedAtText,
		UpdatedAt:     handlerTestUpdatedAtText,
	}
}

func TestNewDatasetHandler(t *testing.T) {
	svc := datasetmocks.NewMockIDatasetService(gomock.NewController(t))

	h, ok := NewDatasetHandler(svc).(*DatasetHandler)
	require.True(t, ok)
	require.Same(t, svc, h.ds)
}

func TestDatasetToProto(t *testing.T) {
	t.Run("maps every exposed field", func(t *testing.T) {
		requireProtoEqual(t, datasetHandlerTestDatasetProto(), datasetToProto(datasetHandlerTestDataset()))
	})

	t.Run("keeps optional fields empty", func(t *testing.T) {
		got := datasetToProto(&dataset.Dataset{
			ID:          12,
			Name:        "bare",
			ProjectName: datasetHandlerTestProject,
			CreatedAt:   handlerTestCreatedAt,
			UpdatedAt:   handlerTestUpdatedAt,
		})

		requireProtoEqual(t, &v1alpha1.Dataset{
			Id:        12,
			Project:   datasetHandlerTestProject,
			Name:      "bare",
			Size:      "0",
			CloneUrls: &v1alpha1.CloneUrls{},
			CreatedAt: handlerTestCreatedAtText,
			UpdatedAt: handlerTestUpdatedAtText,
		}, got)
		require.Empty(t, got.Labels)
	})
}

func TestDatasetTreeEntryToProtoFile(t *testing.T) {
	entries := handlerTestTreeEntries()
	files := handlerTestTreeFiles()

	requireProtoEqual(t, files[1], datasetTreeEntryToProtoFile(entries[1]))

	dir := datasetTreeEntryToProtoFile(entries[0])
	requireProtoEqual(t, files[0], dir)
	require.Nil(t, dir.Commit)
}

func TestDatasetHandlerListDatasetTaskLabels(t *testing.T) {
	ctx := context.Background()

	t.Run("returns task labels", func(t *testing.T) {
		h, svc := datasetHandlerTestSetup(t)
		svc.EXPECT().ListDatasetTaskLabels(ctx).Return([]*model.Label{
			{ID: 1, Name: "text-classification", Category: "task", Scope: "dataset", CreatedAt: handlerTestCreatedAt, UpdatedAt: handlerTestUpdatedAt},
			{ID: 5, Name: "translation", Category: "task", Scope: "dataset", CreatedAt: handlerTestCreatedAt, UpdatedAt: handlerTestUpdatedAt},
		}, nil)

		got, err := h.ListDatasetTaskLabels(ctx, &v1alpha1.ListDatasetTaskLabelsRequest{})
		require.NoError(t, err)
		requireProtoEqual(t, &v1alpha1.ListDatasetTaskLabelsResponse{Items: []*v1alpha1.Label{
			{Id: 1, Name: "text-classification", Category: v1alpha1.Category_TASK, CreatedAt: handlerTestCreatedAtText, UpdatedAt: handlerTestUpdatedAtText},
			{Id: 5, Name: "translation", Category: v1alpha1.Category_TASK, CreatedAt: handlerTestCreatedAtText, UpdatedAt: handlerTestUpdatedAtText},
