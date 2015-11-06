import {
  Button,
  Group,
  Stack,
  Text,
} from '@mantine/core'
import { notifications } from '@mantine/notifications'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { ModalWrapper } from '@/shared/components/ModalWrapper'

import { batchDeleteUsersMutationOptions } from '../users.mutation'
import { getUserDisplayName } from '../users.utils'

import type { User } from '@matrixhub/api-ts/v1alpha1/user.pb'

interface BatchDeleteUsersModalProps {
  opened: boolean
  users: readonly User[]
  onClose: () => void
  onSuccess: () => void
}

export function BatchDeleteUsersModal({
  opened,
  users,
  onClose,
  onSuccess,
}: BatchDeleteUsersModalProps) {
  const { t } = useTranslation()
  const mutation = useMutation(batchDeleteUsersMutationOptions())
  const usernames = users.map(getUserDisplayName).join(', ')

  const handleDelete = async () => {
    const res = await mutation.mutateAsync(users)

    if (res && res.partial) {
      notifications.show({
        color: 'yellow',
        message: t('routes.admin.users.notifications.batchDeletePartialError', {
          successCount: res.successCount,
          failureCount: res.failureCount,
        }),
      })
    } else {
