import {
  Button, Checkbox, Group, InputWrapper, rem, Stack, TextInput,
} from '@mantine/core'
import { ProjectType } from '@matrixhub/api-ts/v1alpha1/project.pb'
import { useMutation, useSuspenseQuery } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

import { useForm } from '@/shared/hooks/useForm'

import { ProjectTypeHintLabel } from '../components/ProjectTypeHintLabel'
import { updateProjectMutationOptions } from '../projects.mutation'
import { projectDetailQueryOptions } from '../projects.query'

const projectRouteApi = getRouteApi('/(auth)/(app)/projects/$projectId')

export function ProjectSettingsPage() {
  const { t } = useTranslation()
  const { projectId } = projectRouteApi.useParams()
  const { data: project } = useSuspenseQuery(projectDetailQueryOptions(projectId))
  const mutation = useMutation(updateProjectMutationOptions())

  const isProxy = !!project.registryUrl

  const projectSettingForm = useForm({
    defaultValues: {
      type: project.type ?? ProjectType.PROJECT_TYPE_PRIVATE,
