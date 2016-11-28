import { Button } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { useTranslation } from 'react-i18next'

import { RegistryFormModal } from './RegistryFormModal'

import type { Registry } from '@matrixhub/api-ts/v1alpha1/registry.pb'

interface EditRegistryActionProps {
  registry: Registry
  disabled?: boolean
}

export function EditRegistryAction({
  registry,
  disabled,
}: EditRegistryActionProps) {
  const { t } = useTranslation()
  const [opened, {
    open,
    close,
  }] = useDisclosure(false)

