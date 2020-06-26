import { Button } from '@mantine/core'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { syncReplicationMutationOptions } from '../replications.mutation'

import type { SyncPolicyItem } from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'

interface SyncReplicationActionProps {
  syncPolicy: SyncPolicyItem
  disabled?: boolean
}

export function SyncReplicationAction({
  syncPolicy,
  disabled,
}: SyncReplicationActionProps) {
  const { t } = useTranslation()
  const mutation = useMutation(syncReplicationMutationOptions())

  const handleSync = async () => {
