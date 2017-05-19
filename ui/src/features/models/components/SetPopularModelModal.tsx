import { rem, Text } from '@mantine/core'
import { useMutation } from '@tanstack/react-query'
import { useRouter } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

import { updateModelSettingMutationOptions } from '@/features/models/models.mutation'
import { ModalWrapper } from '@/shared/components/ModalWrapper'

interface RecommendedModelModalProps {
  opened: boolean
  onClose: () => void
  projectId: string
  modelId: string
  isPopular: boolean
}

export function SetPopularModelModal({
  opened,
  onClose,
  projectId,
  modelId,
  isPopular,
}: RecommendedModelModalProps) {
  const { t } = useTranslation()
