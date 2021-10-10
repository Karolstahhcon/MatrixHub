import {
  SimpleGrid,
  Skeleton,
} from '@mantine/core'
import { useMemo } from 'react'

import { EmptyStatePrompt } from '@/shared/components/EmptyStatePrompt'

import type { ReactNode } from 'react'

export interface ResourceCardGridProps {
  loading?: boolean
  skeletonCount?: number
  skeletonHeight?: number
  emptyTitle?: ReactNode
  emptyHelper?: ReactNode
  children: ReactNode[]
}

const DEFAULT_SKELETON_COUNT = 6
const DEFAULT_SKELETON_HEIGHT = 116

export function ResourceCardGrid({
  loading,
  skeletonCount = DEFAULT_SKELETON_COUNT,
