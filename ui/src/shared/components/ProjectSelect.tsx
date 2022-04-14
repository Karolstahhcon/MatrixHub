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
