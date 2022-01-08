import { Button } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { IconRobot } from '@tabler/icons-react'
import { useQueryClient } from '@tanstack/react-query'
import { getRouteApi } from '@tanstack/react-router'
import {
  startTransition,
  useState,
} from 'react'
import { useTranslation } from 'react-i18next'

import { DeleteRobotModal } from '../components/DeleteRobotModal.tsx'
import { DisableRobotModal } from '../components/DisableRobotModal.tsx'
import { EnableRobotModal } from '../components/EnableRobotModal.tsx'
import { RefreshRobotTokenModal } from '../components/RefreshRobotTokenModal.tsx'
import { RobotTable } from '../components/RobotTable.tsx'
import { RobotTokenInfoModal } from '../components/RobotTokenInfoModal.tsx'
import {
  adminRobotKeys,
  useRobots,
} from '../robot.query.ts'
import {
  isRobotEnabled,
  type RobotAccount,
} from '../robot.utils.ts'

const robotsRouteApi = getRouteApi('/(auth)/admin/robots/')

type RobotActionModalType = 'delete' | 'enable' | 'disable' | 'refreshToken'

interface RobotActionModalPayload {
  type: RobotActionModalType
  robot: RobotAccount
}

interface RobotTokenInfoPayload {
  name: string
  token: string
}

export function RobotsPage() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const navigate = robotsRouteApi.useNavigate()
  const search = robotsRouteApi.useSearch()
  const {
    data,
    isLoading,
    isFetching,
  } = useRobots(search)
  const robots = data?.items ?? []
  const pagination = data?.pagination

  const refreshRobots = () => queryClient.invalidateQueries({
    queryKey: adminRobotKeys.lists(),
  })

  const handleSearchChange = (value: string) => {
    if (value === (search.query ?? '')) {
      return
    }

    void navigate({
      replace: true,
      search: prev => ({
        ...prev,
        page: 1,
        query: value,
      }),
    })
  }

  const handlePageChange = (page: number) => {
    if (page === (search.page ?? 1)) {
      return
    }

    startTransition(() => {
      void navigate({
        search: prev => ({
          ...prev,
          page,
        }),
      })
    })
  }

  const handleCreate = () => {
    void navigate({ to: '/admin/robots/new' })
  }

  const handleEdit = (robot: RobotAccount) => {
    if (robot.id == null) {
      return
    }

    void navigate({
      to: '/admin/robots/$robotId/edit',
      params: { robotId: String(robot.id) },
    })
  }

  const [actionModalOpened, actionModalHandlers] = useDisclosure(false)
  const [actionModalPayload, setActionModalPayload] = useState<RobotActionModalPayload | null>(null)

  const openActionModal = (type: RobotActionModalType, robot: RobotAccount) => {
    setActionModalPayload({
      type,
      robot,
    })
    actionModalHandlers.open()
  }

  const closeActionModal = () => {
    actionModalHandlers.close()
    setActionModalPayload(null)
  }

  const [tokenInfoModalOpened, tokenInfoModalHandlers] = useDisclosure(false)
  const [tokenInfoPayload, setTokenInfoPayload] = useState<RobotTokenInfoPayload | null>(null)

  const openTokenInfoModal = (robot: RobotAccount, token: string) => {
    setTokenInfoPayload({
      name: robot.name ?? '',
