import { html as formatDiffHtml, parse } from 'diff2html'

export const MAX_FILES = 50
const MAX_FILE_DIFF_BYTES = 100_000

const textEncoder = new TextEncoder()

interface CommitDiffRender {
  html: string
  hasTooManyFiles: boolean
}

export function splitCommitDiff(rawDiff: string): string[] {
  return rawDiff.split(/(?=^diff --git )/m).filter(Boolean)
