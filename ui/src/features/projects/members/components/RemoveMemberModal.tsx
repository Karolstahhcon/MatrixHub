import { Text } from '@mantine/core'
import { ProjectRoleType } from '@matrixhub/api-ts/v1alpha1/role.pb'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { ModalWrapper } from '@/shared/components/ModalWrapper'

import { removeMembersMutationOptions } from '../members.mutation'

import type { ProjectMember } from '@matrixhub/api-ts/v1alpha1/project.pb'

interface RemoveMemberModalProps {
  opened: boolean
  onClose: () => void
  projectId: string
  members: ProjectMember[]
}

function getRoleLabel(t: (key: string) => string, role?: ProjectRoleType) {
  switch (role) {
    case ProjectRoleType.ROLE_TYPE_PROJECT_ADMIN:
