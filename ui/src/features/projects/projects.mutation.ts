import {
  ProjectType, Projects, type UpdateProjectRequest,
} from '@matrixhub/api-ts/v1alpha1/project.pb'
import { mutationOptions } from '@tanstack/react-query'

import { authKeys } from '@/features/auth/auth.query'
import i18n from '@/i18n'

import { projectKeys } from './projects.query'

import type { CreateProjectInput } from './projects.schema'
import type { NotificationMeta } from '@/types/tanstack-query'

