import { Projects } from '@matrixhub/api-ts/v1alpha1/project.pb'
import {
  keepPreviousData,
  queryOptions,
  useQuery,
} from '@tanstack/react-query'
import { z } from 'zod'

export const MEMBERS_PAGE_SIZE = 10

export const membersSearchDefaults = {
  q: '',
  page: 1,
}
