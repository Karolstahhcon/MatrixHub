import type { GuideSnippet } from '@/features/models/guides/GuideDrawer'

export type ModelCommandType = 'upload' | 'download'

export const HF_TOKEN_PLACEHOLDER = '<your-matrixhub-token>'

/** Default cache dir used by huggingface_hub when no --local-dir is given. */
export const HF_DEFAULT_CACHE_DIR = '~/.cache/huggingface/hub'

export function buildInstallHfSnippet(): GuideSnippet {
  return {
    lang: 'bash',
    code: 'pip install -U huggingface_hub',
  }
}

/** Download only needs the endpoint; upload also needs a token to write. */
export function buildHfEnvSnippet(type: ModelCommandType, hfEndpoint: string): GuideSnippet {
  const lines = [`export HF_ENDPOINT="${hfEndpoint}"`]

  if (type === 'upload') {
    lines.push(`export HF_TOKEN="${HF_TOKEN_PLACEHOLDER}"`)
  }

  return {
    lang: 'bash',
    code: lines.join('\n'),
  }
}

export function buildDownloadSnippet(modelPath: string): GuideSnippet {
  return {
