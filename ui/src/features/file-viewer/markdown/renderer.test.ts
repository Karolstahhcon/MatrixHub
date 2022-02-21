import {
  describe, expect, it,
} from 'vitest'

import {
  renderMarkdown, slugifyHeading, stripFrontmatter,
} from './renderer'

describe('slugifyHeading', () => {
  it('matches GitHub slugs', () => {
    expect(slugifyHeading('Best Practices')).toBe('best-practices')
    expect(slugifyHeading('Qwen3 Highlights')).toBe('qwen3-highlights')
    expect(slugifyHeading('What\'s new? (v2.0)')).toBe('whats-new-v20')
    expect(slugifyHeading('模型 概览')).toBe('模型-概览')
  })
})

describe('stripFrontmatter', () => {
  it('strips a leading YAML mapping block', () => {
    const input = '---\nlicense: apache-2.0\nbase_model:\n- Qwen/Qwen3-0.6B-Base\n---\n# Title\n'

    expect(stripFrontmatter(input)).toBe('# Title\n')
  })

  it('keeps a document that opens with a thematic break', () => {
    const input = '---\n# Introduction\nSome text.\n---\n# Details\n'

    expect(stripFrontmatter(input)).toBe(input)
  })

  it('leaves content without frontmatter untouched', () => {
    expect(stripFrontmatter('# Title\n')).toBe('# Title\n')
  })
