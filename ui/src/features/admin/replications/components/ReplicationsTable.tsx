import {
  Badge,
  Group,
  Stack,
  Text,
} from '@mantine/core'
import {
  SyncPolicyType,
  TriggerType,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import {
  DataTable,
  type DataTableProps,
  type DataTableRowActionsProps,
} from '@/shared/components/DataTable'
import { FieldHintLabel } from '@/shared/components/FieldHintLabel'
import { TableCellLink } from '@/shared/components/TableCellLink'
import { TruncatedText } from '@/shared/components/TruncatedText'

import { DeleteReplicationAction } from './DeleteReplicationAction'
import { EditReplicationAction } from './EditReplicationAction'
import { SyncReplicationAction } from './SyncReplicationAction'
import { ToggleReplicationAction } from './ToggleReplicationAction'
import {
  formatReplicationBandwidth,
  getCronExpressionDescription,
  getReplicationRowId,
} from '../replications.utils'

import type { Registry } from '@matrixhub/api-ts/v1alpha1/registry.pb'
import type { SyncPolicyItem } from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import type { MRT_ColumnDef } from 'mantine-react-table'

type ReplicationCellProps = Parameters<NonNullable<MRT_ColumnDef<SyncPolicyItem>['Cell']>>[0]

type ReplicationsTableProps = Omit<DataTableProps<SyncPolicyItem>, 'columns'>

const EMPTY_VALUE = '-'

function ReplicationTruncatedTextCell({ cell }: ReplicationCellProps) {
  return <TruncatedText value={cell.getValue<string>() || EMPTY_VALUE} />
}

function ReplicationNameCell({ row }: ReplicationCellProps) {
  const name = row.original.name
  const replicationId = row.original.id

  if (replicationId == null) {
    return (
      <TruncatedText value={name ?? EMPTY_VALUE} fw={500} />
    )
  }

  return (
    <TruncatedText value={name ?? EMPTY_VALUE}>
      <TableCellLink
        fw={500}
        to="/admin/replications/$replicationId/executions"
        params={{ replicationId: String(replicationId) }}
      >
        {name ?? EMPTY_VALUE}
      </TableCellLink>
    </TruncatedText>
  )
}

function ReplicationStatusCell({ row }: ReplicationCellProps) {
  const { t } = useTranslation()
  const isDisabled = row.original.isDisabled

  return (
    <Badge color={isDisabled ? 'red' : 'green'} variant="light">
      {isDisabled
        ? t('routes.admin.replications.status.disabled')
        : t('routes.admin.replications.status.enabled')}
    </Badge>
  )
}

function ReplicationTriggerTypeCell({ row }: ReplicationCellProps) {
  const {
    t,
    i18n,
  } = useTranslation()
  const triggerType = row.original.triggerType
