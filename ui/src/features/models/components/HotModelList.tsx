import {
  Alert, Box, Button, Collapse, Group, Text, UnstyledButton,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { IconChevronDown } from '@tabler/icons-react'
import { useQuery } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

import { catalogModelsQueryOptions } from '@/features/models/models.query.ts'
import { ModelCard } from '@/shared/components/resource-card/ModelCard.tsx'
import { ResourceCardGrid } from '@/shared/components/ResourceCardGrid'

const { useSearch } = getRouteApi('/(auth)/(app)/models/')

export function HotModelList() {
  const { t } = useTranslation()
  const search = useSearch()

  const [opened, { toggle }] = useDisclosure(false)

  const {
    data: {
      items = [],
    } = {},
    isLoading,
    isPending,
    isLoadingError,
    refetch,
  } = useQuery(catalogModelsQueryOptions({
    ...search,
    page: 1,
    pageSize: -1,
