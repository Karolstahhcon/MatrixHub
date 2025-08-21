import {
  Card,
  Grid,
  Stack,
  Text,
} from '@mantine/core'
import {
  SyncTaskStatus,
  TriggerType,
  type SyncTask,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { useTranslation } from 'react-i18next'

import { StatusIndicator } from '@/shared/components/StatusIndicator'
import { formatDateTime } from '@/shared/utils/date'

interface ExecutionInfoProps {
  task?: SyncTask
}

export function ExecutionInfo({ task }: ExecutionInfoProps) {
  const { t } = useTranslation()

  const infoItems = [
    {
      label: 'ID',
      value: task?.id ?? '-',
    },
    {
      label: t('routes.admin.replications.executions.table.status'),
      value: (
        <ExecutionStatusIndicator status={task?.status} />
      ),
    },
    {
      label: t('routes.admin.replications.executions.table.triggerType'),
      value: task?.triggerType === TriggerType.TRIGGER_TYPE_SCHEDULED
        ? t('routes.admin.replications.trigger.scheduled')
        : t('routes.admin.replications.trigger.manual'),
