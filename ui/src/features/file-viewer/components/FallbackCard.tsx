import {
  Box,
  Stack,
  Text,
} from '@mantine/core'
import { useTranslation } from 'react-i18next'

import { formatStorageSize } from '@/shared/utils/format'

import type { FileViewerFile } from '../types'

interface FallbackCardProps {
  file: FileViewerFile
