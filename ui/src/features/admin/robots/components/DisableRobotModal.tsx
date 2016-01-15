import { Text } from '@mantine/core'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { ModalWrapper } from '@/shared/components/ModalWrapper.tsx'

import { disableRobotAccountMutationOptions } from '../robot.mutation.ts'

import type { RobotAccount } from '../robot.utils.ts'

interface DisableRobotModalProps {
  robot: RobotAccount
  opened: boolean
  onClose: () => void
}
