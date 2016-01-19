import {
  SimpleGrid,
  Skeleton,
} from '@mantine/core'
import { useMemo } from 'react'

import { EmptyStatePrompt } from '@/shared/components/EmptyStatePrompt'

import type { ReactNode } from 'react'

export interface ResourceCardGridProps {
  loading?: boolean
