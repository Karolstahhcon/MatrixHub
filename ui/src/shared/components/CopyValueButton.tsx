import {
  ActionIcon,
  Tooltip,
} from '@mantine/core'
import { IconCopy } from '@tabler/icons-react'
import { useTranslation } from 'react-i18next'

import { CopyButton } from '@/shared/components/CopyButton'

import type { ReactNode } from 'react'

interface CopyRenderProps {
  copied: boolean
  copy: () => void
}

interface CopyValueButtonProps {
  value: string
  timeout?: number
  iconSize?: number
  iconClassName?: string
  color?: string
  children?: (props: CopyRenderProps) => ReactNode
}

export function CopyValueButton({
  value,
  color = 'gray',
  timeout = 1500,
  iconClassName,
  iconSize = 16,
  children,
}: CopyValueButtonProps) {
  const { t } = useTranslation()

