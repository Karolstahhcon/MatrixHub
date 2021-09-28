import { Link } from '@tanstack/react-router'

import {
  buildModelBadges,
  buildModelMetaItems,
  buildModelTitle,
} from '@/features/models/models.utils.tsx'
import { BaseCard } from '@/shared/components/resource-card/BaseCard.tsx'

import type { Model } from '@matrixhub/api-ts/v1alpha1/model.pb.ts'

interface ModelCardProps {
  model: Model
  fallbackProjectId?: string
}

export function ModelCard({
