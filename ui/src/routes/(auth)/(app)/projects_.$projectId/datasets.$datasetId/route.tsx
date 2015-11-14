import {
  Box,
  Button,
  Space,
  Tabs,
} from '@mantine/core'
import { type Dataset } from '@matrixhub/api-ts/v1alpha1/dataset.pb.ts'
import { Category } from '@matrixhub/api-ts/v1alpha1/model.pb'
import {
  IconCloudUpload as UploadIcon,
  IconDownload as DownloadIcon,
} from '@tabler/icons-react'
import {
  Outlet,
  Link,
  createFileRoute,
  linkOptions,
  useMatchRoute,
} from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

import { useProjectRole } from '@/features/auth/useProjectRole'
import { ResourceDetailHeader } from '@/shared/components/ResourceDetailHeader'

import { Route as DatasetSettingsRoute } from './settings'
import { Route as DatasetTreeRoute } from './tree/$ref/$'

// TODO: Replace with real API data
const MOCK_DATA: Dataset = {
  size: '595 GB',
  updatedAt: '2021-12-17 12:12',
  labels: [
    {
      id: 1,
      name: 'Text',
      category: Category.TASK,
      createdAt: '2024-01-01T12:00:00Z',
      updatedAt: '2024-01-01T12:00:00Z',
    },
  ],
}

export const Route = createFileRoute(
  '/(auth)/(app)/projects_/$projectId/datasets/$datasetId',
)({
  component: DatasetLayout,
  loader: async () => {
    // TODO fetch dataset detail data, if the project is private and user has no access, will throw error
    // throw new NotFoundRouteError()

    return {
      dataset: MOCK_DATA,
    }
  },
})

function DatasetLayout() {
  const { t } = useTranslation()

  const {
    projectId, datasetId,
  } = Route.useParams()

  const projectRole = useProjectRole(projectId)
  const hasProjectRole = projectRole !== undefined

  const tabRoutes = linkOptions([
    {
      id: 'desc',
      label: t('dataset.detail.desc'),
      to: Route.to,
      params: {
        projectId,
        datasetId,
      },
    },
    {
      id: 'tree',
      label: t('dataset.detail.tree'),
      to: DatasetTreeRoute.to,
      params: {
        projectId,
        datasetId,
        ref: 'testDsd',
