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
}

export function renderCommitDiff(
  rawDiff: string,
  diffTooBigMessage: string,
  rawDiffLinkText: string,
): CommitDiffRender {
  const fileDiffs = splitCommitDiff(rawDiff)
  const diffJson = fileDiffs.slice(0, MAX_FILES).flatMap((fileDiff, fileIndex) => {
    const parsedDiff = parse(fileDiff)

    // Parsing is cheap; rendering (especially side-by-side line matching) is what
    // blows up on large files. Keep the parser's line counts, drop the blocks so
    // the renderer only emits the placeholder message.
