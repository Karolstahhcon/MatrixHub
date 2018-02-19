import {
  Alert,
  Box,
  Center,
  Flex,
  Group,
  Loader,
  Stack,
} from '@mantine/core'
import { IconInfoCircle } from '@tabler/icons-react'
import { getRouteApi, linkOptions } from '@tanstack/react-router'

import { FileViewer, useFileContent } from '@/features/file-viewer'
import { ModelRevisionSelect } from '@/features/models/components/ModelRevisionSelect.tsx'
import { useModelBlob } from '@/features/models/models.query.ts'
import { PathBreadcrumbs } from '@/shared/components/PathBreadcrumbs.tsx'

const { useParams } = getRouteApi('/(auth)/(app)/projects_/$projectId/models/$modelId/blob/$ref/$')

export function ModelBlobPage() {
  const {
    projectId, modelId, ref: revision, _splat: path,
  } = useParams()

  const {
