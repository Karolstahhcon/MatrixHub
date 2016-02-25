import { useTimeout } from '@mantine/hooks'
import { useState } from 'react'

export interface UseClipboardOptions {
  /** Time in ms after which the copied state resets. */
  timeout?: number
  /** Whether to fall back to document.execCommand('copy'). */
  legacy?: boolean
}

export interface UseClipboardReturnValue {
  copied: boolean
  error: Error | null
  copy: (value: string) => void
  reset: () => void
}

function legacyCopy(value: string) {
  const textarea = document.createElement('textarea')

  textarea.value = value
  textarea.style.position = 'absolute'
  textarea.style.opacity = '0'
  textarea.setAttribute('readonly', '')
  document.body.appendChild(textarea)
  textarea.select()

  try {
    return document.execCommand('copy')
  } finally {
    textarea.remove()
  }
}

export function useClipboard({
  timeout = 2000,
  legacy = false,
}: UseClipboardOptions = {}): UseClipboardReturnValue {
  const [copied, setCopied] = useState(false)
  const [error, setError] = useState<Error | null>(null)
  const copiedTimeout = useTimeout(() => setCopied(false), timeout, { autoInvoke: false })

  const handleCopyResult = (succeeded: boolean) => {
    setCopied(succeeded)

    if (succeeded) {
      copiedTimeout.start()
    }
  }

  const copy = (value: string) => {
    setError(null)

    void (async () => {
      let clipboardError: unknown
      const isClipboardApiSupported = window.isSecureContext && 'clipboard' in navigator
      let useLegacy = !isClipboardApiSupported || legacy

      if (!useLegacy) {
        try {
          await navigator.clipboard.writeText(value)
          handleCopyResult(true)

          return
        } catch (copyError) {
          clipboardError = copyError
          useLegacy = true
        }
      }

      if (useLegacy) {
        try {
