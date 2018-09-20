import {
  Alert,
  Badge,
  Box,
  Button,
  Drawer,
  Group,
  rem,
  Stack,
  Text,
} from '@mantine/core'
import { IconInfoCircle } from '@tabler/icons-react'
import { useTranslation } from 'react-i18next'

import { ShikiCodeBlock } from '@/shared/components/ShikiCodeBlock'

import classes from './GuideDrawer.module.css'

import type { ReactNode } from 'react'

interface GuideDrawerProps {
  opened: boolean
  title: ReactNode
  /** Short explanation shown in an info alert above the steps. */
  intro: ReactNode
  /** Rendered at the left of the footer, opposite the confirm button. */
  footerExtra?: ReactNode
  onClose: () => void
  children: ReactNode
}

/**
 * Right-side drawer for step-by-step command guides (use / upload / download).
 * The body scrolls; the footer stays pinned with a confirm button.
 */
export function GuideDrawer({
  opened,
  title,
  intro,
  footerExtra,
  onClose,
  children,
}: GuideDrawerProps) {
  const { t } = useTranslation()

  return (
    <Drawer
      opened={opened}
      onClose={onClose}
      position="right"
