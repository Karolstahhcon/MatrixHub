import {
  Group,
  rem,
  Text,
} from '@mantine/core'
import { Link } from '@tanstack/react-router'
import { Fragment, type ReactNode } from 'react'
import { useTranslation } from 'react-i18next'

import { BaseBadge } from '@/shared/components/badges/BaseBadge'
import { CopyValueButton } from '@/shared/components/CopyValueButton.tsx'

import type { ResourceBadge, ResourceMetaItem } from '@/shared/components/resource-card/BaseCard'

interface ResourceDetailHeaderProps {
  projectId: string
