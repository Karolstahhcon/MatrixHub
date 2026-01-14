import {
  Badge,
  Button,
  Group,
  Text,
} from '@mantine/core'
import {
  RobotAccountExpireStatus,
  RobotAccountProjectScope,
} from '@matrixhub/api-ts/v1alpha1/robot.pb'
import { useTranslation } from 'react-i18next'

import {
  DataTable,
  type DataTableProps,
  type DataTableRowActionsProps,
} from '@/shared/components/DataTable'
import { TruncatedText } from '@/shared/components/TruncatedText'
import { formatDateTime } from '@/shared/utils/date'

import {
  getRobotRowId,
  isRobotEnabled,
  type RobotAccount,
} from '../robot.utils'

import type { MRT_ColumnDef } from 'mantine-react-table'

type RobotCellProps = Parameters<NonNullable<MRT_ColumnDef<RobotAccount>['Cell']>>[0]

type RobotTableProps = Omit<DataTableProps<RobotAccount>, 'columns'> & {
  onEdit?: (robot: RobotAccount) => void
  onDelete?: (robot: RobotAccount) => void
  onToggleStatus?: (robot: RobotAccount) => void
  onRefreshToken?: (robot: RobotAccount) => void
}

function RobotStatusCell({ row }: RobotCellProps) {
  const { t } = useTranslation()
  const enabled = isRobotEnabled(row.original)

  return (
    <Badge
      color={enabled ? 'green' : 'gray'}
      variant="light"
    >
      {enabled
        ? t('routes.admin.robots.boolean.yes')
        : t('routes.admin.robots.boolean.no')}
    </Badge>
  )
}

function RobotNameCell({ row }: RobotCellProps) {
  return <TruncatedText value={row.original.name || '-'} />
}

function RobotPlatformPermissionsCell({ row }: RobotCellProps) {
