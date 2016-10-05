import {
  Checkbox,
  Group,
  Input,
  MultiSelect,
  NumberInput,
  Radio,
  RadioGroup,
  Select,
  Stack,
  Text,
  TextInput,
  Tooltip,
  Textarea,
} from '@mantine/core'
import {
  ResourceType,
  SyncPolicyType,
  TriggerType,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'
import { useStore } from '@tanstack/react-form'
import { useMutation, useQuery } from '@tanstack/react-query'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { registriesQueryOptions } from '@/features/admin/registries/registries.query'
import { FieldHintLabel } from '@/shared/components/FieldHintLabel'
import { ModalWrapper } from '@/shared/components/ModalWrapper'
import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form'

import {
  createReplicationMutationOptions,
  updateReplicationMutationOptions,
} from '../replications.mutation'
import {
  createReplicationFormSchema,
  replicationCronExpressionSchema,
  type ReplicationResourceTypeValue,
  type ReplicationFormValues,
  replicationResourceTypeValues,
  replicationTriggerTypeValues,
  DESCRIPTION_MAX_LENGTH,
} from '../replications.schema'
import {
  convertBandwidthToApiValue,
  getCronExpressionDescription,
  getDefaultBandwidthUnit,
  getDefaultBandwidthValue,
  normalizeFiveFieldCronExpression,
} from '../replications.utils'
import classes from './ReplicationFormModal.module.css'

import type {
  CreateSyncPolicyRequest,
  SyncPolicyItem,
  UpdateSyncPolicyRequest,
} from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'

const BANDWIDTH_MIN = -1
const BANDWIDTH_MAX = 1048576
// The inline label column width is locale-dependent and lives in the CSS module
// (`--inline-field-label-width`), so all three fields share one aligned column.
const inlineFieldProps = {
  variant: 'unstyled' as const,
  leftSectionProps: {
    className: classes.inlineFieldSection,
  },
  classNames: {
    wrapper: classes.inlineFieldWrapper,
    input: classes.inlineFieldInput,
  },
} as const

function renderInlineFieldSection(label: string) {
  return (
    <Tooltip
      label={label}
      openDelay={200}
    >
      <Text
        span
        inherit
        truncate="end"
        className={classes.inlineFieldSectionLabel}
      >
        {label}
      </Text>
    </Tooltip>
  )
}

export interface ReplicationFormModalProps {
  mode: 'create' | 'edit'
  opened: boolean
  syncPolicy?: SyncPolicyItem
  onClose: () => void
}

// The backend accepts multiple resource types. Legacy ALL/empty values are
// normalized to both supported concrete types when editing existing data.
function normalizeResourceTypesForForm(
  resourceTypes?: readonly ResourceType[],
): ReplicationResourceTypeValue[] {
  const selected = (resourceTypes ?? []).filter(
    (resourceType): resourceType is ReplicationResourceTypeValue =>
      replicationResourceTypeValues.includes(resourceType as ReplicationResourceTypeValue),
  )

  return selected.length > 0
    ? selected
    : [...replicationResourceTypeValues]
}

function buildResourceTypesForApi(
  resourceTypes: ReplicationResourceTypeValue[],
): ResourceType[] {
  return resourceTypes
}

function buildPolicyPayload(values: ReplicationFormValues) {
  const bandwidth = convertBandwidthToApiValue(values.bandwidth ?? '', values.bandwidthUnit)
  const resourceTypes = buildResourceTypesForApi(values.resourceTypes)
  const base = {
    description: values.description,
    triggerType: values.triggerType,
    bandwidth,
    isOverwrite: values.isOverwrite,
    ...(values.triggerType === TriggerType.TRIGGER_TYPE_SCHEDULED
      ? {
          triggerTypeSchedule: {
            cron: normalizeFiveFieldCronExpression(values.cronExpression),
          },
        }
      : {}),
  }

  if (values.policyType === SyncPolicyType.SYNC_POLICY_TYPE_PUSH_BASE) {
    return {
      ...base,
      pushBasePolicy: {
        resourceName: values.resourceName ?? '',
        resourceTypes,
        targetRegistryId: values.targetRegistryId,
        targetProjectName: values.targetProjectName ?? '',
      },
    }
  }

  return {
    ...base,
    pullBasePolicy: {
      sourceRegistryId: values.sourceRegistryId,
      resourceName: values.resourceName ?? '',
      resourceTypes,
      targetProjectName: values.targetProjectName ?? '',
    },
  }
}

function getFormDefaults(syncPolicy?: SyncPolicyItem): ReplicationFormValues {
  const isPush = syncPolicy?.policyType === SyncPolicyType.SYNC_POLICY_TYPE_PUSH_BASE
  const rawResourceTypes = isPush
    ? syncPolicy?.pushBasePolicy?.resourceTypes
    : syncPolicy?.pullBasePolicy?.resourceTypes

  return {
    name: syncPolicy?.name ?? '',
    description: syncPolicy?.description ?? '',
    policyType: isPush
      ? SyncPolicyType.SYNC_POLICY_TYPE_PUSH_BASE
      : SyncPolicyType.SYNC_POLICY_TYPE_PULL_BASE,
    triggerType: syncPolicy?.triggerType === TriggerType.TRIGGER_TYPE_SCHEDULED
      ? TriggerType.TRIGGER_TYPE_SCHEDULED
      : TriggerType.TRIGGER_TYPE_MANUAL,
    cronExpression: syncPolicy?.triggerTypeSchedule?.cron ?? '',
    bandwidth: getDefaultBandwidthValue(syncPolicy?.bandwidth),
    bandwidthUnit: getDefaultBandwidthUnit(syncPolicy?.bandwidth),
    isOverwrite: syncPolicy?.isOverwrite ?? true,
    resourceName: isPush
      ? (syncPolicy?.pushBasePolicy?.resourceName ?? '')
      : (syncPolicy?.pullBasePolicy?.resourceName ?? ''),
    resourceTypes: normalizeResourceTypesForForm(rawResourceTypes),
    sourceRegistryId: syncPolicy?.pullBasePolicy?.sourceRegistryId ?? 0,
    targetProjectName: isPush
      ? (syncPolicy?.pushBasePolicy?.targetProjectName ?? '')
      : (syncPolicy?.pullBasePolicy?.targetProjectName ?? ''),
    targetRegistryId: syncPolicy?.pushBasePolicy?.targetRegistryId ?? 0,
  }
}

export function ReplicationFormModal({
  mode,
  opened,
  syncPolicy,
  onClose,
}: ReplicationFormModalProps) {
  const {
    t,
    i18n,
  } = useTranslation()
  const createMutation = useMutation(createReplicationMutationOptions())
  const updateMutation = useMutation(updateReplicationMutationOptions())
  const registriesQuery = useQuery({
    ...registriesQueryOptions({
      page: 1,
      query: '',
    }),
    enabled: opened,
  })

  const registryOptions = useMemo(() => {
    const options = new Map<string, {
      value: string
      label: string
    }>()

    for (const registry of registriesQuery.data?.registries ?? []) {
      if (registry.id == null || registry.id < 1) {
        continue
      }

      const value = String(registry.id)

      if (!value || options.has(value)) {
        continue
      }

      options.set(value, {
        value,
        label: registry.name || registry.url || `#${registry.id}`,
      })
    }

    return Array.from(options.values())
  }, [registriesQuery.data])

  const resourceTypeOptions = useMemo(() => replicationResourceTypeValues.map(
    value => ({
      value,
      label: value === ResourceType.RESOURCE_TYPE_MODEL
        ? t('routes.admin.replications.form.resourceTypes.model')
        : t('routes.admin.replications.form.resourceTypes.dataset'),
    }),
  ), [t])

  const triggerTypeOptions = useMemo(() => replicationTriggerTypeValues.map(
    value => ({
      value,
      label: value === TriggerType.TRIGGER_TYPE_SCHEDULED
        ? t('routes.admin.replications.trigger.scheduled')
        : t('routes.admin.replications.trigger.manual'),
