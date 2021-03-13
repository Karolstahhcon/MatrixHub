import { SystemService } from '@matrixhub/api-ts/v1alpha1/system.pb'
import { queryOptions, useQuery } from '@tanstack/react-query'

export const systemKeys = {
  all: ['system'] as const,
  config: () => [...systemKeys.all, 'config'] as const,
}

