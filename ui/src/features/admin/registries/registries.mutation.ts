import {
  type CreateRegistryRequest,
  type PingRegistryRequest,
  type DeleteRegistryRequest,
  Registries,
  type UpdateRegistryRequest,
} from '@matrixhub/api-ts/v1alpha1/registry.pb'
import { mutationOptions } from '@tanstack/react-query'

import i18n from '@/i18n'

import { adminRegistryKeys } from './registries.query'

import type { NotificationMeta } from '@/types/tanstack-query'

function requireRegistryId(id?: number) {
  if (id == null) {
    throw new Error(i18n.t('routes.admin.registries.errors.missingRegistryId'))
  }

  return id
}

