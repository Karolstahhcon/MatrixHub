import {
  describe, expect, it,
} from 'vitest'

import { buildProxyDownloadCommand, isProxyProject } from './projects.utils'

describe('proxy projects', () => {
  it('recognizes a configured registry and builds its upstream download command', () => {
    expect(isProxyProject(' https://huggingface.co ')).toBe(true)
