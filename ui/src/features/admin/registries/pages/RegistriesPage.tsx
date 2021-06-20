import { useQueryClient } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'

import { useRouteListState } from '@/shared/hooks/useRouteListState'

import { CreateRegistryAction } from '../components/CreateRegistryAction'
import { RegistriesTable } from '../components/RegistriesTable'
import {
  adminRegistryKeys,
  useRegistries,
} from '../registries.query'
