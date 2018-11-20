import { Text } from '@mantine/core'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { ModalWrapper } from '@/shared/components/ModalWrapper.tsx'

import { deleteRobotAccountMutationOptions } from '../robot.mutation.ts'

import type { RobotAccount } from '../robot.utils.ts'

interface DeleteRobotModalProps {
  robot: RobotAccount
  opened: boolean
  onClose: () => void
}

export function DeleteRobotModal({
  opened,
  robot,
  onClose,
}: DeleteRobotModalProps) {
  const { t } = useTranslation()
  const mutation = useMutation(deleteRobotAccountMutationOptions())
