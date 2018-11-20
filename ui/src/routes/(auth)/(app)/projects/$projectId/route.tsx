import {
  Badge, Box, Group, Space, Tabs, Text,
} from '@mantine/core'
import { ProjectRoleType } from '@matrixhub/api-ts/v1alpha1/role.pb'
import { IconApiApp as ProjectIcon } from '@tabler/icons-react'
import { useSuspenseQuery } from '@tanstack/react-query'
import {
  Link,
  createFileRoute,
  linkOptions,
  Outlet,
  useMatchRoute,
  useLocation,
} from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

import { useProjectRole } from '@/features/auth/useProjectRole'
import { projectDetailQueryOptions } from '@/features/projects/projects.query'
import { isProxyProject } from '@/features/projects/projects.utils'

import { Route as ProjectMembersRoute } from './members'
import { Route as ProjectModelsRoute } from './models'
import { Route as ProjectSettingsRoute } from './settings'

export const Route = createFileRoute('/(auth)/(app)/projects/$projectId')({
