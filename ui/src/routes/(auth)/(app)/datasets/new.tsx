import { Title, Box } from '@mantine/core'
import { createFileRoute } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

export const Route = createFileRoute(
  '/(auth)/(app)/datasets/new',
)({
  component: RouteComponent,
})

