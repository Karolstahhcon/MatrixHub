import {
  Code, Skeleton, Stack, Text,
} from '@mantine/core'
import { useTranslation } from 'react-i18next'

import {
  GuideDrawer,
  GuideNote,
  GuideSnippetBlock,
  GuideStep,
} from '@/features/models/guides/GuideDrawer'
import { useSystemConfig } from '@/features/system/system.query'

import {
  buildDownloadSnippet,
  buildHfEnvSnippet,
  buildInstallHfSnippet,
  buildLocalPathPlaceholder,
  buildUploadSnippet,
  HF_DEFAULT_CACHE_DIR,
  type ModelCommandType,
} from './commandGuides'

interface ModelCommandDrawerProps {
  opened: boolean
  type: ModelCommandType
  modelPath: string
  onClose: () => void
}

/** Step-by-step guide for uploading to / downloading from a model repo with the `hf` CLI. */
export function ModelCommandDrawer({
  opened,
  type,
  modelPath,
  onClose,
}: ModelCommandDrawerProps) {
  const { t } = useTranslation()
  const systemConfigQuery = useSystemConfig()
  const hfEndpoint = systemConfigQuery.data?.endpoints?.hfBase || window.location.origin
  const localPath = buildLocalPathPlaceholder(modelPath)

  return (
    <GuideDrawer
      opened={opened}
      onClose={onClose}
      title={t(`model.detail.commandGuide.${type}.title`)}
      intro={t(`model.detail.commandGuide.${type}.intro`)}
    >
      <GuideStep
        index={1}
        title={t(type === 'download'
          ? 'model.detail.commandGuide.download.install'
          : 'model.detail.commandGuide.steps.install')}
      >
        <Text size="sm">{t('model.detail.commandGuide.steps.installHint')}</Text>
