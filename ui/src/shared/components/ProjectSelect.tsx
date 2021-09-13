import {
  CheckIcon,
  Combobox,
  Group,
  InputBase,
  type InputBaseProps,
  ScrollArea,
  Text,
  useCombobox,
} from '@mantine/core'
import { type ComponentProps, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { ProjectTypeBadge } from '@/shared/components/badges/ProjectTypeBadge'
import { FieldHintLabel } from '@/shared/components/FieldHintLabel.tsx'
import { filterByKeyword } from '@/shared/utils'

export interface ProjectSelectOption {
  name?: string
  type?: ComponentProps<typeof ProjectTypeBadge>['type']
}

export interface ProjectSelectProps {
  data?: ProjectSelectOption[]
  value?: string
  onChange?: (value: string) => void
  label?: InputBaseProps['label']
  withAsterisk?: InputBaseProps['withAsterisk']
  inputProps?: Omit<
    InputBaseProps,
    | 'component'
    | 'type'
    | 'children'
    | 'rightSection'
    | 'rightSectionPointerEvents'
    | 'onClick'
  > & {
    onBlur?: () => void
  }
}

interface SelectedProjectDisplayProps {
  name?: string
  type?: ComponentProps<typeof ProjectTypeBadge>['type']
}

function SelectedProjectDisplay({
  name,
  type,
}: SelectedProjectDisplayProps) {
  return (
    <Group gap={6} wrap="nowrap">
      <Text
        title={name}
        size="sm"
        truncate
      >
        {name ?? ''}
      </Text>
      <ProjectTypeBadge
        type={type}
        flex="0 0 auto"
      />
    </Group>
  )
}

const EMPTY_OPTIONS: ProjectSelectOption[] = []

export function ProjectSelect({
  data = EMPTY_OPTIONS,
  value,
  onChange,
  label,
  withAsterisk = true,
  inputProps,
}: ProjectSelectProps) {
  const { t } = useTranslation()
  const [search, setSearch] = useState('')
  const combobox = useCombobox({
    onDropdownClose: () => {
      combobox.resetSelectedOption()
      combobox.focusTarget()
      setSearch('')
    },
    onDropdownOpen: () => {
      combobox.updateSelectedOptionIndex('active', { scrollIntoView: true })
      combobox.focusSearchInput()
    },
