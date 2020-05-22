import {
  Alert,
  Box,
  Button,
  Group,
  Paper,
  rem,
  Stack,
  Text,
  TextInput,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import {
  IconFolderPlus,
  IconInfoCircle,
  IconPlus,
} from '@tabler/icons-react'
import { useMutation } from '@tanstack/react-query'
import { useNavigate, useRouter } from '@tanstack/react-router'
import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'

import { createModelMutationOptions } from '@/features/models/models.mutation'
import { useWritableModelProjects } from '@/features/models/models.query.ts'
import { createModelSchema } from '@/features/models/models.schema'
import { CreateProjectModal } from '@/features/projects/components/CreateProjectModal'
import { ProjectSelect } from '@/shared/components/ProjectSelect'
import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form.ts'

interface ModelCreatePageProps {
  initialProjectId?: string
}

export function ModelCreatePage({ initialProjectId = '' }: ModelCreatePageProps) {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const router = useRouter()
  const [createProjectOpened, createProjectHandlers] = useDisclosure(false)

  const mutation = useMutation(createModelMutationOptions())
  const modelCreateSchema = createModelSchema(t)
  const {
    name: nameValidator, projectId: projectIdValidator,
  } = modelCreateSchema.shape

  const handleNavigateBack = () => {
    if (initialProjectId) {
      return navigate({
        to: '/projects/$projectId/models',
        params: {
          projectId: initialProjectId,
        },
      })
    }

    if (router.history.length > 1) {
      return router.history.back()
    }

    return navigate({ to: '/models' })
  }

  const form = useForm({
    defaultValues: {
      name: '',
      projectId: initialProjectId?.trim(),
    },
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        name: value.name,
        project: value.projectId,
      })

      handleNavigateBack()
    },
  })

  const {
    data: projects = [],
    isLoadingError: isProjectsLoadingError,
    isPending: isProjectsPending,
    refetch: refetchProjects,
  } = useWritableModelProjects()
  const showEmptyProjectPrompt = !initialProjectId
    && !isProjectsPending
    && !isProjectsLoadingError
