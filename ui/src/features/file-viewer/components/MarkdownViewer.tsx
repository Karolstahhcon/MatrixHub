import {
  Alert, Box, Center, Loader,
} from '@mantine/core'
import { useEffect, useState } from 'react'

import { decodeHash } from '../markdown/hash'
import { renderMarkdown } from '../markdown/renderer'

import 'github-markdown-css/github-markdown-light.css'
import './MarkdownViewer.css'

interface MarkdownViewerProps {
  content: string
}

export function MarkdownViewer({ content }: MarkdownViewerProps) {
  const [result, setResult] = useState<{
    content: string
    value: Error | string
  } | null>(null)

  useEffect(() => {
    let cancelled = false

    renderMarkdown(content).then(
      html => !cancelled && setResult({
        content,
        value: html,
      }),
      error => !cancelled && setResult({
        content,
        value: error instanceof Error ? error : new Error(String(error)),
      }),
    )

    return () => {
      cancelled = true
    }
  }, [content])

  const html = result?.content === content && typeof result.value === 'string' ? result.value : null
