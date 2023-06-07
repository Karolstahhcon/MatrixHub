import {
  Box, Center, Loader,
} from '@mantine/core'
import Editor from '@monaco-editor/react'

import classes from './FileViewer.module.css'

interface CodeViewerProps {
  content: string
  language: string
}

