import { IconBuildingWarehouse as AdminRegistriesIcon } from '@tabler/icons-react'
import { createFileRoute, stripSearchParams } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

import { AdminPageHeader } from '@/features/admin/components/admin-page-header'
import { AdminPageLayout } from '@/features/admin/components/admin-page-layout'
import { RegistriesPage } from '@/features/admin/registries/pages/RegistriesPage'
import { registriesQueryOptions } from '@/features/admin/registries/registries.query'
import { registriesSearchDefaults, registriesSearchSchema } from '@/features/admin/registries/registries.schema'

