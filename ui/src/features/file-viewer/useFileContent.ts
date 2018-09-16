import {
  useEffect,
  useState,
} from 'react'

import { classifyFile } from './utils'

import type { File } from '@matrixhub/api-ts/v1alpha1/model.pb'

interface FileContentState {
  content: string | undefined
  isLoading: boolean
  error: Error | null
}

/**
 * Fetches the text content of a file from `file.url`.
 * Only fetches when the file is classified as text (markdown or code).
 * Binary files are skipped — the caller should show a fallback card instead.
 */
export function useFileContent(file: File | undefined): FileContentState {
