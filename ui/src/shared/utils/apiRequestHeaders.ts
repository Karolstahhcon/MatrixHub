import {
  setFetchFn,
  type FetchFn,
} from '@matrixhub/api-ts/fetch.pb'

import i18n, {
  DEFAULT_LANGUAGE,
  normalizeLanguage,
} from '@/i18n'

const API_PATH_PREFIXES = ['/api', '/apis']

type FetchInput = Parameters<typeof window.fetch>[0]

type FetchInit = Parameters<typeof window.fetch>[1]

type MatrixHubWindow = Window & typeof globalThis & {
  __matrixhubApiRequestHeadersInstalled?: boolean
}

export type ApiRequestHeaderProvider = () => HeadersInit | null | undefined

// Keep cross-cutting SDK headers in one provider list so future headers can be
// added here without wrapping every generated SDK call.
const apiRequestHeaderProviders: ApiRequestHeaderProvider[] = [
