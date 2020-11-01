import {
  Anchor,
  Box,
  Group,
  Text,
} from '@mantine/core'
import {
  SSHKeyStatus,
  type SSHKey,
} from '@matrixhub/api-ts/v1alpha1/current_user.pb'
import { useTranslation } from 'react-i18next'

import { formatExpiredAt } from '@/features/profile/profile.utils.ts'
import {
  DataTable,
  type DataTableProps,
  type DataTableRowActionsProps,
} from '@/shared/components/DataTable'
import { formatDateTime } from '@/shared/utils/date'

import type {
  MRT_ColumnDef,
  MRT_Row,
} from 'mantine-react-table'

const StatusCell = ({ row }: { row: MRT_Row<SSHKey> }) => {
  const { t } = useTranslation()
  const status = row.original.status ?? SSHKeyStatus.SSH_KEY_STATUS_UNKNOWN

  const statusColor = {
