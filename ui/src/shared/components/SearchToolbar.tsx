import {
  Flex,
  type FlexProps,
  Group,
} from '@mantine/core'

import { SearchInput, type SearchInputProps } from '@/shared/components/SearchInput'

import classes from './SearchToolbar.module.css'

import type { ReactNode } from 'react'

export interface SearchToolbarProps {
  searchPlaceholder?: string
  searchValue?: SearchInputProps['value']
  onSearchChange?: SearchInputProps['onChange']
  toolbarProps?: Omit<FlexProps, 'children'>
  searchInputProps?: Omit<SearchInputProps, 'placeholder' | 'value' | 'onChange'>
  children?: ReactNode
}

export function SearchToolbar({
