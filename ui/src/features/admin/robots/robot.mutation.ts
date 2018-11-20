import {
  type CreateRobotAccountRequest,
  type DeleteRobotAccountRequest, type GetRobotAccountResponse,
  type RefreshRobotAccountTokenRequest,
  RobotAccountStatus,
  Robots,
  type UpdateRobotAccountRequest,
} from '@matrixhub/api-ts/v1alpha1/robot.pb'
import { mutationOptions } from '@tanstack/react-query'

import i18n from '@/i18n'

import { adminRobotKeys } from './robot.query'

import type { NotificationMeta } from '@/types/tanstack-query'

function requireRobotId(id?: number) {
  if (id == null) {
    throw new Error(i18n.t('routes.admin.robots.errors.missingRobotId'))
  }

  return id
}

function robotStatusMutationOptions({
  status,
  successMessage,
  errorMessage,
}: {
  status: RobotAccountStatus
  successMessage: string
  errorMessage: string
}) {
  return mutationOptions({
    mutationFn: (robot: GetRobotAccountResponse) => Robots.UpdateRobotAccount({
      ...robot,
      status,
    }),
