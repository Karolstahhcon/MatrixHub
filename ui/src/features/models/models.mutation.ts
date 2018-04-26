import {
  type CreateModelRequest,
  type DeleteModelRequest,
  type UpdateModelSettingRequest,
  Models,
} from '@matrixhub/api-ts/v1alpha1/model.pb'
import { mutationOptions } from '@tanstack/react-query'

import { modelKeys } from '@/features/models/models.query'
import i18n from '@/i18n'

import type { NotificationMeta } from '@/types/tanstack-query'

export function deleteModelMutationOptions() {
  return mutationOptions({
    mutationFn: (params: DeleteModelRequest) => Models.DeleteModel(params),
