import {
  createRootRouteWithContext, Outlet, HeadContent,
} from '@tanstack/react-router'
import { lazy, Suspense } from 'react'

import i18n from '@/i18n'

import type { QueryClient } from '@tanstack/react-query'

const TanStackDevtools = import.meta.env.DEV
  ? lazy(() =>
      Promise.all([
        import('@tanstack/react-devtools'),
        import('@tanstack/react-query-devtools'),
        import('@tanstack/react-router-devtools'),
        import('@tanstack/react-form-devtools'),
      ]).then(([devtools, queryDevtools, routerDevtools, formDevtools]) => ({
        default: () => (
          <devtools.TanStackDevtools
            plugins={[
              {
                name: 'TanStack Query',
                render: <queryDevtools.ReactQueryDevtoolsPanel />,
              },
              {
