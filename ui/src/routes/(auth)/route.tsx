import {
  AppShell,
  Avatar,
  Flex,
  Group,
  Menu,
  NavLink,
  Text,
  UnstyledButton,
  rem,
} from '@mantine/core'
import { Login } from '@matrixhub/api-ts/v1alpha1/login.pb'
import {
  IconChevronDown as ArrowDownIcon,
  IconCube as ModelIcon,
  IconFileText as DocsIcon,
  IconLogout as LogOutIcon,
  IconSettings as SettingsIcon,
  IconUser as UserIcon,
  IconApiApp as ProjectIcon,
} from '@tabler/icons-react'
import { useMutation } from '@tanstack/react-query'
import {
  createFileRoute,
  Link,
  linkOptions,
  Outlet,
  useMatchRoute,
  CatchBoundary,
  useRouterState, redirect,
} from '@tanstack/react-router'
import { use } from 'react'
import { useTranslation } from 'react-i18next'

import LogoIcon from '@/assets/svgs/logo.svg?react'
import { CurrentUserContext } from '@/context/current-user-context.tsx'
import { currentUserQueryOptions } from '@/features/auth/auth.query'
import { queryClient } from '@/queryClient'
import { Route as ModelsRoute } from '@/routes/(auth)/(app)/models'
import { Route as CreateModelRoute } from '@/routes/(auth)/(app)/models/new'
import { Route as ProfileRoute } from '@/routes/(auth)/(app)/profile'
import { Route as ProjectsRoute } from '@/routes/(auth)/(app)/projects'
import { Route as ProjectModelRoute } from '@/routes/(auth)/(app)/projects_.$projectId/models.$modelId/route'
import { Route as AdminRoute } from '@/routes/(auth)/admin'
import { LanguageSwitcher } from '@/shared/components/LanguageSwitcher'
import { RouterErrorComponent } from '@/shared/components/RouterErrorComponent'
import { RouteStatusPage } from '@/shared/components/RouteStatusPage'
import {
  isForbiddenRouteError, isNotFoundRouteError,
  isSdkNotFound, isSdkPermissionDenied,
  isSdkUnauthenticated,
} from '@/utils/routerAccess'

export const Route = createFileRoute('/(auth)')({
  component: AuthLayout,
  loader: async () => {
    try {
      const user = await queryClient.ensureQueryData(currentUserQueryOptions())

      return {
        user,
      }
    } catch {
      throw redirect({ to: '/login' })
    }
  },
})

function AppLogo() {
  return (
    <UnstyledButton
      component={Link}
      to="/"
    >
      <Group
        gap={8}
        wrap="nowrap"
      >
        <LogoIcon fontSize={rem(36)} />
        <Text
          fw={600}
          size="xl"
          c="var(--mantine-color-black)"
        >
          MatrixHub
        </Text>
      </Group>
    </UnstyledButton>
  )
}

const DOCS_URL = '/'

function AppNavbar() {
  const { t } = useTranslation()
  const navRoutes = linkOptions([
    {
      label: t('nav.models'),
      icon: ModelIcon,
      to: ModelsRoute.to,
      extraMatch: ProjectModelRoute.to,
    },
    {
      label: t('nav.projectManagement'),
      icon: ProjectIcon,
      to: ProjectsRoute.to,
      repelMatch: [ProjectModelRoute.to],
    },
  ])
  const matchRoute = useMatchRoute()

  const navLinkStyles = (isActive: boolean) => ({
    root: {
      width: 'auto',
      height: '32px',
      borderRadius: 'var(--mantine-radius-lg)',
      fontWeight: '600',
      color: isActive ? 'var(--nl-color)' : '#868E96',
    },
