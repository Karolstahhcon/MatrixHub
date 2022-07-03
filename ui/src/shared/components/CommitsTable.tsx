import {
  Box,
  Group,
  Text,
} from '@mantine/core'
import { type LinkComponentProps } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

import AnchorLink from '@/shared/components/AnchorLink.tsx'
import { CopyValueButton } from '@/shared/components/CopyValueButton'
import { DataTable, type DataTableRowActionsProps } from '@/shared/components/DataTable'
import { formatDateTime } from '@/shared/utils/date'

import type { Commit } from '@matrixhub/api-ts/v1alpha1/model.pb'
import type { Pagination as PaginationData } from '@matrixhub/api-ts/v1alpha1/utils.pb'
import type { MRT_ColumnDef } from 'mantine-react-table'

interface CommitsTableProps {
  commits: Commit[]
  pagination?: PaginationData
  page: number
  loading?: boolean
  onPageChange: (page: number) => void
  getDetailLinkProps: (commit: Commit) => LinkComponentProps
}

type CommitCellProps = Parameters<NonNullable<MRT_ColumnDef<Commit>['Cell']>>[0]

function CommitMessageCell({ row }: CommitCellProps) {
  return (
    <Text size="sm" py="xs">
      {row.original.message || '-'}
    </Text>
  )
}

function CommitHashCell({ row }: CommitCellProps) {
  const commitId = row.original.id

  if (!commitId) {
    return <Text size="sm">-</Text>
  }

  return (
    <Group gap="xs" wrap="nowrap">
      <Text size="sm" truncate>
        {commitId}
      </Text>
      <CopyValueButton
        value={commitId}
        iconSize={12}
