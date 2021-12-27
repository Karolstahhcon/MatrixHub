import { Tabs } from '@mantine/core'

import type { FilterTabDefinition } from './types.ts'

interface ResourceFilterPanelProps {
  tabs: FilterTabDefinition[]
  activeTab: string
  onTabChange: (nextTab: string) => void
}

export function ResourceFilterPanels({
