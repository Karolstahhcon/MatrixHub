import { useQueryClient } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'

import { usePayloadModal } from '@/shared/hooks/usePayloadModal'
import { useRouteListState } from '@/shared/hooks/useRouteListState'

import { BatchDeleteUsersModal } from '../components/BatchDeleteUsersModal'
import { CreateUserAction } from '../components/CreateUserAction'
import { UsersTable } from '../components/UsersTable'
import {
  adminUserKeys,
  useUsers,
} from '../users.query'
import { DEFAULT_USERS_PAGE } from '../users.schema'
import { getUserRowId } from '../users.utils'

import type { User } from '@matrixhub/api-ts/v1alpha1/user.pb'

const usersRouteApi = getRouteApi('/(auth)/admin/users')

export function UsersPage() {
  const queryClient = useQueryClient()
  const navigate = usersRouteApi.useNavigate()
  const search = usersRouteApi.useSearch()
  const {
    data,
    isLoading,
    isFetching,
  } = useUsers(search)
  const users = data?.users ?? []
  const pagination = data?.pagination
  const batchDeleteModal = usePayloadModal<User[]>()

  const refreshUsers = () => queryClient.invalidateQueries({
    queryKey: adminUserKeys.lists(),
  })

  const {
    rowSelection,
    setRowSelection,
    clearRowSelection,
    selectedCount,
    selectedRecords,
    onSearchChange,
    onRefresh,
    onPageChange,
  } = useRouteListState({
