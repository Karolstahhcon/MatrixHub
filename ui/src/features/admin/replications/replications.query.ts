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
} from './replications.schema'

export const adminReplicationKeys = {
  all: ['admin', 'replications'] as const,
  lists: () => [...adminReplicationKeys.all, 'list'] as const,
  list: (search: ReplicationsSearch) => [...adminReplicationKeys.lists(), search] as const,
  details: () => [...adminReplicationKeys.all, 'detail'] as const,
  detail: (syncPolicyId: number) => [...adminReplicationKeys.details(), syncPolicyId] as const,
}

export function replicationsQueryOptions(search: ReplicationsSearch) {
  return queryOptions({
    queryKey: adminReplicationKeys.list(search),
