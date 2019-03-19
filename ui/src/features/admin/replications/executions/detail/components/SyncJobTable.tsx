import {
  ActionIcon,
  Tooltip,
} from '@mantine/core'
import {
  ResourceType,
  SyncJobStatus,
  type SyncJob,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { IconFileText } from '@tabler/icons-react'
import { useTranslation } from 'react-i18next'

import {
  DataTable,
  type DataTableProps,
  type DataTableRowActionsProps,
} from '@/shared/components/DataTable'
import { StatusIndicator } from '@/shared/components/StatusIndicator'
import { formatDateTime } from '@/shared/utils/date'

import {
  syncJobStatusFilterValues,
  syncJobResourceTypeFilterValues,
} from '../execution-detail.schema'

import type { MRT_ColumnDef } from 'mantine-react-table'

type SyncJobTableProps = Omit<
  DataTableProps<SyncJob>,
  'columns' | 'renderToolbar'
> & {
  syncPolicyId: number
  syncTaskId: number
}

type JobCellProps = Parameters<NonNullable<MRT_ColumnDef<SyncJob>['Cell']>>[0]

function getSyncJobLogUrl(syncPolicyId: number, syncTaskId: number, syncJobId: number) {
  return `/api/v1alpha1/sync-policies/${syncPolicyId}/sync-tasks/${syncTaskId}/sync-jobs/${syncJobId}/log`
}

function JobStatusCell({ row }: JobCellProps) {
  const { t } = useTranslation()
  const status = row.original.status

  if (status === SyncJobStatus.SYNC_JOB_STATUS_RUNNING) {
    return (
      <StatusIndicator
        type="info"
        label={t('routes.admin.replications.executions.detail.status.running')}
      />
    )
  }

  if (status === SyncJobStatus.SYNC_JOB_STATUS_STOPPED) {
    return (
      <StatusIndicator
        type="warning"
        label={t('routes.admin.replications.executions.detail.status.stopped')}
      />
    )
  }

  if (status === SyncJobStatus.SYNC_JOB_STATUS_SUCCEEDED) {
    return (
      <StatusIndicator
        type="success"
        label={t('routes.admin.replications.executions.detail.status.succeeded')}
      />
    )
  }

  if (status === SyncJobStatus.SYNC_JOB_STATUS_FAILED) {
    return (
