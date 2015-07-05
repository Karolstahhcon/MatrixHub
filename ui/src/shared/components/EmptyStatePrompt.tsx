import {
  Center,
  Stack,
  Text,
} from '@mantine/core'
import { useTranslation } from 'react-i18next'

import type { ReactNode } from 'react'

export interface EmptyStatePromptProps {
  title?: ReactNode
  helper?: ReactNode
}
