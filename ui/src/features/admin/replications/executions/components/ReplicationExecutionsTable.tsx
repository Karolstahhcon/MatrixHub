import { Text } from '@mantine/core'
import {
  SyncTaskStatus,
  TriggerType,
  type SyncTask,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import {
  DataTable,
  type DataTableProps,
  type DataTableRowActionsProps,
} from '@/shared/components/DataTable'
import { StatusIndicator } from '@/shared/components/StatusIndicator'
import { TableCellLink } from '@/shared/components/TableCellLink'
import { formatDateTime } from '@/shared/utils/date'

import { StopReplicationExecutionAction } from './StopReplicationExecutionAction'
import { replicationExecutionStatusFilterValues } from '../executions.schema'
import {
  formatReplicationExecutionDuration,
  getReplicationExecutionRowId,
} from '../executions.utils'

import type { MRT_ColumnDef } from 'mantine-react-table'

type ExecutionCellProps = Parameters<NonNullable<MRT_ColumnDef<SyncTask>['Cell']>>[0]

function formatReplicationExecutionSuccessRate(task: SyncTask) {
  const totalItems = task.totalItems == null ? 0 : Number(task.totalItems)
  const successfulItems = task.successfulItems == null ? 0 : Number(task.successfulItems)

  if (!Number.isFinite(totalItems) || totalItems <= 0) {
    return '-'
  }

  const successRate = ((Number.isFinite(successfulItems) ? successfulItems : 0) / totalItems) * 100

  return `${Number(successRate.toFixed(2))}%`
}

function formatReplicationExecutionTotalItems(task: SyncTask) {
  const totalItems = task.totalItems == null ? 0 : Number(task.totalItems)

  if (!Number.isFinite(totalItems)) {
    return '-'
  }

