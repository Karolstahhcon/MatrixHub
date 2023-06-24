import { useQueryClient } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'

import { useRouteListState } from '@/shared/hooks/useRouteListState'

import { CreateReplicationAction } from '../components/CreateReplicationAction'
import { ReplicationsTable } from '../components/ReplicationsTable'
import {
  adminReplicationKeys,
  useReplications,
} from '../replications.query'
