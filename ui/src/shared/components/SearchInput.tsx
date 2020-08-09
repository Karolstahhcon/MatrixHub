import { TextInput, type TextInputProps } from '@mantine/core'
import { useDebouncedCallback } from '@mantine/hooks'
import { IconSearch } from '@tabler/icons-react'
import {
  startTransition,
  useEffect,
  useState,
} from 'react'

export interface SearchInputProps extends Omit<TextInputProps, 'value' | 'onChange'> {
  value?: string
  onChange?: (value: string) => void
  debounceMs?: number
}

const DEFAULT_DEBOUNCE_MS = 300

