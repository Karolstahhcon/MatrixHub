import { Button } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { useTranslation } from 'react-i18next'

import { ReplicationFormModal } from './ReplicationFormModal'

import type { SyncPolicyItem } from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'

interface EditReplicationActionProps {
  syncPolicy: SyncPolicyItem
  disabled?: boolean
}

export function EditReplicationAction({
  syncPolicy,
  disabled,
}: EditReplicationActionProps) {
  const { t } = useTranslation()
