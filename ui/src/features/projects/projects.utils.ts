export function isProxyProject(registryUrl?: string) {
  return Boolean(registryUrl?.trim())
}
