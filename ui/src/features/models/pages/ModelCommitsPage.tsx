import {
  Flex, Group,
  Stack,
} from '@mantine/core'
import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { getRouteApi, linkOptions } from '@tanstack/react-router'

import { ModelRevisionSelect } from '@/features/models/components/ModelRevisionSelect.tsx'
import { modelCommitsQueryOptions } from '@/features/models/models.query'
import { CommitsTable } from '@/shared/components/CommitsTable'
import { PathBreadcrumbs } from '@/shared/components/PathBreadcrumbs.tsx'

const {
  useNavigate, useParams, useSearch,
} = getRouteApi('/(auth)/(app)/projects_/$projectId/models/$modelId/commits/$ref/')

export function ModelCommitsPage() {
  const navigate = useNavigate()
