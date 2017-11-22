import {
  Box,
  Popover,
  UnstyledButton,
} from '@mantine/core'
import { IconSortDescending as SortDescIcon } from '@tabler/icons-react'
import {
  useId,
  useState,
} from 'react'
import { useTranslation } from 'react-i18next'

import classes from './SortDropdown.module.css'

import type {
  CSSProperties,
  ReactNode,
} from 'react'

export type SortOrder = 'asc' | 'desc'

export interface SortDropdownOption {
  value: string
  label: string
  icon?: ReactNode
  disabled?: boolean
}

interface SortDropdownProps {
  fieldOptions: readonly SortDropdownOption[]
  fieldValue: string
