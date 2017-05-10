import {
  Box,
  Flex,
  Group,
  rem,
  Text,
} from '@mantine/core'
import { IconCopy } from '@tabler/icons-react'
import { useTranslation } from 'react-i18next'

import { CopyValueButton } from '@/shared/components/CopyValueButton'

interface CommitHashCardProps {
  commitId?: string
}

export function CommitHashCard({ commitId }: CommitHashCardProps) {
  const { t } = useTranslation()

  const shownCommitId = commitId ?? ''
  const shortCommitId = shownCommitId.slice(0, 7)

  if (!shortCommitId) {
    return null
  }

  return (
    <Group gap={8} wrap="nowrap">
      <Text fw={600} size="sm">
        {t('shared.commitDetail.commit')}
      </Text>
      <Box
