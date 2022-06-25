import { SyncPolicy } from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import {
  keepPreviousData,
  queryOptions,
  useQuery,
} from '@tanstack/react-query'

import i18n from '@/i18n'

import {
  DEFAULT_REPLICATIONS_PAGE_SIZE,
  type ReplicationsSearch,
