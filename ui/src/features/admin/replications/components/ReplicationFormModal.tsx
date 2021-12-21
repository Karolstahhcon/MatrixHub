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
