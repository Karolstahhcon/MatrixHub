import type { HighlighterCore } from 'shiki/core'

/**
 * Grammar loaders keyed by language. Loaded lazily — only grammars actually
 * used are fetched (cpp alone is ~800 kB).
 */
const LANG_LOADERS: Record<string, () => Promise<unknown>> = {
  bash: () => import('@shikijs/langs/bash'),
  shell: () => import('@shikijs/langs/shell'),
  python: () => import('@shikijs/langs/python'),
  javascript: () => import('@shikijs/langs/javascript'),
  typescript: () => import('@shikijs/langs/typescript'),
  json: () => import('@shikijs/langs/json'),
  yaml: () => import('@shikijs/langs/yaml'),
  markdown: () => import('@shikijs/langs/markdown'),
  go: () => import('@shikijs/langs/go'),
  rust: () => import('@shikijs/langs/rust'),
  c: () => import('@shikijs/langs/c'),
  cpp: () => import('@shikijs/langs/cpp'),
  java: () => import('@shikijs/langs/java'),
  sql: () => import('@shikijs/langs/sql'),
  html: () => import('@shikijs/langs/html'),
  css: () => import('@shikijs/langs/css'),
  diff: () => import('@shikijs/langs/diff'),
  docker: () => import('@shikijs/langs/docker'),
  toml: () => import('@shikijs/langs/toml'),
