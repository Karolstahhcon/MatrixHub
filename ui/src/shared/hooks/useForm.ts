import { useForm as _useForm } from '@tanstack/react-form'
import { useEffect, useRef } from 'react'
import { useTranslation } from 'react-i18next'

/**
 * Wraps TanStack's useForm with automatic i18n error refresh.
