import {
  Box, Group, Stack, Tabs, Text,
} from '@mantine/core'
import { IconUser } from '@tabler/icons-react'
import {
  createFileRoute, Link, linkOptions, Outlet, useMatchRoute,
} from '@tanstack/react-router'
import { use } from 'react'
import { useTranslation } from 'react-i18next'

import DefaultAvatarIcon from '@/assets/svgs/default-avatar.svg?react'
import { CurrentUserContext } from '@/context/current-user-context'
import { Route as AccessTokenRoute } from '@/routes/(auth)/(app)/profile/access-token'
import { Route as SecurityRoute } from '@/routes/(auth)/(app)/profile/security'
import { Route as SshKeysRoute } from '@/routes/(auth)/(app)/profile/ssh-keys'

export const Route = createFileRoute('/(auth)/(app)/profile')({
  component: Profile,
})

function Profile() {
  const { t } = useTranslation()
  const currentUser = use(CurrentUserContext)

  const profileTabs = linkOptions([
    {
      label: t('profile.securitySetting'),
      value: 'security',
      to: SecurityRoute.to,
    },
    {
      label: t('profile.accessToken'),
      value: 'access-token',
      to: AccessTokenRoute.to,
    },
    {
      label: t('profile.sshKey.publicKey'),
      value: 'ssh-keys',
      to: SshKeysRoute.to,
    },
