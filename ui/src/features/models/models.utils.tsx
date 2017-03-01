import {
  Category, type Label, type Model,
} from '@matrixhub/api-ts/v1alpha1/model.pb'
import {
  IconClock, IconCube, IconApiApp,
} from '@tabler/icons-react'
import humanFormat from 'human-format'

import i18n from '@/i18n'
import { BaseBadge } from '@/shared/components/badges/BaseBadge'
import { LibraryBadge } from '@/shared/components/badges/LibraryBadge'
import { ParameterCountBadge } from '@/shared/components/badges/ParameterCountBadge'
import { TaskBadge } from '@/shared/components/badges/TaskBadge'
import { formatDateTime } from '@/shared/utils/date'
import { formatStorageSize } from '@/shared/utils/format'

import type { ResourceBadge, ResourceMetaItem } from '@/shared/components/resource-card/BaseCard'
import type { TFunction } from 'i18next'

export function buildModelTitle(model: Model, projectId: string) {
  const projectName = model.project ?? projectId
  const modelName = model.name?.trim()

  return `${projectName} / ${modelName || '-'}`
}

export function buildModelBadges(
  model: Model,
  options: {
    taskCategory: Category
    libraryCategory: Category
  } = {
    taskCategory: Category.TASK,
    libraryCategory: Category.LIBRARY,
  },
): ResourceBadge[] {
  const badges: ResourceBadge[] = []
  const taskLabels = getLabelsByCategory(model.labels, options.taskCategory)
  const libraryLabels = getLabelsByCategory(model.labels, options.libraryCategory)

  for (const name of taskLabels) {
    badges.push({
      key: `task-${name}`,
      content: <TaskBadge task={name} maw={132} />,
    })
  }
