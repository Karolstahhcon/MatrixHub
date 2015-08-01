import { Text } from '@mantine/core'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { ModalWrapper } from '@/shared/components/ModalWrapper'

import { deleteProjectMutationOptions } from '../projects.mutation'

import type { Project } from '@matrixhub/api-ts/v1alpha1/project.pb'

export interface DeleteProjectModalProps {
  project: Project | null
  opened: boolean
  onClose: () => void
}

export function DeleteProjectModal({
  project,
  opened,
  onClose,
}: DeleteProjectModalProps) {
  const { t } = useTranslation()
  const mutation = useMutation(deleteProjectMutationOptions())

  const handleConfirm = () => {
