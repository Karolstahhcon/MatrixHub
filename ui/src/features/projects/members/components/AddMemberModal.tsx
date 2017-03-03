import { Select, Stack } from '@mantine/core'
import { MemberType } from '@matrixhub/api-ts/v1alpha1/project.pb'
import { ProjectRoleType } from '@matrixhub/api-ts/v1alpha1/role.pb'
import { Users } from '@matrixhub/api-ts/v1alpha1/user.pb'
import { useStore } from '@tanstack/react-form'
import {
  useMutation,
  useQuery,
} from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

import { ModalWrapper } from '@/shared/components/ModalWrapper'
import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form'

import { useProjectRoleDescription, useProjectRoleOptions } from '../member.utils'
import { addMemberMutationOptions } from '../members.mutation'

interface AddMemberModalProps {
  opened: boolean
  onClose: () => void
  projectId: string
}

const defaultValues = {
  memberType: MemberType.MEMBER_TYPE_USER as string,
  memberId: '' as string,
  role: ProjectRoleType.ROLE_TYPE_PROJECT_VIEWER as string,
}

export function AddMemberModal({
  opened,
  onClose,
  projectId,
}: AddMemberModalProps) {
  const { t } = useTranslation()
  const mutation = useMutation(addMemberMutationOptions())

  const requiredString = z.string().min(1, t('common.validation.required'))

  const form = useForm({
    defaultValues,
    onSubmit: async ({ value }) => {
      await mutation.mutateAsync({
        name: projectId,
        memberType: value.memberType as MemberType,
        memberId: Number(value.memberId),
        role: value.role as ProjectRoleType,
      })
      handleClose()
    },
  })

