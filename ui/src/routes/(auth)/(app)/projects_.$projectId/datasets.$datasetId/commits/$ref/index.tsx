import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute(
  '/(auth)/(app)/projects_/$projectId/datasets/$datasetId/commits/$ref/',
)({
  component: RouteComponent,
})

function RouteComponent() {
  const {
