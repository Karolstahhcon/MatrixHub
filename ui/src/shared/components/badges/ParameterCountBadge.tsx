import {
  type PolymorphicFactory,
  polymorphicFactory,
  useProps,
} from '@mantine/core'
import { IconBinaryTree } from '@tabler/icons-react'

import { BaseBadge, type BaseBadgeProps } from '@/shared/components/badges/BaseBadge'

export interface ParameterCountBadgeProps extends Omit<BaseBadgeProps, 'icon' | 'label'> {
  parameterCount: string
}

export type ParameterCountBadgeFactory = PolymorphicFactory<{
  props: ParameterCountBadgeProps
  defaultComponent: 'div'
  defaultRef: HTMLDivElement
  stylesNames: 'root' | 'label'
}>

export const ParameterCountBadge = polymorphicFactory<ParameterCountBadgeFactory>((_props, ref) => {
  const {
    parameterCount,
    ...badgeProps
  } = useProps('ParameterCountBadge', {}, _props)

  return (
    <BaseBadge
