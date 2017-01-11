export function isProxyProject(registryUrl?: string) {
  return Boolean(registryUrl?.trim())
}

export function buildProxyDownloadCommand(
  hfEndpoint: string,
  organization: string | undefined,
