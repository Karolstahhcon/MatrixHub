import {
  Textarea,
  TextInput,
} from '@mantine/core'
import { useMutation } from '@tanstack/react-query'
import dayjs from 'dayjs'
import { useTranslation } from 'react-i18next'

import { ExpireAtField } from '@/features/profile/components/ExpireAtField'
import { createSshKeyMutationOptions } from '@/features/profile/profile.mutation'
import { createSshKeySchema } from '@/features/profile/profile.schema'
import { ModalWrapper } from '@/shared/components/ModalWrapper'
import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form.ts'

interface CreateSshKeyModalProps {
  opened: boolean
  onClose: () => void
}

export function CreateSshKeyModal({
  opened,
  onClose,
}: CreateSshKeyModalProps) {
  const { t } = useTranslation()
  const mutation = useMutation(createSshKeyMutationOptions())
  const sshKeySchema = createSshKeySchema(t)

  const form = useForm({
    defaultValues: {
      name: '',
      publicKey: '',
      expireAt: '',
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        ...value,
        expireAt: value.expireAt ? String(dayjs(value.expireAt).unix()) : '',
      })
      handleClose()
    },
  })

  const handleClose = () => {
    onClose()
    form.reset()
    mutation.reset()
  }

  return (
    <ModalWrapper
      title={t('profile.sshKey.create.title')}
      opened={opened}
      onClose={handleClose}
