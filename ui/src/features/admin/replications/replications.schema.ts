import {
  ResourceType,
  SyncPolicyType,
  TriggerType,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { z } from 'zod'

import i18n from '@/i18n'
import { withNameStartCharRule, withNameValidCharsRule } from '@/shared/validation'

import { isFiveFieldCronExpression } from './replications.utils'

export const DEFAULT_REPLICATIONS_PAGE = 1
export const DEFAULT_REPLICATIONS_PAGE_SIZE = 10
export const DESCRIPTION_MAX_LENGTH = 50

export const replicationResourceTypeValues = [
  ResourceType.RESOURCE_TYPE_MODEL,
  ResourceType.RESOURCE_TYPE_DATASET,
] as const
export type ReplicationResourceTypeValue = typeof replicationResourceTypeValues[number]

export const replicationTriggerTypeValues = [
  TriggerType.TRIGGER_TYPE_MANUAL,
  TriggerType.TRIGGER_TYPE_SCHEDULED,
] as const

export const replicationBandwidthUnitValues = ['Kbps', 'Mbps'] as const
export type ReplicationBandwidthUnit = typeof replicationBandwidthUnitValues[number]

export const replicationsSearchDefaults = {
  page: DEFAULT_REPLICATIONS_PAGE,
  query: undefined as string | undefined,
}

export const replicationsSearchSchema = z.object({
  page: z.coerce.number().int().positive().default(replicationsSearchDefaults.page).catch(replicationsSearchDefaults.page),
  query: z.string().trim().optional().catch(replicationsSearchDefaults.query),
})

export type ReplicationsSearch = z.infer<typeof replicationsSearchSchema>

function t(key: string) {
  return i18n.getFixedT(i18n.resolvedLanguage ?? i18n.language)(key)
}

export function replicationCronExpressionSchema() {
  return z
    .string()
    .trim()
    .min(1, t('routes.admin.replications.validation.cronRequired'))
    .refine(
      value => isFiveFieldCronExpression(value),
