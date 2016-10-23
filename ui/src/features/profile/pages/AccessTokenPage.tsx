import {
  Alert,
  Button,
  Flex,
  Group,
  rem,
  Stack,
  Text,
  TextInput,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { CurrentUser, type CreateAccessTokenRequest } from '@matrixhub/api-ts/v1alpha1/current_user.pb'
import {
  IconInfoCircle,
  IconKey,
} from '@tabler/icons-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import dayjs from 'dayjs'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import z from 'zod'

import { AccessTokenTable } from '@/features/profile/components/AccessTokenTable'
import { ExpireAtField } from '@/features/profile/components/ExpireAtField'
import { profileKeys, useAccessTokens } from '@/features/profile/profile.query'
import { expireAtSchema } from '@/features/profile/profile.schema.ts'
import { CopyValueButton } from '@/shared/components/CopyValueButton'
import { ModalWrapper } from '@/shared/components/ModalWrapper'
import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form.ts'

export function AccessTokenPage() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const expireAtFieldSchema = expireAtSchema(t)

  const {
    data, isFetching,
  } = useAccessTokens()
  const tokens = data?.items ?? []

  const [hintVisible, setHintVisible] = useState(true)
  const [createOpened, {
    open: openCreate, close: closeCreate,
  }] = useDisclosure(false)

  const [copyOpened, {
    open: openCopy, close: closeCopy,
  }] = useDisclosure(false)

  const handleRefresh = () => {
    queryClient.invalidateQueries({ queryKey: profileKeys.accessTokens })
  }

  const nameSchema = z.string().min(1, { error: t('common.validation.fieldRequired', { field: t('profile.tokenName') }) })

  const [newToken, setNewToken] = useState<string>('')

  const {
    mutate: createToken, isPending: isCreating,
  } = useMutation({
    mutationFn: (value: CreateAccessTokenRequest) => CurrentUser.CreateAccessToken(value),
    meta: {
      successMessage: t('profile.tokenCreated'),
      invalidates: [profileKeys.accessTokens],
    },
    onSuccess: (res) => {
      setNewToken(res.token ?? '')
      handleCreateClose()

      if (res.token) {
        openCopy()
      }
    },
  })

  const form = useForm({
    defaultValues: {
      name: '',
      expireAt: '',
    },
    onSubmit: ({ value }) => {
      createToken({
        ...value,
        expireAt: value.expireAt ? String(dayjs(value.expireAt).unix()) : '',
      })
    },
  })

  const handleCreateClose = () => {
    closeCreate()
    form.reset()
  }

  const handleCopyClose = () => {
    closeCopy()
    setNewToken('')
  }

  return (
    <Stack gap="sm">
      {hintVisible && (
        <Alert
          icon={<IconInfoCircle size={20} />}
          variant="light"
          color="cyan"
          withCloseButton
          onClose={() => setHintVisible(false)}
          styles={{ icon: { marginRight: 6 } }}
        >
          <Text size="sm">{t('profile.tokenHint')}</Text>
        </Alert>
      )}

      <AccessTokenTable
        tokens={tokens}
        fetching={isFetching}
        onRefresh={handleRefresh}
        toolbarExtra={(
          <Button
            leftSection={<IconKey size={16} />}
            onClick={openCreate}
          >
            {t('profile.createToken')}
          </Button>
        )}
      />

      <ModalWrapper
        title={t('profile.createToken')}
        opened={createOpened}
        onClose={handleCreateClose}
