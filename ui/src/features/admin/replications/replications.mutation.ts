import {
  type CreateSyncPolicyRequest,
  type CreateSyncTaskRequest,
  type CreateSyncTaskResponse,
  type DeleteSyncPolicyRequest,
  SyncPolicy,
  type UpdateSyncPolicyRequest,
  type UpdateSyncPolicySwitchRequest,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { mutationOptions } from '@tanstack/react-query'

import i18n from '@/i18n'

import { adminReplicationKeys } from './replications.query'

import type { NotificationMeta } from '@/types/tanstack-query'

function requireSyncPolicyId(id?: number) {
  if (id == null) {
    throw new Error(i18n.t('routes.admin.replications.errors.missingSyncPolicyId'))
  }

  return id
}

export function createReplicationMutationOptions() {
  return mutationOptions({
    mutationFn: (input: CreateSyncPolicyRequest) => SyncPolicy.CreateSyncPolicy(input),
    meta: {
      successMessage: i18n.t('routes.admin.replications.notifications.createSuccess'),
