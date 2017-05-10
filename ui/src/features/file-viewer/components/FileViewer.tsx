import {
  Alert,
  Center,
  Loader,
  Paper,
} from '@mantine/core'
import {
  lazy,
  Suspense,
  useState,
} from 'react'

import { classifyFile, getMonacoLanguage } from '../utils'
import { FallbackCard } from './FallbackCard'
import { FileToolbar } from './FileToolbar'

import type { FileViewerFile } from '../types'

type ViewMode = 'preview' | 'code'

const MarkdownViewer = lazy(() =>
