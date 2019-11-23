import { useQueryClient } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'

import { useRouteListState } from '@/shared/hooks/useRouteListState'
import { DEFAULT_PAGE } from '@/utils/constants'

import { ReplicationExecutionsTable } from '../components/ReplicationExecutionsTable'
import {
  adminReplicationExecutionKeys,
  useReplicationExecutions,
} from '../executions.query'
import { parseReplicationId } from '../executions.utils'

const replicationExecutionsRouteApi = getRouteApi('/(auth)/admin/replications_/$replicationId/executions')

export function ReplicationExecutionsPage() {
  const queryClient = useQueryClient()
  const navigate = replicationExecutionsRouteApi.useNavigate()
  const search = replicationExecutionsRouteApi.useSearch()
  const { replicationId } = replicationExecutionsRouteApi.useParams()
  const syncPolicyId = parseReplicationId(replicationId)
  const {
    data,
    isLoading,
    isFetching,
  } = useReplicationExecutions(syncPolicyId, search)
  const executions = data?.executions ?? []
  const pagination = data?.pagination

  const refreshExecutions = () => queryClient.invalidateQueries({
    queryKey: adminReplicationExecutionKeys.lists(syncPolicyId),
  })
