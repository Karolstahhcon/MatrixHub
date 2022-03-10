import {
  Box, Button, Flex, Group, TextInput,
} from '@mantine/core'
import { IconClockHour4, IconSearch } from '@tabler/icons-react'
import { useQuery } from '@tanstack/react-query'
import {
  getRouteApi, Link, linkOptions,
} from '@tanstack/react-router'
import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { ModelRevisionSelect } from '@/features/models/components/ModelRevisionSelect.tsx'
import { modelCommitsQueryOptions, useModelTree } from '@/features/models/models.query.ts'
import { PathBreadcrumbs } from '@/shared/components/PathBreadcrumbs.tsx'
import { PathNotFound } from '@/shared/components/PathNotFound'
import { RepoFileTree } from '@/shared/components/repo-file-tree'

import type { File as RepoFile } from '@matrixhub/api-ts/v1alpha1/model.pb'

const { useParams } = getRouteApi('/(auth)/(app)/projects_/$projectId/models/$modelId/tree/$ref/$')

export function ModelTreePage() {
  const { t } = useTranslation()

  const {
    projectId, modelId, ref, _splat: treePath,
  } = useParams()

  const {
    data: { items } = {}, isLoading, isError,
  } = useModelTree(projectId, modelId, {
    revision: ref,
    path: treePath,
  })
  const { data } = useQuery(modelCommitsQueryOptions(projectId, modelId, {
    revision: ref,
    page: 1,
    pageSize: 1,
  }))

  const files = useMemo(() => items ?? [], [items])

  const [searchValue, setSearchValue] = useState('')

  const filteredFiles = useMemo(() => {
    if (!searchValue.trim()) {
      return files
    }
    const q = searchValue.toLocaleLowerCase()

    return files.filter(f => f.name?.toLocaleLowerCase().includes(q))
  }, [files, searchValue])

  const latestCommit = files.length > 0 ? findLatestCommit(files) : undefined

  const pathNotFound = !isLoading && (isError || (treePath && files.length === 0))

  return (
    <Box pt="sm" pb="xl">
      <Flex
        justify="space-between"
        align="center"
