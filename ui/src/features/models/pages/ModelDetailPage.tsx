import {
  Box,
  Button,
  Menu,
  Tabs,
} from '@mantine/core'
import { ProjectRoleType } from '@matrixhub/api-ts/v1alpha1/role.pb.ts'
import {
  IconChevronDown,
  IconCloudUpload,
  IconDownload,
  IconTerminal2,
} from '@tabler/icons-react'
import {
  Link,
  getRouteApi,
  linkOptions,
  useMatchRoute,
} from '@tanstack/react-router'
import { type ReactNode } from 'react'
import { useTranslation } from 'react-i18next'

import { useProjectRole } from '@/features/auth/useProjectRole'
import { type ModelCommandType } from '@/features/models/command-guide/commandGuides'
import { ModelCommandDrawer } from '@/features/models/command-guide/ModelCommandDrawer'
import { buildModelBadges, buildModelMetaItems } from '@/features/models/models.utils'
import { UseModelDrawer } from '@/features/models/use-model/UseModelDrawer'
import {
  ENGINE_LABELS,
  resolveUseModelTask,
  USE_MODEL_ENGINES,
  type UseModelEngine,
} from '@/features/models/use-model/useModelGuides'
import { ResourceDetailHeader } from '@/shared/components/ResourceDetailHeader'
import { usePayloadModal } from '@/shared/hooks/usePayloadModal'

interface ModelDetailPageProps {
  children: ReactNode
}

const {
  useLoaderData,
  useParams,
} = getRouteApi('/(auth)/(app)/projects_/$projectId/models/$modelId')

export function ModelDetailPage({
  children,
}: ModelDetailPageProps) {
  const {
