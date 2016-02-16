import {
  Group,
  Stack,
  Text,
  rem, Alert,
} from '@mantine/core'
import { IconInfoCircle } from '@tabler/icons-react'
import { useTranslation } from 'react-i18next'

import { CopyValueButton } from '@/shared/components/CopyValueButton.tsx'
import { ModalWrapper } from '@/shared/components/ModalWrapper.tsx'

import type { ReactNode } from 'react'

interface RobotTokenInfoModalProps {
  opened: boolean
  title: ReactNode
  name: string
  token: string
  onClose: () => void
}

interface TokenInfoRowProps {
  label: string
  value: string
  copyable?: boolean
}

function TokenInfoRow({
  label,
  value,
  copyable,
}: TokenInfoRowProps) {
  return (
    <Group gap="lg" wrap="nowrap">
      <Text size="sm" c="dimmed">
        {label}
      </Text>

      <Group gap="sm" align="center" wrap="nowrap">
        <Text
          size="sm"
          style={{ wordBreak: 'break-all' }}
        >
          {value}
        </Text>

        {copyable && (
          <CopyValueButton
            value={value}
            iconSize={16}
          />
        )}
      </Group>
    </Group>
  )
}

export function RobotTokenInfoModal({
  opened,
  title,
  name,
  token,
  onClose,
