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
		}}, got)
	})

	t.Run("hides service errors", func(t *testing.T) {
		h, svc := datasetHandlerTestSetup(t)
		svc.EXPECT().ListDatasetTaskLabels(ctx).Return(nil, errors.New("database is locked"))

		got, err := h.ListDatasetTaskLabels(ctx, &v1alpha1.ListDatasetTaskLabelsRequest{})
		require.Nil(t, got)
		requireHandlerTestStatus(t, err, codes.Internal, "failed to list task labels")
	})
}

func TestDatasetHandlerListDatasets(t *testing.T) {
	ctx := context.Background()
	request := &v1alpha1.ListDatasetsRequest{
		Label:    []string{"text-classification"},
		Search:   "im",
		Sort:     "updated_at",
		Project:  datasetHandlerTestProject,
		Page:     2,
		PageSize: 5,
		Popular:  true,
	}

	t.Run("forwards filters and paginates the result", func(t *testing.T) {
		h, svc := datasetHandlerTestSetup(t)
		svc.EXPECT().ListDatasets(ctx, gomock.Any()).DoAndReturn(
			func(_ context.Context, filter *model.Filter) ([]*dataset.Dataset, int64, error) {
				require.Equal(t, []string{"text-classification"}, filter.Label)
				require.Equal(t, "im", filter.Search)
				require.Equal(t, "updated_at", filter.Sort)
				require.Equal(t, datasetHandlerTestProject, filter.Project)
				require.Equal(t, int32(2), filter.Page)
				require.Equal(t, int32(5), filter.PageSize)
				require.NotNil(t, filter.Popular)
				require.True(t, *filter.Popular)
				return []*dataset.Dataset{datasetHandlerTestDataset()}, 12, nil
			})

		got, err := h.ListDatasets(ctx, request)
		require.NoError(t, err)
		requireProtoEqual(t, &v1alpha1.ListDatasetsResponse{
			Items:      []*v1alpha1.Dataset{datasetHandlerTestDatasetProto()},
			Pagination: &v1alpha1.Pagination{Total: 12, Page: 2, PageSize: 5, Pages: 3},
		}, got)
	})

	t.Run("passes a false popular flag through", func(t *testing.T) {
		h, svc := datasetHandlerTestSetup(t)
		svc.EXPECT().ListDatasets(ctx, gomock.Any()).DoAndReturn(
			func(_ context.Context, filter *model.Filter) ([]*dataset.Dataset, int64, error) {
				require.NotNil(t, filter.Popular)
				require.False(t, *filter.Popular)
				return nil, 0, nil
			})

		got, err := h.ListDatasets(ctx, &v1alpha1.ListDatasetsRequest{Page: 1, PageSize: 20})
		require.NoError(t, err)
		require.Empty(t, got.Items)
		requireProtoEqual(t, &v1alpha1.Pagination{Total: 0, Page: 1, PageSize: 20, Pages: 0}, got.Pagination)
	})

	t.Run("hides service errors", func(t *testing.T) {
		h, svc := datasetHandlerTestSetup(t)
		svc.EXPECT().ListDatasets(ctx, gomock.Any()).Return(nil, int64(0), errors.New("database is locked"))

		got, err := h.ListDatasets(ctx, request)
		require.Nil(t, got)
		requireHandlerTestStatus(t, err, codes.Internal, "failed to list datasets")
	})
}

func TestDatasetHandlerRejectsInvalidRequests(t *testing.T) {
	cases := []struct {
		name  string
		field string
		call  func(context.Context, *DatasetHandler) error
	}{
		{
			name:  "GetDataset without project",
			field: "GetDatasetRequest.Project",
			call: func(ctx context.Context, h *DatasetHandler) error {
				_, err := h.GetDataset(ctx, &v1alpha1.GetDatasetRequest{Name: datasetHandlerTestName})
				return err
			},
		},
		{
			name:  "GetDataset with reserved name",
			field: "GetDatasetRequest.Name",
			call: func(ctx context.Context, h *DatasetHandler) error {
				_, err := h.GetDataset(ctx, &v1alpha1.GetDatasetRequest{Project: datasetHandlerTestProject, Name: "datasets"})
				return err
			},
		},
		{
			name:  "CreateDataset without name",
			field: "CreateDatasetRequest.Name",
			call: func(ctx context.Context, h *DatasetHandler) error {
				_, err := h.CreateDataset(ctx, &v1alpha1.CreateDatasetRequest{Project: datasetHandlerTestProject})
				return err
			},
		},
		{
			name:  "CreateDataset with reserved name",
			field: "CreateDatasetRequest.Name",
			call: func(ctx context.Context, h *DatasetHandler) error {
				_, err := h.CreateDataset(ctx, &v1alpha1.CreateDatasetRequest{Project: datasetHandlerTestProject, Name: "models"})
				return err
			},
		},
		{
			name:  "DeleteDataset with reserved name",
			field: "DeleteDatasetRequest.Name",
			call: func(ctx context.Context, h *DatasetHandler) error {
				_, err := h.DeleteDataset(ctx, &v1alpha1.DeleteDatasetRequest{Project: datasetHandlerTestProject, Name: "resolve"})
				return err
			},
		},
		{
			name:  "ListDatasetRevisions without project",
			field: "ListDatasetRevisionsRequest.Project",
			call: func(ctx context.Context, h *DatasetHandler) error {
				_, err := h.ListDatasetRevisions(ctx, &v1alpha1.ListDatasetRevisionsRequest{Name: datasetHandlerTestName})
				return err
			},
		},
		{
			name:  "ListDatasetCommits without name",
			field: "ListDatasetCommitsRequest.Name",
			call: func(ctx context.Context, h *DatasetHandler) error {
				_, err := h.ListDatasetCommits(ctx, &v1alpha1.ListDatasetCommitsRequest{Project: datasetHandlerTestProject})
				return err
			},
		},
		{
			name:  "GetDatasetCommit without commit id",
			field: "GetDatasetCommitRequest.Id",
			call: func(ctx context.Context, h *DatasetHandler) error {
				_, err := h.GetDatasetCommit(ctx, &v1alpha1.GetDatasetCommitRequest{Project: datasetHandlerTestProject, Name: datasetHandlerTestName})
				return err
			},
		},
		{
			name:  "GetDatasetTree with reserved name",
			field: "GetDatasetTreeRequest.Name",
