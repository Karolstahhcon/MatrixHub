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
