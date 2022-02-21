import { Select } from '@mantine/core'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import {
  DEFAULT_LANGUAGE,
  LANGUAGE_STORAGE_KEY,
  normalizeLanguage,
} from '@/i18n'

export function LanguageSwitcher() {
  const { i18n } = useTranslation()

  const currentLanguage = useMemo(() => {
    const raw = normalizeLanguage(i18n.language)

    return raw ?? DEFAULT_LANGUAGE
  }, [i18n.language])

  const handleChange = (value: string | null) => {
