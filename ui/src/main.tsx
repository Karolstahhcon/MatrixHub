import { MantineProvider } from '@mantine/core'
import { Notifications } from '@mantine/notifications'
import { QueryClientProvider } from '@tanstack/react-query'
import { RouterProvider } from '@tanstack/react-router'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import { installApiRequestHeaders } from '@/shared/utils/apiRequestHeaders'

import './i18n/index.ts'
import '@mantine/core/styles.css'
import '@mantine/notifications/styles.css'
