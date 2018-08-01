import { Button } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { IconHomePlus } from '@tabler/icons-react'
import { getRouteApi } from '@tanstack/react-router'
import { useEffect, useEffectEvent } from 'react'
import { useTranslation } from 'react-i18next'

import { RegistryFormModal } from './RegistryFormModal'

const registriesRouteApi = getRouteApi('/(auth)/admin/registries')

export function CreateRegistryAction() {
  const { t } = useTranslation()
  const navigate = registriesRouteApi.useNavigate()
  const { create } = registriesRouteApi.useSearch()
