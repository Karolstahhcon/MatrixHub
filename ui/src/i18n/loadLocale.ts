type LocaleResource = Record<string, unknown>

const modules = import.meta.glob<Record<string, LocaleResource>>(
  '../locales/**/*.json',
  {
    eager: true,
  },
)

/**
 * Sort locale files so shallower paths load first.
 * Deeper paths load later and override duplicate keys
 * (e.g. `shared/components/A.json` overrides keys in `shared.json`).
 */
function sortLocalePaths(paths: string[]) {
  return [...paths].sort((a, b) => {
    const depthDiff = a.split('/').length - b.split('/').length

    if (depthDiff !== 0) {
      return depthDiff
    }

    return a.localeCompare(b)
  })
}

function isLocaleResource(value: unknown): value is LocaleResource {
  return typeof value === 'object' && value !== null && !Array.isArray(value)
}

function warnDuplicateKey(keyPath: string, sourcePath: string) {
  if (import.meta.env.DEV) {
    console.warn(
      `[i18n] Duplicate key "${keyPath}": overridden by "${sourcePath}".`,
    )
  }
}

/**
 * Wrap a value in nested objects following the given path segments.
