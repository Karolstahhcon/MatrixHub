import {
  Anchor,
  Badge,
  Group,
  Text,
} from '@mantine/core'
import { RegistryStatus } from '@matrixhub/api-ts/v1alpha1/registry.pb'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import {
  DataTable,
  type DataTableProps,
  type DataTableRowActionsProps,
} from '@/shared/components/DataTable'
import { TruncatedText } from '@/shared/components/TruncatedText'
import { formatDateTime } from '@/shared/utils/date'

import { DeleteRegistryAction } from './DeleteRegistryAction'
import { EditRegistryAction } from './EditRegistryAction'
import { registryProviderLabelKeys } from '../registries.form'
import { getRegistryRowId } from '../registries.utils'

import type { Registry } from '@matrixhub/api-ts/v1alpha1/registry.pb'
import type { MRT_ColumnDef } from 'mantine-react-table'

type RegistryCellProps = Parameters<NonNullable<MRT_ColumnDef<Registry>['Cell']>>[0]

type RegistriesTableProps = Omit<DataTableProps<Registry>, 'columns'>

function RegistryNameCell({ row }: RegistryCellProps) {
  const name = row.original.name

  return <TruncatedText value={name || '-'} />
}

function RegistryTypeCell({ row }: RegistryCellProps) {
  const { t } = useTranslation()
  const type = row.original.type
  const labelKey = type == null ? undefined : registryProviderLabelKeys[type]
  const label = labelKey ? t(labelKey) : '-'

  return <Text size="sm">{label}</Text>
}

function RegistryStatusCell({ row }: RegistryCellProps) {
  const { t } = useTranslation()
  const status = row.original.status

  let color = 'gray'
  let label = '-'

  if (status === RegistryStatus.REGISTRY_STATUS_HEALTHY) {
    color = 'green'
    label = t('routes.admin.registries.status.healthy')
  } else if (status === RegistryStatus.REGISTRY_STATUS_UNHEALTHY) {
    color = 'red'
    label = t('routes.admin.registries.status.unhealthy')
  }

  return (
    <Badge color={color} variant="light">
      {label}
    </Badge>
  )
}

function RegistryUrlCell({ row }: RegistryCellProps) {
  const url = row.original.url

  if (!url) {
    return <Text size="sm">-</Text>
  }

  return (
    <TruncatedText value={url}>
      <Anchor size="sm" href={url} target="_blank" rel="noopener noreferrer">
        {url}
      </Anchor>
    </TruncatedText>
  )
}

function RegistryInsecureCell({ row }: RegistryCellProps) {
  const { t } = useTranslation()
  const insecure = row.original.insecure

  if (insecure == null) {
    return <Text size="sm">-</Text>
  }

  const verified = !insecure

  return (
    <Badge color={verified ? 'green' : 'red'} variant="light">
      {verified
        ? t('routes.admin.registries.insecure.verified')
        : t('routes.admin.registries.insecure.unverified')}
    </Badge>
  )
