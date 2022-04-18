import {
  Box, Code, Skeleton,
} from '@mantine/core'
import {
  useEffect, useRef, useState,
} from 'react'

import { CopyValueButton } from '@/shared/components/CopyValueButton'
import { highlightCode } from '@/shared/utils/shiki'

import classes from './ShikiCodeBlock.module.css'

interface ShikiCodeBlockProps {
  code: string
  lang: string
  /** Copy button in the top-right corner. Defaults to true. */
  copyable?: boolean
}

/**
 * Syntax-highlighted, copyable code block. Renders a plain `<Code>` while
 * Shiki (and the requested grammar) loads, then swaps in the highlighted HTML.
 */
export function ShikiCodeBlock({
  code,
  lang,
  copyable = true,
}: ShikiCodeBlockProps) {
  const [highlighted, setHighlighted] = useState<{
    code: string
    html: string
  } | null>(null)

  useEffect(() => {
    let cancelled = false

    highlightCode(code, lang).then(
      html => !cancelled && setHighlighted({
        code,
        html,
      }),
      () => !cancelled && setHighlighted({
        code,
        html: '',
      }),
    )

    return () => {
      cancelled = true
