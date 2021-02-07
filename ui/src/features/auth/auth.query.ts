import { CurrentUser } from '@matrixhub/api-ts/v1alpha1/current_user.pb'
import { queryOptions, useQuery } from '@tanstack/react-query'

import { queryClient } from '@/queryClient'
// -- Query key factory --
export const authKeys = {
  all: ['auth'] as const,
