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

