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
  all: ['admin', 'registries'] as const,
  lists: () => [...adminRegistryKeys.all, 'list'] as const,
  // Only paging/search affect the response, so UI-only params (e.g. `create`)
  // must not become part of the cache key.
  list: (search: RegistriesSearch) => [
    ...adminRegistryKeys.lists(),
    {
      page: search.page,
      query: search.query,
    },
  ] as const,
  allList: () => [...adminRegistryKeys.lists(), 'all'] as const,
}

