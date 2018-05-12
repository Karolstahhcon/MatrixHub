import { Button } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { useTranslation } from 'react-i18next'

import { Icon } from '@/routes/(auth)/admin/replications'

import { ReplicationFormModal } from './ReplicationFormModal'

export function CreateReplicationAction() {
  const { t } = useTranslation()
  const [opened, {
    open, close,
  }] = useDisclosure(false)
