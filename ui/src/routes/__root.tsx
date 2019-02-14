import {
  createRootRouteWithContext, Outlet, HeadContent,
} from '@tanstack/react-router'
import { lazy, Suspense } from 'react'

import i18n from '@/i18n'

import type { QueryClient } from '@tanstack/react-query'

const TanStackDevtools = import.meta.env.DEV
  ? lazy(() =>
      Promise.all([
