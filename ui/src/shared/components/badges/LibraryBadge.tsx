import {
  type PolymorphicFactory,
  polymorphicFactory,
  useProps,
} from '@mantine/core'

import PytorchIcon from '@/assets/svgs/pytorch.svg?react'
import { BaseBadge, type BaseBadgeProps } from '@/shared/components/badges/BaseBadge'

import type { ReactNode } from 'react'

export interface LibraryBadgeProps extends Omit<BaseBadgeProps, 'icon' | 'label'> {
  library: string
}

export type LibraryBadgeFactory = PolymorphicFactory<{
  props: LibraryBadgeProps
  defaultComponent: 'div'
  defaultRef: HTMLDivElement
  stylesNames: 'root' | 'label'
