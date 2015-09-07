import {
  Alert,
  Anchor,
  Box,
  Button,
  Center,
  Group,
  Paper,
  Space,
  Stack,
  Text,
} from '@mantine/core'
import { IconClock, IconCube } from '@tabler/icons-react'
import { useQuery } from '@tanstack/react-query'
import { getRouteApi, Link } from '@tanstack/react-router'
import { startTransition } from 'react'
import { useTranslation } from 'react-i18next'

import { catalogModelsQueryOptions } from '@/features/models/models.query'
import { Pagination } from '@/shared/components/Pagination'
import { ModelCard } from '@/shared/components/resource-card/ModelCard.tsx'
import { ResourceCardGrid } from '@/shared/components/ResourceCardGrid'
import { SearchToolbar } from '@/shared/components/SearchToolbar'
import {
  SortDropdown,
  type SortDropdownOption,
} from '@/shared/components/SortDropdown'
import { DEFAULT_PAGE_SIZE } from '@/utils/constants.ts'

const {
  useNavigate, useSearch,
} = getRouteApi('/(auth)/(app)/models/')

export function AllModelList() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const search = useSearch()

  const {
    data: {
      items = [],
      pagination,
    } = {},
    isLoading,
    isPending,
    isLoadingError,
    refetch,
  } = useQuery(catalogModelsQueryOptions(search))

  const hasActiveFilters = Boolean(search.q || search.task || search.library || search.project)
  const modelCount = pagination?.total ?? items.length
  const showCreatePrompt = !isLoading && !isPending && !isLoadingError
    && !hasActiveFilters && modelCount === 0

  const sortFieldOptions: SortDropdownOption[] = [
    {
      value: 'updatedAt',
      label: t('projects.detail.modelsPage.sortFieldUpdatedAt'),
      icon: <IconClock size={16} />,
    },
  ]

  const cardElements = items.map((model) => {
    const projectId = model.project?.trim() ?? '-'
    const modelName = model.name?.trim() ?? '-'

    return <ModelCard key={`${projectId}/${modelName}`} model={model} />
  })

  return (
    <Box>
      <Group>
        <Text fz="md" fw={600} lh="20px" mb="sm">
          {t('model.list.allModels') }
        </Text>
      </Group>

      <Stack gap={0}>
        <SearchToolbar
          searchPlaceholder={t('model.list.placeholder.modelName')}
          searchValue={search.q}
          onSearchChange={(nextQuery) => {
            void navigate({
              replace: true,
              search: prev => ({
                ...prev,
                q: nextQuery,
                page: 1,
              }),
            })
          }}
        >
          <SortDropdown
            fieldOptions={sortFieldOptions}
            fieldValue={search.sort}
            order={search.order}
            onFieldChange={(nextField) => {
              if (sortFieldOptions.find(o => o.value === nextField)?.disabled) {
                return
              }
