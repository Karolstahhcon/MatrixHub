import { Projects } from '@matrixhub/api-ts/v1alpha1/project.pb'
import { Robots } from '@matrixhub/api-ts/v1alpha1/robot.pb'
import { Roles } from '@matrixhub/api-ts/v1alpha1/role.pb'
import {
  keepPreviousData,
  queryOptions,
  useQuery,
} from '@tanstack/react-query'

import type { RobotsSearch } from '@/routes/(auth)/admin/robots'
import type { NotificationMeta } from '@/types/tanstack-query'

export const adminRobotKeys = {
  all: ['admin', 'robots'] as const,
  lists: () => [...adminRobotKeys.all, 'list'] as const,
  list: (search: RobotsSearch) => [...adminRobotKeys.lists(), search] as const,
  details: () => [...adminRobotKeys.all, 'detail'] as const,
  detail: (robotId: number) => [...adminRobotKeys.details(), robotId] as const,
  projects: () => [...adminRobotKeys.all, 'projects'] as const,
