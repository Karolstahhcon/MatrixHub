import { ProjectRoleType } from '@matrixhub/api-ts/v1alpha1/role.pb.ts'
import { createFileRoute } from '@tanstack/react-router'

import { ModelSettingsPage } from '@/features/models/pages/ModelSettingsPage'
import { ensureProjectAccess } from '@/utils/routerAccess'

export const Route = createFileRoute(
  '/(auth)/(app)/projects_/$projectId/models/$modelId/settings/',
)({
  beforeLoad: async ({ params }) => {
