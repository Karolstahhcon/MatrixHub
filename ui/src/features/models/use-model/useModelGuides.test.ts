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
