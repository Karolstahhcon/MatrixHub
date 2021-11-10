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
