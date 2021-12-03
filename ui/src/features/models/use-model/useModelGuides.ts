import { Category, type Model } from '@matrixhub/api-ts/v1alpha1/model.pb'

import { getLabelsByCategory } from '@/features/models/models.utils'

export const USE_MODEL_ENGINES = ['transformers', 'vllm', 'sglang'] as const

export type UseModelEngine = (typeof USE_MODEL_ENGINES)[number]

/** Model capability the generated Transformers snippets are specialised for. */
export type UseModelTask = 'text-generation' | 'image-text-to-text'

export const ENGINE_LABELS: Record<UseModelEngine, string> = {
  transformers: 'Transformers',
  vllm: 'vLLM',
  sglang: 'SGLang',
}

export const HF_TOKEN_PLACEHOLDER = '<your-matrixhub-token>'

export const ENGINE_DOC_URLS: Record<UseModelEngine, string> = {
  transformers: 'https://huggingface.co/docs/transformers/index',
  vllm: 'https://docs.vllm.ai/en/latest/',
  sglang: 'https://docs.sglang.ai/',
}

export interface UseModelSnippet {
  code: string
  lang: 'bash' | 'python'
}

/**
 * Pick the task the snippets should target from the model's TASK labels.
 * Only text generation and image-text-to-text are supported; anything else
 * requires the user to choose a supported task explicitly.
 */
export function resolveUseModelTask(model: Pick<Model, 'labels'>): UseModelTask | null {
  const tasks = getLabelsByCategory(model.labels, Category.TASK).map(task => task.toLowerCase())

  if (tasks.some(task => task === 'image-text-to-text' || task === 'image-to-text')) {
    return 'image-text-to-text'
  }

  if (tasks.includes('text-generation')) {
    return 'text-generation'
  }

  return null
}

export function buildEnvSnippet(hfEndpoint: string): UseModelSnippet {
  return {
    lang: 'bash',
    code: `export HF_ENDPOINT="${hfEndpoint}"\nexport HF_TOKEN="${HF_TOKEN_PLACEHOLDER}"`,
  }
}

export function buildInstallSnippet(engine: UseModelEngine): UseModelSnippet {
  const commands: Record<UseModelEngine, string> = {
    transformers: 'pip install -U transformers torch accelerate',
    vllm: 'pip install vllm',
    sglang: 'pip install "sglang[all]"',
  }

  return {
    lang: 'bash',
    code: commands[engine],
  }
}

export function buildServeSnippet(engine: Exclude<UseModelEngine, 'transformers'>, modelPath: string): UseModelSnippet {
  const commands = {
    vllm: `vllm serve "${modelPath}" --port 8000`,
    sglang: `python -m sglang.launch_server \\\n    --model-path "${modelPath}" \\\n    --host 0.0.0.0 --port 30000`,
  }

  return {
    lang: 'bash',
    code: commands[engine],
  }
}

