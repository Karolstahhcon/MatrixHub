import {
  Badge,
  Group,
  Text,
} from '@mantine/core'
import { UserSource } from '@matrixhub/api-ts/v1alpha1/user.pb'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import {
  DataTable,
  type DataTableProps,
  type DataTableRowActionsProps,
} from '@/shared/components/DataTable'
import { formatDateTime } from '@/shared/utils/date'

import { DeleteUserAction } from './DeleteUserAction'
import { ResetUserPasswordAction } from './ResetUserPasswordAction'
import { ToggleUserAdminAction } from './ToggleUserAdminAction'
import { getUserRowId } from '../users.utils'

import type { User } from '@matrixhub/api-ts/v1alpha1/user.pb'
import type { MRT_ColumnDef } from 'mantine-react-table'

type UserCellProps = Parameters<NonNullable<MRT_ColumnDef<User>['Cell']>>[0]

type UsersTableProps = Omit<DataTableProps<User>, 'columns'>

function UserAdminCell({ row }: UserCellProps) {
  const { t } = useTranslation()
  const isAdmin = !!row.original.isAdmin

  return (
    <Badge
      color={isAdmin ? 'green' : 'red'}
      variant="light"
    >
      {isAdmin
        ? t('routes.admin.users.boolean.yes')
        : t('routes.admin.users.boolean.no')}
    </Badge>
  )
}
