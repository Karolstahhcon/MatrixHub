import { SyncTaskStatus } from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { z } from 'zod'

import { DEFAULT_PAGE } from '@/utils/constants'

export const replicationExecutionStatusFilterValues = [
  SyncTaskStatus.SYNC_TASK_STATUS_RUNNING,
