import {
  Group,
  Paper,
  Stack,
  Text,
  UnstyledButton,
} from '@mantine/core'
import { type Project } from '@matrixhub/api-ts/v1alpha1/project.pb.ts'
import { useTranslation } from 'react-i18next'

import { ProjectTypeBadge } from '@/shared/components/badges/ProjectTypeBadge'
import { BaseFilterPanel } from '@/shared/components/resource-filter-panel/BaseFilterPanel'
import { filterByKeyword } from '@/shared/utils'

interface ProjectCardProps {
  option: Project
  selected?: boolean
  onSelect?: (option: string) => void
}

function ProjectCard({
  option,
  selected = false,
  onSelect,
}: ProjectCardProps) {
  const { t } = useTranslation()

  return (
    <UnstyledButton
      w="100%"
      onClick={() => onSelect?.(option.name as string)}
    >
      <Paper
        radius="md"
        px="sm"
        pt="sm"
        pb={8}
        withBorder
        style={{
          borderColor: selected
            ? 'var(--mantine-primary-color-6)'
            : 'var(--mantine-color-gray-3)',
          transition: 'all 120ms ease',
        }}
      >
        <Group
          justify="space-between"
          wrap="nowrap"
          miw={0}
        >
          <Text
            fw={600}
            c="var(--mantine-color-gray-9)"
            size="sm"
            lh={1}
            truncate
          >
            {option.name}
          </Text>
          <ProjectTypeBadge type={option.type} flex="0 0 auto" />
        </Group>

        <Group
          justify="flex-end"
          gap={24}
          mt="md"
          wrap="nowrap"
        >
          <Group
            gap={8}
            wrap="nowrap"
