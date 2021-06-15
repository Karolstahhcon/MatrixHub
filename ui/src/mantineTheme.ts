import { createTheme, rem } from '@mantine/core'
import { type CSSVariablesResolver } from '@mantine/core'

import { themeComponents } from './mantine-theme/components'

export const mantineTheme = createTheme({
  primaryColor: 'cyan',
  components: themeComponents,
})

export const cssVariablesResolver: CSSVariablesResolver = () => ({
  variables: {
    '--app-size-icon-basic': rem(20),
    '--app-size-icon-md': rem(16),
    '--app-size-icon-sm': rem(16),
    '--app-size-radius-mmd': rem(6),
  },
  light: {
    '--mantine-color-text': 'var(--mantine-color-gray-9)',
    '--app-color-text-default': 'var(--mantine-color-gray-9)',
    '--app-color-text-label': 'var(--mantine-color-gray-7)',
    '--app-color-text-link': 'var(--mantine-color-blue-6)',
    '--app-color-text-title': 'var(--mantine-color-gray-9)',
    '--app-color-text-error-default': 'var(--mantine-color-red-6)',
    '--app-color-text-error-disabled': 'var(--mantine-color-red-4)',
    '--app-color-background-body': 'var(--mantine-color-white)',
    '--app-color-background-success-filled-hover': 'var(--mantine-color-teal-7)',
    '--app-color-border-default': '1px solid var(--mantine-color-gray-3)',
    '--app-color-border-error-default': '1px solid var(--mantine-color-red-6)',
    '--app-color-gray-10': 'var(--mantine-color-gray-0)',
    '--app-color-gray-20': 'var(--mantine-color-gray-1)',
    '--app-color-gray-30': 'var(--mantine-color-gray-2)',
