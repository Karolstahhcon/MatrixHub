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

  const memberTypeOptions = [
    {
      value: MemberType.MEMBER_TYPE_USER,
      label: t('projects.detail.membersPage.memberType.user'),
    },
    {
      value: MemberType.MEMBER_TYPE_GROUP,
      disabled: true, // Group member type is not supported yet
      label: t('projects.detail.membersPage.memberType.group'),
    },
  ]

  const memberType = useStore(form.store, s => s.values.memberType)

  const roleOptions = useProjectRoleOptions()
  const roleDescription = useProjectRoleDescription()

  const { data: usersData } = useQuery({
    queryKey: ['users', 'list'],
    queryFn: () => Users.ListUsers({
      page: 1,
      pageSize: -1,
    }),
    enabled: opened && memberType === MemberType.MEMBER_TYPE_USER,
  })

  const userOptions = (usersData?.users ?? []).map(u => ({
    value: String(u.id ?? ''),
    label: u.username ?? '',
  }))

  const handleClose = () => {
    form.reset()
    onClose()
  }

  const canSubmit = useStore(form.store, s => s.canSubmit && !s.isSubmitting)
  const isSubmitting = useStore(form.store, s => s.isSubmitting)

  const handleConfirm = () => {
    void form.handleSubmit()
  }

  return (
    <ModalWrapper
      opened={opened}
      onClose={handleClose}
      title={t('projects.detail.membersPage.addMemberModal.title')}
      confirmLoading={isSubmitting}
      onConfirm={canSubmit ? handleConfirm : undefined}
    >
      <Stack gap="md">
        <form.Field
          name="memberType"
