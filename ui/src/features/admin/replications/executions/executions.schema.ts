import { SyncTaskStatus } from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { z } from 'zod'

import { DEFAULT_PAGE } from '@/utils/constants'

export const replicationExecutionStatusFilterValues = [
  SyncTaskStatus.SYNC_TASK_STATUS_RUNNING,
  SyncTaskStatus.SYNC_TASK_STATUS_STOPPED,
  SyncTaskStatus.SYNC_TASK_STATUS_SUCCEEDED,
  SyncTaskStatus.SYNC_TASK_STATUS_FAILED,
] as const

export type ReplicationExecutionStatusFilterValue
  = typeof replicationExecutionStatusFilterValues[number]

