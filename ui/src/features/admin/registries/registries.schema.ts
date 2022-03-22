import { z } from 'zod'

export const DEFAULT_REGISTRIES_PAGE = 1
export const DEFAULT_REGISTRIES_PAGE_SIZE = 10

export const registriesSearchDefaults = {
  page: DEFAULT_REGISTRIES_PAGE,
  query: undefined as string | undefined,
  create: undefined as boolean | undefined,
}

export const registriesSearchSchema = z.object({
