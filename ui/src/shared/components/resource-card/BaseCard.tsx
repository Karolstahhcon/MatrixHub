import {
  Card,
  Group,
  Stack,
  Text,
} from '@mantine/core'
import { IconDots } from '@tabler/icons-react'
import { Fragment } from 'react'

import { BaseBadge } from '@/shared/components/badges/BaseBadge.tsx'

import classes from './BaseCard.module.css'

import type { ReactNode } from 'react'

export interface ResourceBadge {
  key: string | number
  content?: ReactNode
  icon?: ReactNode
  label?: ReactNode
}

export interface ResourceMetaItem {
  key: string
  label: string
  icon?: ReactNode
  value: ReactNode
}

interface BaseCardProps {
  title?: ReactNode
  badges?: ResourceBadge[]
  maxBadges?: number
  metaItems?: ResourceMetaItem[]
  renderMeta?: () => ReactNode
  renderRoot?: (props: Record<string, unknown>) => ReactNode
}

const EMPTY_BADGES: ResourceBadge[] = []
const EMPTY_META_ITEMS: ResourceMetaItem[] = []
const DEFAULT_MAX_BADGES = 3

export function BaseCard({
  title,
  badges = EMPTY_BADGES,
