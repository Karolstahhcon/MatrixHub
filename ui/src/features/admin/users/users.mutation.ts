import {
  type CreateUserRequest,
  type DeleteUserRequest,
  type ResetUserPasswordRequest,
  type SetUserSysAdminRequest,
  type User,
  Users,
} from '@matrixhub/api-ts/v1alpha1/user.pb'
import { mutationOptions } from '@tanstack/react-query'

import i18n from '@/i18n'

import { adminUserKeys } from './users.query'

import type { NotificationMeta } from '@/types/tanstack-query'

function requireUserId(id?: number) {
  if (id == null) {
    throw new Error(i18n.t('routes.admin.users.errors.missingUserId'))
  }

  return id
}

type UserSysAdminMutationInput = Pick<SetUserSysAdminRequest, 'id'>

function userSysAdminMutationOptions({
  sysadminFlag,
  successMessage,
  errorMessage,
}: {
  sysadminFlag: boolean
  successMessage: string
  errorMessage: string
}) {
  return mutationOptions({
    mutationFn: ({ id }: UserSysAdminMutationInput) =>
      Users.SetUserSysAdmin({
        id: requireUserId(id),
        sysadminFlag,
      }),
    meta: {
      successMessage,
      errorMessage,
      invalidates: [adminUserKeys.lists()],
    } satisfies NotificationMeta,
  })
}

export function createUserMutationOptions() {
  return mutationOptions({
    mutationFn: (input: CreateUserRequest) => Users.CreateUser(input),
    meta: {
      successMessage: i18n.t('routes.admin.users.notifications.createSuccess'),
      errorMessage: i18n.t('routes.admin.users.notifications.createError'),
      invalidates: [adminUserKeys.lists()],
    } satisfies NotificationMeta,
  })
}

export function resetUserPasswordMutationOptions() {
