import {
  describe, expect, it,
} from 'vitest'

import {
  buildDownloadSnippet,
  buildHfEnvSnippet,
  buildLocalPathPlaceholder,
  buildUploadSnippet,
} from './commandGuides'

describe('command guide snippets', () => {
  it('only asks for a token when uploading', () => {
    expect(buildHfEnvSnippet('download', 'https://hub.example.com').code)
      .toBe('export HF_ENDPOINT="https://hub.example.com"')
    expect(buildHfEnvSnippet('upload', 'https://hub.example.com').code)
      .toContain('export HF_TOKEN=')
