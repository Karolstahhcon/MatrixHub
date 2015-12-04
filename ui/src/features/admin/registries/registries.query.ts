import { Registries } from '@matrixhub/api-ts/v1alpha1/registry.pb'
import {
  keepPreviousData,
  queryOptions,
  useQuery,
} from '@tanstack/react-query'

import {
  DEFAULT_REGISTRIES_PAGE_SIZE,
  type RegistriesSearch,
} from './registries.schema'

export const adminRegistryKeys = {
