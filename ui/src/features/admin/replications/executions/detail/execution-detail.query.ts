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
    'jobs',
  ] as const,
  jobList: (syncPolicyId: number, syncTaskId: number, search: ExecutionDetailSearch) => [
    ...adminReplicationExecutionDetailKeys.jobs(syncPolicyId, syncTaskId),
    search,
  ] as const,
}

export function executionDetailQueryOptions(
  syncPolicyId: number,
  syncTaskId: number,
) {
  return queryOptions({
    queryKey: adminReplicationExecutionDetailKeys.item(syncPolicyId, syncTaskId),
    queryFn: async () => {
      const response = await SyncPolicy.GetSyncTask({
        syncPolicyId,
        syncTaskId,
      })

