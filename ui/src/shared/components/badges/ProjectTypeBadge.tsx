import {
  Badge,
  type BadgeProps,
} from '@mantine/core'
import { ProjectType } from '@matrixhub/api-ts/v1alpha1/project.pb'
import { useTranslation } from 'react-i18next'

interface ProjectTypeBadgeProps extends Omit<BadgeProps, 'children'> {
  type?: ProjectType
}

export function ProjectTypeBadge({
  type,
  ...badgeProps
}: ProjectTypeBadgeProps) {
  const { t } = useTranslation()

  if (!type) {
