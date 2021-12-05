import {
  Group,
  Loader,
  Stack,
  Text,
} from '@mantine/core'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { SearchInput } from '@/shared/components/SearchInput'

import type { ReactNode } from 'react'

interface BaseFilterTabProps {
  loading: boolean
  searchPlaceholder: string
  children?: (params: { keyword: string }) => ReactNode
}

export function BaseFilterPanel({
  searchPlaceholder,
  children,
  loading,
}: BaseFilterTabProps) {
  const { t } = useTranslation()
  const [keyword, setKeyword] = useState('')
  const content = children?.({ keyword }) ?? null

  return (
