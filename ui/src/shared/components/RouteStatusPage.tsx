import {
  Box, Button, Stack, Text, Title,
} from '@mantine/core'
import { Link } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

interface RouteStatusPageProps {
  code: 403 | 404
  title?: string
  description?: string
  fullScreen?: boolean
}

export function RouteStatusPage({
  code,
  title,
  description,
  fullScreen = false,
}: RouteStatusPageProps) {
  const { t } = useTranslation()

  return (
    <Box
      h={fullScreen ? '100vh' : '100%'}
      pos="relative"
      style={{
        overflow: 'hidden',
      }}
    >
      <Box
        aria-hidden
        pos="absolute"
        fw="900"
        c="gray.1"
        lh="1"
