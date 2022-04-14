import {
  Button,
  Checkbox,
  PasswordInput,
  Stack,
  TextInput,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { IconUserPlus } from '@tabler/icons-react'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { ModalWrapper } from '@/shared/components/ModalWrapper'
import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form'

import { createUserMutationOptions } from '../users.mutation'
import { createUserFormSchema } from '../users.schema'

interface CreateUserModalProps {
  opened: boolean
  onClose: () => void
}

function CreateUserModal({
  opened,
  onClose,
}: CreateUserModalProps) {
  const { t } = useTranslation()
  const mutation = useMutation(createUserMutationOptions())
  const form = useForm({
    defaultValues: {
      username: '',
      password: '',
      confirmPassword: '',
      isAdmin: false,
    },
    validators: {
      onSubmit: createUserFormSchema(),
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        username: value.username,
        password: value.password,
        isAdmin: value.isAdmin,
      })
      onClose()
    },
  })

  return (
    <ModalWrapper
      size="xs"
      opened={opened}
      onClose={onClose}
      title={t('routes.admin.users.createModal.title')}
      confirmLoading={form.state.isSubmitting}
      onConfirm={form.handleSubmit}
    >

      <Stack gap="lg">
        <form.Field
          name="username"
        >
          {field => (
            <TextInput
              required
              autoComplete="username"
              label={t('routes.admin.users.form.username')}
              value={field.state.value}
              onChange={event => field.handleChange(event.currentTarget.value)}
              onBlur={field.handleBlur}
              error={fieldError(field)}
            />
          )}
