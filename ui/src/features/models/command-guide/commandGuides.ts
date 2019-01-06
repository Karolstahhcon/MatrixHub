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

