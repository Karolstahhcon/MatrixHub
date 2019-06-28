import {
  ActionIcon,
  Checkbox,
  Group,
  Input,
  Select,
  Stack,
  Switch,
  Text,
  TextInput,
  Tooltip,
} from '@mantine/core'
import {
  IconExternalLink, IconInfoCircle, IconPlus, IconRefresh,
} from '@tabler/icons-react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { useEffect, useEffectEvent } from 'react'
import { Trans, useTranslation } from 'react-i18next'

import { allRegistriesQueryOptions } from '@/features/admin/registries/registries.query'
import { useCurrentUser } from '@/features/auth/auth.query'
import AnchorLink from '@/shared/components/AnchorLink'
import { FieldHintLabel } from '@/shared/components/FieldHintLabel.tsx'
import { ModalWrapper } from '@/shared/components/ModalWrapper'
import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form'

import { createProjectMutationOptions } from '../projects.mutation'
import {
  organizationSchema, projectNameSchema, registryIdSchema,
} from '../projects.schema'
import { ProjectTypeHintLabel } from './ProjectTypeHintLabel'

import type { MantineSize } from '@mantine/core'
import type { ReactNode } from 'react'

export interface CreateProjectModalProps {
  opened: boolean
  onClose: () => void
  onCreated?: (projectName: string) => void | Promise<void>
}

/**
 * Link to registry creation. Opened in a new tab so the half-filled project
 * form stays untouched, with the usual icon marking the tab switch.
 */
function RegistryCreateLink({
  children,
  size,
}: {
  children?: ReactNode
  size?: MantineSize
}) {
  return (
    <AnchorLink
      to="/admin/registries"
      search={{ create: true }}
      target="_blank"
      size={size}
      inherit={!size}
    >
      <Group component="span" gap={2} align="center" wrap="nowrap" display="inline-flex">
        {children}
        <IconExternalLink size={14} />
      </Group>
    </AnchorLink>
  )
}

export function CreateProjectModal({
  opened,
  onClose,
  onCreated,
}: CreateProjectModalProps) {
  const { t } = useTranslation()
  const mutation = useMutation(createProjectMutationOptions())
  const { data: currentUser } = useCurrentUser()

  const form = useForm({
    defaultValues: {
      name: '',
      isPublic: false,
      enabledProxy: false,
      registryId: undefined as number | undefined,
      organization: undefined as string | undefined,
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync(value)
      await onCreated?.(value.name)
      onClose()
    },
  })

  // Reset form when modal opens
  const resetForm = useEffectEvent(() => {
    if (opened) {
      form.reset()
      mutation.reset()
    }
  })

  useEffect(() => {
    resetForm()
  }, [opened])

  // Fetch registries for the dropdown when proxy is enabled
  const registriesQuery = useQuery({
    ...allRegistriesQueryOptions(),
    enabled: opened,
  })

  const registryOptions = (registriesQuery.data ?? []).map(r => ({
    value: String(r.id),
    label: r.name ?? r.url ?? '',
  }))

  // Both states are decided only once the request has succeeded, so loading
  // and error states show neither the guidance nor the shortcut, and the row
  // does not shift as the request settles.
  const showEmptyRegistryHint
    = registriesQuery.isSuccess && registryOptions.length === 0
  const showCreateRegistryShortcut
    = registriesQuery.isSuccess && registryOptions.length > 0

  // Opened in a new tab so the half-filled project form stays untouched here.
  const registryLink = <RegistryCreateLink />

  // Sits next to whichever registry action is currently shown: the create
  // shortcut when registries exist, the guidance line when none do.
  const refreshRegistriesButton = (
    <Tooltip label={t('projects.createModal.refreshRegistries')}>
      <ActionIcon
        variant="subtle"
        color="gray"
        aria-label={t('projects.createModal.refreshRegistries')}
        loading={registriesQuery.isFetching}
        onClick={() => void registriesQuery.refetch()}
      >
        <IconRefresh size={16} />
      </ActionIcon>
    </Tooltip>
  )

  const handleSubmit = () => {
    void form.handleSubmit()
  }

  return (
    <ModalWrapper
