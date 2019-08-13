import {
  Button,
  PasswordInput,
  Stack,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { ModalWrapper } from '@/shared/components/ModalWrapper'
import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form'

import { resetUserPasswordMutationOptions } from '../users.mutation'
import { resetUserPasswordFormSchema } from '../users.schema'

import type { User } from '@matrixhub/api-ts/v1alpha1/user.pb'

interface ResetUserPasswordActionProps {
  user: User
  disabled?: boolean
}

interface ResetUserPasswordModalProps {
  opened: boolean
  user: User
  onClose: () => void
}

function ResetUserPasswordModal({
  opened,
  user,
  onClose,
}: ResetUserPasswordModalProps) {
  const { t } = useTranslation()
  const mutation = useMutation(resetUserPasswordMutationOptions())
  const form = useForm({
    defaultValues: {
      password: '',
      confirmPassword: '',
    },
    validators: {
      onChange: resetUserPasswordFormSchema(),
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        id: user.id,
        password: value.password,
      })
      onClose()
    },
  })

  return (
    <ModalWrapper
      opened={opened}
      onClose={onClose}
      title={t('routes.admin.users.resetPasswordModal.title')}
      confirmLoading={form.state.isSubmitting}
      onConfirm={form.handleSubmit}
      size="xs"
    >

      <Stack gap="lg">
