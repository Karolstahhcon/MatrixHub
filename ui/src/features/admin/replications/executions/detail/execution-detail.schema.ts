import {
  ResourceType,
  SyncJobStatus,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { z } from 'zod'

import { DEFAULT_PAGE } from '@/utils/constants'

export const syncJobStatusFilterValues = [
  SyncJobStatus.SYNC_JOB_STATUS_RUNNING,
  SyncJobStatus.SYNC_JOB_STATUS_SUCCEEDED,
  SyncJobStatus.SYNC_JOB_STATUS_FAILED,
  SyncJobStatus.SYNC_JOB_STATUS_STOPPED,
