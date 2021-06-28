import {
  Button,
  Stack,
  Text,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { ModalWrapper } from '@/shared/components/ModalWrapper'

import { deleteReplicationMutationOptions } from '../replications.mutation'
import { getReplicationDisplayName } from '../replications.utils'

import type { SyncPolicyItem } from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'

interface DeleteReplicationActionProps {
  syncPolicy: SyncPolicyItem
  disabled?: boolean
}

interface DeleteReplicationModalProps {
  opened: boolean
  syncPolicy: SyncPolicyItem
  loading?: boolean
  onClose: () => void
  onConfirm: () => void
}

function DeleteReplicationModal({
  opened,
  syncPolicy,
  loading,
  onClose,
  onConfirm,
}: DeleteReplicationModalProps) {
  const { t } = useTranslation()

  return (
    <ModalWrapper
      opened={opened}
      onClose={onClose}
      type="danger"
      size="sm"
      closeOnClickOutside={false}
      title={t('routes.admin.replications.deleteModal.title')}
      confirmLoading={loading}
      onConfirm={onConfirm}
    >
      <Stack gap="sm">
        <Text size="sm">
