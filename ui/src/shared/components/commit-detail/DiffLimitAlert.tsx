import {
  Alert,
  Anchor,
  Text,
} from '@mantine/core'
import { type MouseEvent } from 'react'
import {
  Trans,
  useTranslation,
} from 'react-i18next'

interface DiffLimitAlertProps {
  maxFiles: number
  onOpenRawDiff: (event: MouseEvent<HTMLAnchorElement>) => void
}

export function DiffLimitAlert({
  maxFiles,
  onOpenRawDiff,
}: DiffLimitAlertProps) {
  const { t } = useTranslation()
