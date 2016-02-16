import {
  Group,
  Text,
  Tooltip,
} from '@mantine/core'
import { useTranslation } from 'react-i18next'

import IconQuestion from '@/assets/svgs/question.svg?react'

import type { TooltipProps } from '@mantine/core'
import type { ReactNode } from 'react'

interface FieldHintLabelProps {
  label: ReactNode
  hint: ReactNode
  tooltipProps?: Omit<TooltipProps, 'children' | 'label'>
}

const DEFAULT_TOOLTIP_WIDTH = 240

export function FieldHintLabel({
  label,
  hint,
  tooltipProps,
}: FieldHintLabelProps) {
  const { t } = useTranslation()

  return (
    <Group
      component="span"
      gap={4}
      align="center"
      wrap="nowrap"
      style={{
