import { Button } from '@mantine/core'
import {
  SyncTaskStatus,
  type SyncTask,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { stopReplicationExecutionMutationOptions } from '../executions.mutation'
