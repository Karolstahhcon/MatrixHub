import { createFileRoute, redirect } from '@tanstack/react-router'

import { currentUserQueryOptions } from '@/features/auth/auth.query'
import { LoginPage } from '@/features/login/pages/LoginPage'
import { getErrorMessage, queryClient } from '@/queryClient'

const ALREADY_LOGGED_IN_ERROR_MESSAGE = 'ALREADY_LOGGED_IN'
