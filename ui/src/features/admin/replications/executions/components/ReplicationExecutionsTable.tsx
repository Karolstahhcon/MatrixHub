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

  return new Intl.NumberFormat().format(totalItems)
}

type ReplicationExecutionsTableProps = Omit<
  DataTableProps<SyncTask>,
  'columns' | 'renderToolbar'
> & {
  syncPolicyId: number
}

function ExecutionTaskCell({
  row, table,
}: ExecutionCellProps) {
  const syncPolicyId = (
    table.options.meta as { syncPolicyId?: number } | undefined
  )?.syncPolicyId

  const executionId = row.original.id

  if (executionId == null || syncPolicyId == null) {
    return <Text size="sm" fw={600}>{executionId ?? '-'}</Text>
  }

  return (
    <TableCellLink
      to="/admin/replications/$replicationId/executions/$executionId"
      params={{
        replicationId: String(syncPolicyId),
        executionId: String(executionId),
      }}
    >
      {executionId}
    </TableCellLink>
  )
}

function ExecutionStatusCell({ row }: ExecutionCellProps) {
  const { t } = useTranslation()
  const status = row.original.status

  if (status === SyncTaskStatus.SYNC_TASK_STATUS_RUNNING) {
    return (
      <StatusIndicator
        type="info"
        label={t('routes.admin.replications.executions.status.running')}
      />
    )
  }

  if (status === SyncTaskStatus.SYNC_TASK_STATUS_STOPPED) {
    return (
      <StatusIndicator
        type="warning"
        label={t('routes.admin.replications.executions.status.stopped')}
      />
    )
  }

  if (status === SyncTaskStatus.SYNC_TASK_STATUS_SUCCEEDED) {
    return (
      <StatusIndicator
        type="success"
        label={t('routes.admin.replications.executions.status.succeeded')}
      />
    )
  }

  if (status === SyncTaskStatus.SYNC_TASK_STATUS_FAILED) {
    return (
      <StatusIndicator
        type="danger"
        label={t('routes.admin.replications.executions.status.failed')}
      />
    )
  }

  return (
    <StatusIndicator
      label={t('routes.admin.replications.executions.status.unknown')}
    />
  )
}

function ExecutionActionsCell({
  row,
  table,
}: DataTableRowActionsProps<SyncTask>) {
  const syncPolicyId = (
    table.options.meta as { syncPolicyId?: number } | undefined
  )?.syncPolicyId

  if (syncPolicyId == null) {
    return <Text size="sm">-</Text>
  }

  return (
    <StopReplicationExecutionAction
      syncPolicyId={syncPolicyId}
