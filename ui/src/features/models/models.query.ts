import {
  type GetModelBlobRequest,
  type GetModelTreeRequest,
  type ListModelCommitsRequest,
  type ListModelsRequest,
  Models,
} from '@matrixhub/api-ts/v1alpha1/model.pb'
import { ProjectPermissionFilter, Projects } from '@matrixhub/api-ts/v1alpha1/project.pb'
import { queryOptions, useQuery } from '@tanstack/react-query'

import { DEFAULT_PAGE_SIZE } from '@/utils/constants.ts'

import type { ModelsCatalogSearch } from '@/routes/(auth)/(app)/models'
import type { ProjectModelsSearch } from '@/routes/(auth)/(app)/projects/$projectId/models'

// -- Query key factory --

export const modelKeys = {
  all: ['models'] as const,

  lists: () => [...modelKeys.all, 'list'] as const,
  listByProject: (projectId: string, params: ListModelsRequest) => [...modelKeys.lists(), projectId, params] as const,
  listByCategory: (params: ListModelsRequest) => [...modelKeys.lists(), params] as const,

  details: () => [...modelKeys.all, 'detail'] as const,
  detail: (projectId: string, modelName: string) => [...modelKeys.details(), projectId, modelName] as const,

  taskLabels: () => [...modelKeys.all, 'task-labels'] as const,
  libraryLabels: () => [...modelKeys.all, 'library-labels'] as const,
  projects: (permissionFilter: ProjectPermissionFilter) => [...modelKeys.all, 'projects', permissionFilter] as const,

  commits: () => [...modelKeys.all, 'commit-list'] as const,
  commitsList: (projectId: string, modelName: string, params: Pick<ListModelCommitsRequest, 'revision' | 'page' | 'pageSize'>) => [
    ...modelKeys.commits(), projectId, modelName, params,
  ] as const,

  commitDetails: () => [...modelKeys.all, 'commit-detail'] as const,
  commitDetail: (projectId: string, modelName: string, commitId: string) => [
    ...modelKeys.commitDetails(),
    projectId,
    modelName,
    commitId,
  ] as const,
}

export const modelRevisionKeys = {
  all: ['modelRevisions'] as const,
  detail: (projectId: string, modelName: string) => [...modelRevisionKeys.all, projectId, modelName] as const,
}

export const modelTreeKeys = {
  all: ['modelTree'] as const,
  detail: (
    projectId: string,
    modelName: string,
    params?: Pick<GetModelTreeRequest, 'revision' | 'path'>,
  ) => [...modelTreeKeys.all, projectId, modelName, params] as const,
}

export const modelBlobKeys = {
  all: ['modelBlob'] as const,
  detail: (
    projectId: string,
    modelName: string,
    params?: Pick<GetModelBlobRequest, 'revision' | 'path'>,
  ) => [...modelBlobKeys.all, projectId, modelName, params] as const,
}

// -- Query options factory --

export function projectModelsQueryOptions(projectId: string, search: ProjectModelsSearch) {
  const sortParam = toSortParam(search.sort, search.order)

  return queryOptions({
    queryKey: modelKeys.listByProject(projectId, {
      search: search.q,
      sort: sortParam,
      page: search.page,
    }),
    queryFn: () => Models.ListModels({
      project: projectId,
      search: search.q || undefined,
      sort: sortParam,
      page: search.page,
      pageSize: DEFAULT_PAGE_SIZE,
    }),
  })
}

export function catalogModelsQueryOptions(
  search: ModelsCatalogSearch & {
    popular?: boolean
    pageSize?: number
  }) {
  const sortParam = toSortParam(search.sort, search.order)

  return queryOptions({
    queryKey: modelKeys.listByCategory({
      search: search.q,
