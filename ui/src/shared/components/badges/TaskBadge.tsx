import {
  type PolymorphicFactory,
  polymorphicFactory,
  useProps,
} from '@mantine/core'
import { IconPhotoUp } from '@tabler/icons-react'

import { BaseBadge, type BaseBadgeProps } from '@/shared/components/badges/BaseBadge'

import type { ReactNode } from 'react'

export interface TaskBadgeProps extends Omit<BaseBadgeProps, 'icon' | 'label'> {
  task: string
}

export type TaskBadgeFactory = PolymorphicFactory<{
  props: TaskBadgeProps
  defaultComponent: 'div'
  defaultRef: HTMLDivElement
  stylesNames: 'root' | 'label'
}>

const TASK_ICON_ENTRIES = [
  {
    matcher: /(image|vision|photo)/i,
    icon: <IconPhotoUp size={16} color="var(--mantine-color-blue-4)" />,
