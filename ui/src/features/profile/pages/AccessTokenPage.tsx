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
