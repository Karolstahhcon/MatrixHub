import {
  Button,
  Group,
  Paper,
  Stack,
  Title,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import {
  IconApiApp as ProjectIcon,
  IconHomePlus,
} from '@tabler/icons-react'
import { useQueryClient } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { CreateProjectModal } from '../components/CreateProjectModal'
import { DeleteProjectModal } from '../components/DeleteProjectModal'
import { ProjectsTable } from '../components/ProjectsTable'
import { projectKeys, useProjects } from '../projects.query'

import type { Project } from '@matrixhub/api-ts/v1alpha1/project.pb'

const projectsRouteApi = getRouteApi('/(auth)/(app)/projects/')

export function ProjectsPage() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const navigate = projectsRouteApi.useNavigate()
  const search = projectsRouteApi.useSearch()

  const {
    data, isLoading, isFetching,
  } = useProjects({
    query: search.query ?? '',
    page: search.page ?? 1,
  })

  const projects = data?.projects ?? []
  const pagination = data?.pagination

  const [createOpened, createHandlers] = useDisclosure(false)
  const [deleteOpened, deleteHandlers] = useDisclosure(false)
  const [deleteTarget, setDeleteTarget] = useState<Project | null>(null)

  const handleSearchChange = (value: string) => {
    if (value === (search.query ?? '')) {
      return
    }

    void navigate({
      replace: true,
      search: prev => ({
        ...prev,
        page: 1,
        query: value,
      }),
    })
  }

  const handleDelete = (project: Project) => {
    setDeleteTarget(project)
    deleteHandlers.open()
  }

  const handleRefresh = () => {
    void queryClient.invalidateQueries({ queryKey: projectKeys.all })
  }

  const handlePageChange = (page: number) => {
    void navigate({
      search: prev => ({
        ...prev,
        page,
      }),
    })
  }

  return (
    <Stack gap="lg" pt="lg">
