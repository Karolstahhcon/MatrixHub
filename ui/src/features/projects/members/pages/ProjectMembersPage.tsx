import { Box, Button } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { ProjectRoleType } from '@matrixhub/api-ts/v1alpha1/role.pb'
import { IconPlus } from '@tabler/icons-react'
import { getRouteApi } from '@tanstack/react-router'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { useCurrentUser } from '@/features/auth/auth.query'
import { useProjectRole } from '@/features/auth/useProjectRole'

import { AddMemberModal } from '../components/AddMemberModal'
import { EditRoleModal } from '../components/EditRoleModal'
import { MembersTable } from '../components/MembersTable'
import { RemoveMemberModal } from '../components/RemoveMemberModal'
import { useMembers } from '../members.query'

import type { ProjectMember } from '@matrixhub/api-ts/v1alpha1/project.pb'
import type { MRT_RowSelectionState } from 'mantine-react-table'

const membersRouteApi = getRouteApi('/(auth)/(app)/projects/$projectId/members/')

export function ProjectMembersPage() {
  const { t } = useTranslation()
  const { projectId } = membersRouteApi.useParams()
  const navigate = membersRouteApi.useNavigate()
  const search = membersRouteApi.useSearch()
  const currentRole = useProjectRole(projectId)
  const { data: currentUser } = useCurrentUser()

  const isAdmin = currentRole === ProjectRoleType.ROLE_TYPE_PROJECT_ADMIN

  const {
    data, isLoading, isFetching, refetch,
  } = useMembers(projectId, {
    q: search.q,
    page: search.page,
  })

  const members = data?.members ?? []
  const pagination = data?.pagination

  const [rowSelection, setRowSelection] = useState<MRT_RowSelectionState>({})

  // Modal states
  const [addOpened, addHandlers] = useDisclosure(false)
  const [editOpened, editHandlers] = useDisclosure(false)
  const [removeOpened, removeHandlers] = useDisclosure(false)
  const [batchRemoveOpened, batchRemoveHandlers] = useDisclosure(false)

  const [processMember, setProcessMember] = useState<ProjectMember | null>(null)

  const handleSearchChange = (value: string) => {
    if (value === search.q) {
      return
    }

    setRowSelection({})
    void navigate({
      replace: true,
      search: prev => ({
        ...prev,
        page: 1,
        q: value,
      }),
    })
  }

  const handlePageChange = (page: number) => {
    setRowSelection({})
    void navigate({
      search: prev => ({
        ...prev,
        page,
      }),
    })
  }

  const handleRefresh = () => {
    setRowSelection({})
    void refetch()
  }

  const selectedMembers = members.filter((m) => {
    const key = `${m.memberType}:${m.memberId}`

    return !!rowSelection[key]
  })

  const handleEditMember = (member: ProjectMember) => {
    setProcessMember(member)
    editHandlers.open()
  }

  const handleRemoveMember = (member: ProjectMember) => {
