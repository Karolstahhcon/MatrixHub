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
  modelPath,
  defaultTask,
  onClose,
}: UseModelDrawerProps) {
  const { t } = useTranslation()
  const systemConfigQuery = useSystemConfig()
  const [task, setTask] = useState<UseModelTask | null>(defaultTask)

  const hfEndpoint = systemConfigQuery.data?.endpoints?.hfBase || window.location.origin
  const engineLabel = ENGINE_LABELS[engine]
  const prompts: SnippetPrompts = {
    textPrompt: t('model.detail.useModel.prompts.text'),
    imageUrl: t('model.detail.useModel.prompts.imageUrl'),
    imagePrompt: t('model.detail.useModel.prompts.image'),
  }

  const taskSelect = (
    <Group gap="xs" wrap="nowrap">
      <Text size="xs" c="dimmed">{t('model.detail.useModel.currentTask')}</Text>
      <Select
        size="xs"
        w={150}
        allowDeselect={false}
        value={task}
        placeholder={t('model.detail.useModel.selectTask')}
        onChange={value => value && setTask(value as UseModelTask)}
        data={TASK_OPTIONS.map(option => ({
          value: option,
          label: t(`model.detail.useModel.tasks.${option}`),
        }))}
      />
    </Group>
  )

  const taskRequired = (
    <Alert color="yellow" variant="light">
      {t('model.detail.useModel.unsupportedTask')}
    </Alert>
  )

  const envStep = (
    <Step
      index={3}
      title={t('model.detail.useModel.steps.configure')}
      hint={t('model.detail.useModel.steps.configureHint')}
    >
      {systemConfigQuery.isPending
        ? <Skeleton height={72} radius="sm" />
        : <SnippetBlock snippet={buildEnvSnippet(hfEndpoint)} />}
    </Step>
  )

  return (
    <GuideDrawer
      opened={opened}
      onClose={onClose}
      title={t('model.detail.useModel.title', { engine: engineLabel })}
      intro={t(`model.detail.useModel.intro.${engine}`)}
      footerExtra={(
