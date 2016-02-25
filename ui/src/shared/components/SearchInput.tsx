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

export function SearchInput({
  value = '',
  onChange,
  debounceMs = DEFAULT_DEBOUNCE_MS,
  placeholder,
  w = 260,
  styles,
  ...restTextInputProps
}: SearchInputProps) {
  const [inputValue, setInputValue] = useState(value)

  const resolvedStyles = typeof styles === 'function'
    ? styles
    : {
        ...styles,
        input: {
          height: 32,
          minHeight: 32,
