import {
  Button,
  Checkbox,
  Group,
  PasswordInput,
  Select,
  Stack,
  Text,
  TextInput,
  Textarea,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { RegistryType } from '@matrixhub/api-ts/v1alpha1/registry.pb'
import { useStore } from '@tanstack/react-form'
import { useMutation } from '@tanstack/react-query'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { FieldHintLabel } from '@/shared/components/FieldHintLabel'
import { ModalWrapper } from '@/shared/components/ModalWrapper'
import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form'

import {
  buildCreateRegistryRequest,
  buildPingRegistryRequest,
  buildUpdateRegistryRequest,
  createRegistryFormSchema,
  editRegistryFormSchema,
  getRegistryDefaultUrl,
  getRegistryFormValues,
  getRegistryUrlCopyKeys,
  isRegistryUrlPristine,
  isLikelyRepositoryUrl,
  REGISTRY_DESCRIPTION_MAX_LENGTH,
  registryProviderLabelKeys,
  registryDescriptionSchema,
  registryNameSchema,
  registryUrlSchema,
  sanitizeRegistryName,
  supportedRegistryTypes,
} from '../registries.form'
import {
  createRegistryMutationOptions,
  pingRegistryMutationOptions,
  updateRegistryMutationOptions,
} from '../registries.mutation'

import type { Registry } from '@matrixhub/api-ts/v1alpha1/registry.pb'

export interface RegistryFormModalProps {
  mode: 'create' | 'edit'
  opened: boolean
  registry?: Registry
  onClose: () => void
}

export function RegistryFormModal({
  mode,
  opened,
  registry,
  onClose,
}: RegistryFormModalProps) {
  const { t } = useTranslation()
  const submitMutation = useMutation(
    mode === 'create'
      ? createRegistryMutationOptions()
      : updateRegistryMutationOptions(),
  )
  const pingMutation = useMutation(pingRegistryMutationOptions())
  const form = useForm({
    defaultValues: getRegistryFormValues(registry),
    validators: {
      onSubmit: mode === 'create'
        ? createRegistryFormSchema
        : editRegistryFormSchema,
    },
    onSubmit: async ({ value }) => {
      if (mode === 'create') {
        await submitMutation.mutateAsync(buildCreateRegistryRequest(value))
      } else if (registry?.id != null) {
        await submitMutation.mutateAsync(buildUpdateRegistryRequest(registry.id, value))
