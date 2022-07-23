import {
  AppShell,
  Group,
  NavLink,
  ScrollArea,
  Stack,
  Text,
  rem,
} from '@mantine/core'
import { IconSettings as SettingsIcon } from '@tabler/icons-react'
import {
  createFileRoute,
  Link,
  linkOptions,
  Outlet,
  redirect,
  useMatchRoute,
} from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

import { currentUserQueryOptions } from '@/features/auth/auth.query'
import { queryClient } from '@/queryClient'
import { Route as AdminRegistriesRoute, Icon as AdminRegistriesIcon } from '@/routes/(auth)/admin/registries'
import { Route as AdminReplicationsRoute, Icon as AdminReplicationsIcon } from '@/routes/(auth)/admin/replications'
import { Route as AdminRobotsRoute, Icon as AdminRobotsIcon } from '@/routes/(auth)/admin/robots/index'
import { Route as AdminUsersRoute, Icon as AdminUsersIcon } from '@/routes/(auth)/admin/users'
import { setAdminContentViewport } from '@/utils/setContentViewport'

export const Route = createFileRoute('/(auth)/admin')({
  beforeLoad: async () => {
    const user = await queryClient.ensureQueryData(currentUserQueryOptions())

    if (!user.isAdmin) {
      throw redirect({
        to: '/',
      })
    }
  },
  component: AdminLayout,
  staticData: {
    navName: 'Admin',
  },
})

function AdminNavbar() {
  const { t } = useTranslation()
  const navRoutes = linkOptions([
    {
      label: t('admin.users'),
      icon: AdminUsersIcon,
      to: AdminUsersRoute.to,
    },
    {
      label: t('admin.robots'),
      icon: AdminRobotsIcon,
      to: AdminRobotsRoute.to,
    },
    {
      label: t('admin.registries'),
      icon: AdminRegistriesIcon,
      to: AdminRegistriesRoute.to,
