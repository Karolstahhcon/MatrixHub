import { RegistryType } from '@matrixhub/api-ts/v1alpha1/registry.pb'
import { z } from 'zod'

import i18n from '@/i18n'

import type {
  CreateRegistryRequest,
  PingRegistryRequest,
  Registry,
  RegistryBasicCredential,
  UpdateRegistryRequest,
} from '@matrixhub/api-ts/v1alpha1/registry.pb'

export const REGISTRY_NAME_MAX_LENGTH = 64
export const REGISTRY_NAME_MIN_LENGTH = 2
export const REGISTRY_DESCRIPTION_MAX_LENGTH = 50
export const supportedRegistryTypes = [
  RegistryType.REGISTRY_TYPE_HUGGINGFACE,
  RegistryType.REGISTRY_TYPE_MATRIXHUB,
] as const
export const registryProviderLabelKeys: Partial<Record<RegistryType, string>> = {
  [RegistryType.REGISTRY_TYPE_HUGGINGFACE]: 'routes.admin.registries.provider.huggingFace',
  [RegistryType.REGISTRY_TYPE_MATRIXHUB]: 'routes.admin.registries.provider.matrixHub',
}

export interface RegistryUrlCopyKeys {
  hint: string
  description: string
  placeholder: string
  repositoryUrlWarning: string
}

const huggingFaceUrlCopyKeys: RegistryUrlCopyKeys = {
  hint: 'routes.admin.registries.form.urlHint.huggingFace',
  description: 'routes.admin.registries.form.urlDescription.huggingFace',
  placeholder: 'routes.admin.registries.form.urlPlaceholder.huggingFace',
  repositoryUrlWarning: 'routes.admin.registries.form.urlRepositoryWarning.huggingFace',
}

/* The target URL is a base address, so the guidance differs per provider:
Hugging Face accepts the upstream site or a mirror, MatrixHub expects an instance address. */
const registryUrlCopyKeys: Partial<Record<RegistryType, RegistryUrlCopyKeys>> = {
  [RegistryType.REGISTRY_TYPE_HUGGINGFACE]: huggingFaceUrlCopyKeys,
  [RegistryType.REGISTRY_TYPE_MATRIXHUB]: {
    hint: 'routes.admin.registries.form.urlHint.matrixHub',
    description: 'routes.admin.registries.form.urlDescription.matrixHub',
    placeholder: 'routes.admin.registries.form.urlPlaceholder.matrixHub',
    repositoryUrlWarning: 'routes.admin.registries.form.urlRepositoryWarning.matrixHub',
  },
}

export function getRegistryUrlCopyKeys(type: RegistryType): RegistryUrlCopyKeys {
  return registryUrlCopyKeys[type] ?? huggingFaceUrlCopyKeys
}

/* Hugging Face has a well-known upstream address that can be prefilled,
while a MatrixHub instance address is deployment specific and must be typed in. */
const registryDefaultUrls: Partial<Record<RegistryType, string>> = {
  [RegistryType.REGISTRY_TYPE_HUGGINGFACE]: 'https://huggingface.co',
  [RegistryType.REGISTRY_TYPE_MATRIXHUB]: '',
}

export function getRegistryDefaultUrl(type: RegistryType): string {
  return registryDefaultUrls[type] ?? ''
}

/* A URL still equal to its provider default (or empty) was never customized,
so switching providers may overwrite it without asking. */
export function isRegistryUrlPristine(url: string, type: RegistryType): boolean {
  const trimmedUrl = url.trim()

  return trimmedUrl.length === 0 || trimmedUrl === getRegistryDefaultUrl(type)
}

/* The stored URL is a base address that later gets joined with resource paths,
so trailing slashes are dropped to keep one canonical form. */
export function normalizeRegistryUrl(url: string): string {
  return url.trim().replace(/\/+$/u, '')
}

/* Paths that only ever appear on a concrete repository address, never on a
base address. Deliberately narrow: the goal is to catch the common
copy-the-model-page mistake, not to validate every possible URL shape. */
const registryRepositoryPathPrefixes = ['models', 'datasets', 'spaces']

/**
 * Detect a URL that looks like a specific repository rather than a base address,
 * e.g. `https://huggingface.co/Qwen/Qwen3-32B` or `https://hf-mirror.com/models/...`.
 * This is advisory only — the address may still be valid, so it never blocks submit.
 */
export function isLikelyRepositoryUrl(url: string): boolean {
  const normalizedUrl = normalizeRegistryUrl(url)

  if (!normalizedUrl) {
    return false
  }

  let segments: string[]

  try {
    segments = new URL(normalizedUrl).pathname
      .split('/')
      .filter(Boolean)
  } catch {
    return false
  }

  const [firstSegment] = segments

  if (firstSegment == null) {
    return false
  }

  if (registryRepositoryPathPrefixes.includes(firstSegment.toLowerCase())) {
    return true
  }

  return segments.length >= 2
}

export interface RegistryFormValues {
  type: RegistryType
  name: string
  description: string
  url: string
  username: string
  password: string
  verifyRemoteCert: boolean
}

export const defaultRegistryFormValues: RegistryFormValues = {
  type: RegistryType.REGISTRY_TYPE_HUGGINGFACE,
  name: '',
  description: '',
  url: '',
  username: '',
  password: '',
  verifyRemoteCert: true,
}

/* Name validation:
No spaces (including full-width/half-width) at the start, end, or within; auto-trimmed;
length 2-64 characters */
export function sanitizeRegistryName(value: string) {
