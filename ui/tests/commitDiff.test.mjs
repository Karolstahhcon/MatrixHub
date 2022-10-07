import assert from 'node:assert/strict'
import { test } from 'node:test'
import { fileURLToPath } from 'node:url'

import { createServer } from 'vite'

const root = fileURLToPath(new URL('../', import.meta.url))

function addedFileDiff(filename, content) {
  return [
    `diff --git a/${filename} b/${filename}`,
    'new file mode 100644',
    '--- /dev/null',
    `+++ b/${filename}`,
    '@@ -0,0 +1 @@',
    `+${content}`,
    '',
  ].join('\n')
}

function sizedFileDiff(filename, byteLength, fill = 'x') {
  const emptyDiff = addedFileDiff(filename, '')
  const contentLength = byteLength - Buffer.byteLength(emptyDiff)

  assert.ok(contentLength >= 0)

  const fillBytes = Buffer.byteLength(fill)
  const content = fill.repeat(Math.floor(contentLength / fillBytes))
