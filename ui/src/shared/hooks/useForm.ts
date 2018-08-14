import { useForm as _useForm } from '@tanstack/react-form'
import { useEffect, useRef } from 'react'
import { useTranslation } from 'react-i18next'

/**
 * Wraps TanStack's useForm with automatic i18n error refresh.
 * When the language changes, re-validates only fields that already have errors
 * so their messages update to the new language.
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const useForm = ((...args: any[]) => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const form = (_useForm as any)(...args)
