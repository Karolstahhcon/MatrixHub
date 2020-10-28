import { Button } from '@mantine/core'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import {
  revokeUserAdminMutationOptions,
  setUserAdminMutationOptions,
} from '../users.mutation'

import type { User } from '@matrixhub/api-ts/v1alpha1/user.pb'

interface ToggleUserAdminActionProps {
  user: User
  disabled?: boolean
}

export function ToggleUserAdminAction({
  user,
  disabled,
}: ToggleUserAdminActionProps) {
  const { t } = useTranslation()
  const isAdmin = !!user.isAdmin
  const mutation = useMutation(
    isAdmin
      ? revokeUserAdminMutationOptions()
      : setUserAdminMutationOptions(),
  )

  const handleClick = async () => {
