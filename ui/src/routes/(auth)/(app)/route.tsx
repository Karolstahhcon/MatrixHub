import { Box, ScrollArea } from '@mantine/core'
import {
  createFileRoute,
  Outlet,
} from '@tanstack/react-router'

import { setContentViewport } from '@/utils/setContentViewport'
export const Route = createFileRoute('/(auth)/(app)')({
  component: AppLayout,
})

function AppLayout() {
  return (
    <ScrollArea
