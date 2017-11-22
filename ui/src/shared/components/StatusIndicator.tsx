import {
  Box,
  Group,
  Text,
  type TextProps,
} from '@mantine/core'

import type { ReactNode } from 'react'

export type StatusIndicatorType = 'danger' | 'info' | 'neutral' | 'success' | 'warning'

export interface StatusIndicatorProps {
  label: ReactNode
  type?: StatusIndicatorType
  dotSize?: number
  textProps?: TextProps
}

const statusIndicatorColors: Record<StatusIndicatorType, string> = {
  danger: 'var(--mantine-color-error)',
  info: 'var(--mantine-primary-color-filled)',
  neutral: 'var(--mantine-color-dimmed)',
  success: 'var(--mantine-color-green-filled)',
  warning: 'var(--mantine-color-yellow-filled)',
