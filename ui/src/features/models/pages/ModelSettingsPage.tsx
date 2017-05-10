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

