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
  defaultRef: HTMLDivElement
  stylesNames: 'root' | 'label'
}>

type BadgeStylesObject = Exclude<NonNullable<BadgeProps['styles']>, (...args: never[]) => unknown>

const DEFAULT_BADGE_STYLES: BadgeStylesObject = {
  root: {
    backgroundColor: 'var(--mantine-color-gray-1)',
    paddingInline: 12,
    border: '1px solid var(--mantine-color-gray-1)',
  },
  label: {
    paddingInline: 0,
    textTransform: 'none',
  },
}

const defaultProps: Partial<BaseBadgeProps> = {
  h: 24,
  radius: 'lg',
}
