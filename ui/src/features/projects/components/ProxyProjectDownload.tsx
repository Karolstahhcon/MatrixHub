import {
  Alert,
  Anchor,
  Box,
  Button,
  Drawer,
  Group,
  rem,
  SimpleGrid,
  Skeleton,
  Stack,
  Text,
  TextInput,
  ThemeIcon,
  Title,
} from '@mantine/core'
import {
  IconExternalLink,
  IconInfoCircle,
} from '@tabler/icons-react'
import { useTranslation } from 'react-i18next'

import IconProxyProjectEmpty from '@/assets/svgs/proxy-project-empty.svg?react'
import { createModelSchema } from '@/features/models/models.schema'
import { useSystemConfig } from '@/features/system/system.query'
import { ShikiCodeBlock } from '@/shared/components/ShikiCodeBlock'
import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form'

import { buildProxyDownloadCommand } from '../projects.utils'

import type { ReactNode } from 'react'

const MODEL_DOWNLOAD_DOC_URL = '/docs/operations/model-repo/upload-download/'
const HF_ENDPOINT_PLACEHOLDER = window.location.origin

interface ProxyProjectDownloadProps {
  organization?: string
  remoteOrganization?: string
  requiresToken?: boolean
}

interface ProxyProjectDownloadDrawerProps extends ProxyProjectDownloadProps {
  opened: boolean
  onClose: () => void
}

function StepTitle({
  number,
  children,
}: {
  number: number
  children: ReactNode
}) {
  return (
    <Group gap="xs">
      <ThemeIcon size="sm" radius="xl" variant="filled">
        <Text size="xs" lh="sm" fw={700}>{number}</Text>
      </ThemeIcon>
      <Text size="sm" fw={400}>{children}</Text>
    </Group>
  )
}

function ProxyProjectDownloadFlow({
  organization,
  remoteOrganization,
  requiresToken,
  layout,
}: ProxyProjectDownloadProps & {
  layout: 'drawer' | 'guide'
}) {
  const { t } = useTranslation()
  const systemConfigQuery = useSystemConfig()
  const form = useForm({
    defaultValues: {
      modelName: '',
    },
  })
  const { name: modelNameSchema } = createModelSchema(t).shape
  const hfEndpoint = systemConfigQuery.data?.endpoints?.hfBase || HF_ENDPOINT_PLACEHOLDER

  const steps = [
    (
      <Stack key="model-name" gap="sm">
        <StepTitle number={1}>{t('projects.detail.proxyDownload.steps.modelName')}</StepTitle>
        <form.Field
          name="modelName"
          validators={{ onChange: modelNameSchema }}
        >
          {field => (
            <TextInput
              required
              label={t('projects.detail.proxyDownload.modelNameLabel')}
              leftSection={remoteOrganization ? `${remoteOrganization} /` : undefined}
              leftSectionPointerEvents="none"
