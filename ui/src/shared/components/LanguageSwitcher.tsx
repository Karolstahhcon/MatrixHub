import { Select } from '@mantine/core'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import {
  DEFAULT_LANGUAGE,
  LANGUAGE_STORAGE_KEY,
  normalizeLanguage,
} from '@/i18n'

