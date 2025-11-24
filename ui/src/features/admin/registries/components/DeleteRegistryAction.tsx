import {
  Button,
  Stack,
  Text,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { ModalWrapper } from '@/shared/components/ModalWrapper'

import { deleteRegistryMutationOptions } from '../registries.mutation'
import { getRegistryDisplayName } from '../registries.utils'

import type { Registry } from '@matrixhub/api-ts/v1alpha1/registry.pb'

interface DeleteRegistryActionProps {
  registry: Registry
  disabled?: boolean
}

interface DeleteRegistryModalProps {
  opened: boolean
  registry: Registry
  loading?: boolean
  onClose: () => void
