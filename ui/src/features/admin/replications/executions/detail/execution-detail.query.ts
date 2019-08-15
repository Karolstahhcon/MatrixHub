import { SyncPolicy } from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import {
  keepPreviousData,
  queryOptions,
  useQuery,
} from '@tanstack/react-query'

import { DEFAULT_PAGE_SIZE } from '@/utils/constants'

import { type ExecutionDetailSearch } from './execution-detail.schema'

export const adminReplicationExecutionDetailKeys = {
  all: ['admin', 'replications', 'executions', 'detail'] as const,
  item: (syncPolicyId: number, syncTaskId: number) => [
    ...adminReplicationExecutionDetailKeys.all,
    syncPolicyId,
    syncTaskId,
  ] as const,
  jobs: (syncPolicyId: number, syncTaskId: number) => [
    ...adminReplicationExecutionDetailKeys.item(syncPolicyId, syncTaskId),
