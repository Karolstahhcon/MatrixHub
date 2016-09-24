import { Category } from '@matrixhub/api-ts/v1alpha1/model.pb'
import {
  describe, expect, it,
} from 'vitest'

import {
  buildDockerSnippet,
  buildEnvSnippet,
  buildLowLevelSnippet,
  buildPipelineSnippet,
  buildTestRequestSnippet,
  resolveUseModelTask,
  type SnippetPrompts,
} from './useModelGuides'

const prompts: SnippetPrompts = {
  textPrompt: 'Who are you?',
  imageUrl: '<image-url>',
  imagePrompt: 'Describe this image',
}

describe('resolveUseModelTask', () => {
  it('detects image-text-to-text from TASK labels regardless of case', () => {
    expect(resolveUseModelTask({
      labels: [{
        category: Category.TASK,
        name: 'Image-Text-to-Text',
      }],
    })).toBe('image-text-to-text')
  })

  it('detects text generation from TASK labels', () => {
    expect(resolveUseModelTask({
      labels: [{
        category: Category.TASK,
        name: 'text-generation',
      }],
    })).toBe('text-generation')
  })

  it('does not guess a task for unsupported or missing TASK labels', () => {
    expect(resolveUseModelTask({
      labels: [{
        category: Category.LIBRARY,
        name: 'image-text-to-text',
      }],
    })).toBeNull()
    expect(resolveUseModelTask({})).toBeNull()
  })
})

describe('snippet builders', () => {
  it('uses the configured HF endpoint', () => {
    expect(buildEnvSnippet('https://hub.example.com').code)
      .toContain('export HF_ENDPOINT="https://hub.example.com"')
  })
