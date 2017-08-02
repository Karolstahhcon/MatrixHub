import {
  Anchor,
  Group,
  SegmentedControl,
  Text,
} from '@mantine/core'
import { IconClock } from '@tabler/icons-react'
import { useTranslation } from 'react-i18next'

import { formatDateTime } from '@/shared/utils/date'
import { formatStorageSize } from '@/shared/utils/format'

import type { FileViewerFile } from '../types'

type ViewMode = 'preview' | 'code'

interface FileToolbarProps {
  file: FileViewerFile
  /** Whether this file is markdown — controls Preview/Code toggle visibility */
  isMarkdown?: boolean
  /** Current view mode (only relevant when isMarkdown is true) */
  viewMode?: ViewMode
  /** Callback when the user switches between Preview and Code */
  onViewModeChange?: (mode: ViewMode) => void
}

export function FileToolbar({
  file,
  isMarkdown = false,
  viewMode = 'code',
  onViewModeChange,
}: FileToolbarProps) {
  const { t } = useTranslation()

  return (
    <Group
      justify="space-between"
      px="sm"
      py="xs"
      bg="var(--mantine-color-default-hover)"
    >
      {/* Left side: commit info */}
      <Group
        gap="sm"
        style={{
          flex: 1,
          minWidth: 0,
          flexWrap: 'nowrap',
        }}
      >
        {file.commit?.authorName
          ? (
              <Text size="sm" c="dimmed" truncate>
                {file.commit.authorName}
                {' '}
                {file.commit.message}
              </Text>
            )
          : null}

        {file.commit?.id
          ? (
              <Text
                size="xs"
                ff="monospace"
                px={6}
                py={2}
