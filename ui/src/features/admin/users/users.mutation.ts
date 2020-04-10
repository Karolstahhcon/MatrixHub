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
