export function filterByKeyword<T extends { name?: string }>(items: T[], keyword: string) {
  const normalized = keyword.trim().toLowerCase()

  if (!normalized) {
    return items
  }

  return items.filter(item => item.name?.toLowerCase()?.includes(normalized))
}

export function parseConventionalCommit(commitMessage: string) {
  const msg = commitMessage.trim()
  const lines = msg.split('\n').map(line => line.trimEnd())

  const result: Record<string, string | null> = {
    header: '',
    body: null,
    footer: null,
  }

  result.header = lines[0] || ''

