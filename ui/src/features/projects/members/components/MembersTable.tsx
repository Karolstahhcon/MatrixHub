import {
  Anchor,
  Avatar,
  Group,
  Text,
} from '@mantine/core'
import { MemberType } from '@matrixhub/api-ts/v1alpha1/project.pb'
import { ProjectRoleType } from '@matrixhub/api-ts/v1alpha1/role.pb'
import { useTranslation } from 'react-i18next'

import {
  DataTable,
  type DataTableProps,
  type DataTableRowActionsProps,
} from '@/shared/components/DataTable'

import type { ProjectMember } from '@matrixhub/api-ts/v1alpha1/project.pb'
import type { MRT_ColumnDef } from 'mantine-react-table'
import type { ReactNode } from 'react'

export interface MembersTableProps extends Omit<DataTableProps<ProjectMember>, 'columns'> {
  currentRole?: ProjectRoleType
  currentUsername?: string
  onEditRole: (member: ProjectMember) => void
  onRemove: (member: ProjectMember) => void
  toolbarExtra?: ReactNode
}

type MemberCellProps = Parameters<NonNullable<MRT_ColumnDef<ProjectMember>['Cell']>>[0]

function MemberNameCell({ row }: MemberCellProps) {
  const name = row.original.memberName

  return (
    <Group gap="sm" wrap="nowrap">
      <Avatar size={24} radius="xl" color="gray">
        {name?.charAt(0)?.toUpperCase()}
      </Avatar>
      <Text fw={600} size="sm" truncate>
        {name || '-'}
      </Text>
    </Group>
  )
}

function useRoleLabel() {
  const { t } = useTranslation()

  return (role?: ProjectRoleType) => {
    switch (role) {
      case ProjectRoleType.ROLE_TYPE_PROJECT_ADMIN:
        return t('projects.detail.membersPage.role.admin')
