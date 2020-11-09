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
    + 'x'.repeat(contentLength % fillBytes)
  const rawDiff = addedFileDiff(filename, content)

  assert.equal(Buffer.byteLength(rawDiff), byteLength)

  return rawDiff
}

function addedLinesDiff(filename, lines) {
  return [
    `diff --git a/${filename} b/${filename}`,
    'new file mode 100644',
    '--- /dev/null',
    `+++ b/${filename}`,
    `@@ -0,0 +1,${lines.length} @@`,
    ...lines.map(line => `+${line}`),
    '',
  ].join('\n')
}

test('commit diffs follow the Hugging Face file size and file count limits', async () => {
  const server = await createServer({
    appType: 'custom',
    configFile: false,
    optimizeDeps: { noDiscovery: true },
    root,
    server: { middlewareMode: true },
  })

  try {
    const {
      MAX_FILES,
      renderCommitDiff,
      splitCommitDiff,
    } = await server.ssrLoadModule('/src/shared/components/commit-detail/commitDiff.ts')
    const message = 'The diff for this file is too large to render.'
    const linkText = 'See raw diff'

    const atLimitResult = renderCommitDiff(
      sizedFileDiff('at-limit.txt', 100_000, '界'),
      message,
      linkText,
    )

    assert.equal(atLimitResult.hasTooManyFiles, false)
    assert.equal(atLimitResult.html.includes(message), false)
    assert.match(atLimitResult.html, /at-limit\.txt/)

    const oversizedModifiedDiff = [
      'diff --git a/over-limit.txt b/over-limit.txt',
      '--- a/over-limit.txt',
      '+++ b/over-limit.txt',
      '@@ -1 +1,2 @@',
      '-old',
      '+new',
      `+${'界'.repeat(40_000)}`,
