import { Button } from '@mantine/core'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { switchReplicationMutationOptions } from '../replications.mutation'

import type { SyncPolicyItem } from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'

interface ToggleReplicationActionProps {
  syncPolicy: SyncPolicyItem
  disabled?: boolean
}

export function ToggleReplicationAction({
  syncPolicy,
  disabled,
}: ToggleReplicationActionProps) {
  const { t } = useTranslation()
  const mutation = useMutation(switchReplicationMutationOptions())
  const nextIsDisabled = !syncPolicy.isDisabled

  const handleToggle = async () => {
    if (syncPolicy.id == null) {
      return
    }

    await mutation.mutateAsync({
      syncPolicyId: syncPolicy.id,
      isDisabled: nextIsDisabled,
