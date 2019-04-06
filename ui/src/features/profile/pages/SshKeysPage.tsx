import {
  Alert,
  Anchor,
  Button,
  Stack,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import {
  IconExternalLink,
  IconInfoCircle,
} from '@tabler/icons-react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { CreateSshKeyModal } from '@/features/profile/components/CreateSshKeyModal'
import { DeleteSshKeyModal } from '@/features/profile/components/DeleteSshKeyModal'
import { SshKeysTable } from '@/features/profile/components/SshKeysTable'
import { profileKeys, sshKeysQueryOptions } from '@/features/profile/profile.query'

import type { SSHKey } from '@matrixhub/api-ts/v1alpha1/current_user.pb'

const SSH_KEY_DOC_URL = '/docs/operations/profile/ssh-key/'

export function SshKeysPage() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const {
    data,
    isFetching,
  } = useQuery(sshKeysQueryOptions())

  const [createOpened, createHandlers] = useDisclosure(false)
  const [deleteOpened, deleteHandlers] = useDisclosure(false)
  const [hintVisible, setHintVisible] = useState(true)
  const [deleteTarget, setDeleteTarget] = useState<SSHKey | null>(null)

  const handleDelete = (key: SSHKey) => {
    setDeleteTarget(key)
    deleteHandlers.open()
  }

  const handleRefresh = () => {
    void queryClient.invalidateQueries({ queryKey: profileKeys.sshKeys })
  }

  const handleDeleteClose = () => {
    deleteHandlers.close()
    setDeleteTarget(null)
  }

  return (
    <Stack gap="sm">
