import { Flex } from '@mantine/core'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/(auth)/(app)/datasets/')({
  component: RouteComponent,
  staticData: {
