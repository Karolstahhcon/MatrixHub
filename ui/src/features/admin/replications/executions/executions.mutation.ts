import {
  type StopSyncTaskRequest,
  SyncPolicy,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { mutationOptions } from '@tanstack/react-query'

import i18n from '@/i18n'

import { adminReplicationExecutionKeys } from './executions.query'

import type { NotificationMeta } from '@/types/tanstack-query'
