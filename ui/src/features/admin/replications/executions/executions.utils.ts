import {
  SyncTaskStatus,
  type SyncTask,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import dayjs from 'dayjs'

import i18n from '@/i18n'
import { toMilliseconds } from '@/shared/utils/date'

export function getReplicationExecutionRowId(task: SyncTask) {
  return String(task.id ?? `${task.startedTimestamp ?? '0'}-${task.status ?? 'unknown'}`)
}

export function parseReplicationId(replicationId: string) {
