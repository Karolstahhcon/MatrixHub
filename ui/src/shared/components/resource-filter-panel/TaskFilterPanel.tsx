import { Group } from '@mantine/core'
import { type Label } from '@matrixhub/api-ts/v1alpha1/model.pb.ts'

import { TaskBadge } from '@/shared/components/badges/TaskBadge'
import { BaseFilterPanel } from '@/shared/components/resource-filter-panel/BaseFilterPanel'
import { filterByKeyword } from '@/shared/utils'

interface TaskFilterPanelProps {
  options: Label[]
  loading: boolean
  searchPlaceholder?: string
  selectedNames: string[]
  onSelect: (name: string) => void
}

export function TaskFilterPanel({
  options,
  loading,
  searchPlaceholder = '',
  selectedNames,
  onSelect,
}: TaskFilterPanelProps) {
  return (
    <BaseFilterPanel
      loading={loading}
      searchPlaceholder={searchPlaceholder}
    >
      {({ keyword }) => {
        const filteredOptions = filterByKeyword(options, keyword)

