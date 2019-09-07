import {
  Card,
  Grid,
  Group,
  Stack,
  Text,
  ThemeIcon,
  rem,
} from '@mantine/core'
import { type SyncTask } from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import {
  IconCircleCheck,
  IconAlertCircle,
  IconMessageCircleBolt,
  IconAlertTriangle,
} from '@tabler/icons-react'
import { useTranslation } from 'react-i18next'

interface ExecutionStatsCardsProps {
  task?: SyncTask
}

export function ExecutionStatsCards({ task }: ExecutionStatsCardsProps) {
  const { t } = useTranslation()

  const safeNumber = (value: string | number | undefined) => Number(value ?? 0) || 0
  const total = safeNumber(task?.totalItems)
  const success = safeNumber(task?.successfulItems)
  const failed = safeNumber(task?.failedItems)
  const stopped = safeNumber(task?.stoppedItems)
  const running = Math.max(0, total - success - failed - stopped)

  const stats = [
    {
      label: t('routes.admin.replications.executions.detail.status.succeeded'),
      value: success,
      icon: IconCircleCheck,
      color: 'var(--mantine-color-teal-6)',
      bgColor: 'var(--mantine-color-teal-light)',
    },
    {
      label: t('routes.admin.replications.executions.detail.status.failed'),
      value: failed,
      icon: IconAlertCircle,
      color: 'var(--mantine-color-red-6)',
      bgColor: 'var(--mantine-color-red-light)',
    },
    {
      label: t('routes.admin.replications.executions.detail.status.running'),
      value: running,
      icon: IconMessageCircleBolt,
      color: 'var(--mantine-primary-color-filled)',
      bgColor: 'var(--mantine-primary-color-light)',
    },
    {
      label: t('routes.admin.replications.executions.detail.status.stopped'),
      value: stopped,
      icon: IconAlertTriangle,
      color: 'var(--mantine-color-yellow-6)',
      bgColor: 'var(--mantine-color-orange-light)',
    },
  ]

  return (
    <Grid gutter="md">
      {stats.map(stat => (
