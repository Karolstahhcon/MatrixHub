import {
  Alert, Button, Checkbox, rem, Stack, Text, TextInput,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { IconAlertTriangle, IconTrash } from '@tabler/icons-react'
import { useMutation, useSuspenseQuery } from '@tanstack/react-query'
import { getRouteApi, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { SetPopularModelModal } from '@/features/models/components/SetPopularModelModal.tsx'
import { deleteModelMutationOptions } from '@/features/models/models.mutation'
import { modelQueryOptions } from '@/features/models/models.query.ts'

const { useParams } = getRouteApi('/(auth)/(app)/projects_/$projectId/models/$modelId/settings/')

export function ModelSettingsPage() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const {
    projectId, modelId,
  } = useParams()

  const [inputValue, setInputValue] = useState('')
  const { data: model } = useSuspenseQuery(modelQueryOptions(projectId, modelId))

  const [popularModalOpened, popularModalHandlers] = useDisclosure(false)
  const deleteMutation = useMutation(deleteModelMutationOptions())

  const fullName = `${projectId}/${modelId}`
  const isPopular = Boolean(model.popular)

  const handleDelete = async () => {
    await deleteMutation.mutateAsync({
      project: projectId,
      name: modelId,
    })

    // Fixme: confirm navigation to which page
    void navigate({
      to: '/projects/$projectId/models',
      params: {
        projectId,
      },
    })
  }

  return (
    <Stack
      pt={20}
      gap={28}
      align="flex-start"
    >
      <Stack
        gap="md"
        align="flex-start"
      >
        <Text fw="600" lh={rem(24)}>
          {t('model.settings.recommended.title')}
        </Text>

        <Checkbox
          checked={model.popular}
          label={t('model.settings.recommended.set.label')}
          onClick={popularModalHandlers.open}
        />

        <SetPopularModelModal
          opened={popularModalOpened}
          projectId={projectId}
          modelId={modelId}
          isPopular={isPopular}
          onClose={popularModalHandlers.close}
        />
      </Stack>

      <Stack
        gap="md"
        align="flex-start"
      >
        <Text fw="600" lh={rem(24)}>
          {t('model.settings.delete.title')}
        </Text>

        <Alert
          icon={<IconAlertTriangle size={20} />}
          variant="light"
          color="var(--mantine-color-yellow-6)"
        >
          <Text size="sm" lh={rem(20)} c="var(--mantine-color-gray-9)">
            {t('model.settings.delete.warning', { name: fullName })}
          </Text>
        </Alert>

        <Stack gap={8}>
