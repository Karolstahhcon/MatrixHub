import {
  RobotAccountProjectScope,
  type GetRobotAccountResponse,
} from '@matrixhub/api-ts/v1alpha1/robot.pb'
import { z } from 'zod'

import i18n from '@/i18n'

import type { TFunction } from 'i18next'

export const ROBOT_DESCRIPTION_MAX_LENGTH = 50
export const DEFAULT_ROBOT_EXPIRE_DAYS = 30

export const robotProjectScopeSchema = z.enum([
  RobotAccountProjectScope.ROBOT_ACCOUNT_PROJECT_SCOPE_ALL,
  RobotAccountProjectScope.ROBOT_ACCOUNT_PROJECT_SCOPE_SELECTED,
])

export const robotExpiryModeSchema = z.enum(['never', 'days'])

export const createRobotAccountFormSchema = (t: TFunction) => z.object({
  name: z.string().trim().min(1, t('common.validation.fieldRequired', { field: t('routes.admin.robots.fields.name') })),
  description: z.string().max(ROBOT_DESCRIPTION_MAX_LENGTH, t('common.validation.maxLength', {
    max: ROBOT_DESCRIPTION_MAX_LENGTH,
    field: t('routes.admin.robots.fields.description'),
  })),
  expiryMode: robotExpiryModeSchema,
  expireDays: z.number().int().min(1, t('common.validation.minNumber', {
    min: 1,
    field: t('routes.admin.robots.fields.expiry'),
  })).optional(),
  platformPermissions: z.array(z.string()),
  projectPermissions: z.array(z.string()),
  projectScope: robotProjectScopeSchema,
  projects: z.array(z.string()),
}).superRefine((value, ctx) => {
  if (value.expiryMode === 'days' && value.expireDays == null) {
    ctx.addIssue({
      code: 'custom',
      message: t('common.validation.fieldRequired', {
        field: t('routes.admin.robots.fields.expiry'),
      }),
      path: ['expireDays'],
    })
  }
})

export const ROBOT_TOKEN_MIN_LENGTH = 8
export const ROBOT_TOKEN_MAX_LENGTH = 20
export const ROBOT_TOKEN_COMPLEXITY_PATTERN = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).*$/

export const refreshRobotTokenFormDefaults = {
  autoGenerate: true,
  token: '',
  confirmToken: '',
