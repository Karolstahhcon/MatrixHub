import {
  Select,
  Stack,
  TextInput,
} from '@mantine/core'
import { useStore } from '@tanstack/react-form'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

import { ModalWrapper } from '@/shared/components/ModalWrapper'
import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form'

import { useProjectRoleDescription, useProjectRoleOptions } from '../member.utils'
import { updateMemberRoleMutationOptions } from '../members.mutation'

import type { ProjectMember } from '@matrixhub/api-ts/v1alpha1/project.pb'
import type { ProjectRoleType } from '@matrixhub/api-ts/v1alpha1/role.pb'

interface EditRoleModalProps {
  opened: boolean
  onClose: () => void
  projectId: string
  member: ProjectMember | null
}

export function EditRoleModal({
  opened,
  onClose,
  projectId,
  member,
}: EditRoleModalProps) {
  const { t } = useTranslation()
  const mutation = useMutation(updateMemberRoleMutationOptions())

