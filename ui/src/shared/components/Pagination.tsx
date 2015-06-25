import {
  Group,
  Pagination as MantinePagination,
  type PaginationProps as MantinePaginationProps,
  Text,
  type TextProps,
} from '@mantine/core'
import { startTransition } from 'react'
import { useTranslation } from 'react-i18next'

import type { ReactNode } from 'react'

export interface PaginationProps {
  total: number
  totalPages: number
  page: number
  onPageChange: (page: number) => void
  totalLabel?: ReactNode
  totalLabelProps?: TextProps
  paginationProps?: Omit<
