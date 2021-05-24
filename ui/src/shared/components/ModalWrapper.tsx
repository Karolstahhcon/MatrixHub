import {
  Button, Group, Modal, Stack, Text,
  ScrollArea,
} from '@mantine/core'
import {
  IconAlertCircle, IconCircleCheckFilled, IconInfoCircle,
} from '@tabler/icons-react'
import { useTranslation } from 'react-i18next'

import type { ModalProps } from '@mantine/core'
import type { ReactNode } from 'react'

const modalTypeIcons = {
  info: <IconInfoCircle size={20} color="var(--mantine-color-cyan-6)" />,
  danger: <IconAlertCircle size={20} color="var(--mantine-color-red-6)" />,
  success: <IconCircleCheckFilled size={20} color="var(--mantine-color-green-6)" />,
}

export type ModalWrapperProps = Omit<ModalProps, 'title'> & {
  type?: 'info' | 'danger' | 'success'
  title: ReactNode
  footer?: ReactNode
  confirmLoading?: boolean
  onConfirm?: () => void
}

export function ModalWrapper({
