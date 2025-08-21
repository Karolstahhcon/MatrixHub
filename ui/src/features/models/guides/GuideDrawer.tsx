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
