import {
  Alert,
  Button,
  Center, Group, Stack, Text,
} from '@mantine/core'
import { IconAlertCircle, IconRefresh } from '@tabler/icons-react'
import { useQueryClient } from '@tanstack/react-query'
import { useRouter } from '@tanstack/react-router'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { getErrorMessage } from '@/queryClient'

interface RouterErrorComponentProps {
  error: unknown
}

export function RouterErrorComponent({ error }: RouterErrorComponentProps) {
