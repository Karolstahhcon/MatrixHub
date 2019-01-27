import {
  Text,
  Tooltip,
} from '@mantine/core'
import {
  useRef,
  useState,
} from 'react'

import type {
  TextProps,
  TooltipProps,
} from '@mantine/core'
import type { MouseEvent, ReactNode } from 'react'

export interface TruncatedTextProps extends Omit<TextProps, 'truncate' | 'children'> {
  /** Full value. Used both as rendered content and as the tooltip label. */
  value: ReactNode
  /** Tooltip label when it should differ from `value` (e.g. `value` is a link). */
  tooltipLabel?: ReactNode
  /** Rendered content when it should differ from `value` (e.g. an `Anchor`). */
  children?: ReactNode
  /** Extra props forwarded to the tooltip. */
  tooltipProps?: Omit<TooltipProps, 'children' | 'label'>
  /** Composed with the internal overflow check rather than replacing it. */
  onMouseEnter?: (event: MouseEvent<HTMLDivElement>) => void
}

/**
 * Single-line cell text that truncates on overflow and reveals the full value
 * on hover. The tooltip only activates when the text actually overflows, so
 * short values do not get a redundant hover card.
 */
export function TruncatedText({
  value,
  tooltipLabel,
  children,
  tooltipProps,
  onMouseEnter,
  ...textProps
}: TruncatedTextProps) {
