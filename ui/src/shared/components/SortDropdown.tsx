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
  order: SortOrder
  orderLabel?: string
  refreshing?: boolean
  width?: number
  orderWidth?: number
  onFieldChange: (value: string) => void
  onToggleOrder: () => void
}

const DEFAULT_WIDTH = 204
const DEFAULT_ORDER_WIDTH = 80

export function SortDropdown({
  fieldOptions,
  fieldValue,
  order,
  orderLabel,
  refreshing = false,
  width = DEFAULT_WIDTH,
  orderWidth = DEFAULT_ORDER_WIDTH,
  onFieldChange,
  onToggleOrder,
}: SortDropdownProps) {
  const { t } = useTranslation()
  const [opened, setOpened] = useState(false)
  const menuId = useId()
  const resolvedOrderLabel = orderLabel
    ?? t(order === 'asc' ? 'shared.sortOrderAsc' : 'shared.sortOrderDesc')
  const fieldWidth = width - orderWidth
  const selectedOption = fieldOptions.find(option => option.value === fieldValue)
    ?? fieldOptions[0]
  const cssVariables = {
    '--sort-width': `${width}px`,
    '--sort-field-width': `${fieldWidth}px`,
    '--sort-order-width': `${orderWidth}px`,
  } as CSSProperties

  return (
    <Popover
      opened={opened}
      onChange={setOpened}
      width={width}
      position="bottom-start"
      offset={8}
      withArrow
      arrowSize={10}
      withinPortal={false}
    >
      <Popover.Target>
        <Box className={classes.anchor} style={cssVariables}>
          <Box className={classes.control}>
            <UnstyledButton
              type="button"
              className={classes.fieldButton}
              aria-controls={menuId}
              aria-expanded={opened}
              aria-haspopup="menu"
              onClick={() => {
                setOpened(current => !current)
              }}
            >
              <span className={classes.label}>
