import {
  Badge,
  type BadgeProps,
  Group,
  type PolymorphicFactory,
  polymorphicFactory,
  Text,
  useProps,
} from '@mantine/core'

import type { ReactNode } from 'react'

export interface BaseBadgeProps extends Omit<BadgeProps, 'children'> {
  icon?: ReactNode
  label?: ReactNode
  children?: ReactNode
}

export type BaseBadgeFactory = PolymorphicFactory<{
  props: BaseBadgeProps
  defaultComponent: 'div'
