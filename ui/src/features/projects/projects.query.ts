import { Projects, ProjectPermissionFilter } from '@matrixhub/api-ts/v1alpha1/project.pb'
import {
  keepPreviousData,
  queryOptions,
  useQuery,
} from '@tanstack/react-query'

import { DEFAULT_PROJECTS_PAGE_SIZE } from '@/utils/constants'

export interface ProjectsSearch {
  query: string
