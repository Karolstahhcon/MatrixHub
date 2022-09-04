import { z } from 'zod'

export const DEFAULT_USERS_PAGE = 1
export const DEFAULT_USERS_PAGE_SIZE = 10

import i18n from '@/i18n'

// TODO: need confirm how to get error message in a better way
export function t(key: string) {
  return i18n.getFixedT(i18n.resolvedLanguage ?? i18n.language)(key)
}

export const usersSearchDefaults = {
  page: DEFAULT_USERS_PAGE,
  query: undefined as string | undefined,
}

export const usersSearchSchema = z.object({
