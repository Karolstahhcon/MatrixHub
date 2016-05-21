import {
  type StopSyncTaskRequest,
  SyncPolicy,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { mutationOptions } from '@tanstack/react-query'

import i18n from '@/i18n'

import { adminReplicationExecutionKeys } from './executions.query'

import type { NotificationMeta } from '@/types/tanstack-query'

function requireSyncTaskId(id?: number) {
  if (id == null) {
    throw new Error(i18n.t('routes.admin.replications.executions.errors.missingSyncTaskId'))
  }

  return id
}

export function stopReplicationExecutionMutationOptions(syncPolicyId: number) {
  return mutationOptions({
