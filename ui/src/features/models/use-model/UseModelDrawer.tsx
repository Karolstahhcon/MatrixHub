import {
  Accordion,
  Alert,
  Anchor,
  Badge,
  Group,
  rem,
  Select,
  Skeleton,
  Stack,
  Tabs,
  Text,
  Tooltip,
} from '@mantine/core'
import {
  IconExternalLink,
  IconInfoCircle,
} from '@tabler/icons-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import {
  GuideDrawer,
  GuideSnippetBlock as SnippetBlock,
  GuideStep as Step,
} from '@/features/models/guides/GuideDrawer'
import guideClasses from '@/features/models/guides/GuideDrawer.module.css'
import { useSystemConfig } from '@/features/system/system.query'

import classes from './UseModelDrawer.module.css'
import {
  buildDockerSnippet,
  buildEnvSnippet,
  buildInstallSnippet,
  buildLowLevelSnippet,
  buildPipelineSnippet,
  buildServeSnippet,
  buildTestRequestSnippet,
  ENGINE_DOC_URLS,
  ENGINE_LABELS,
  type SnippetPrompts,
  type UseModelEngine,
  type UseModelTask,
} from './useModelGuides'

const TASK_OPTIONS: UseModelTask[] = ['text-generation', 'image-text-to-text']

interface UseModelDrawerProps {
  opened: boolean
  engine: UseModelEngine
  modelPath: string
  /** Task detected from the model's labels; the user may override it inside the drawer. */
  defaultTask: UseModelTask | null
  onClose: () => void
}

export function UseModelDrawer({
  opened,
  engine,
