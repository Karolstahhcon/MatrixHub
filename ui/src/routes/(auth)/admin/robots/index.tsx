import { IconRobot } from '@tabler/icons-react'
import { createFileRoute, stripSearchParams } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

import { AdminPageHeader } from '@/features/admin/components/admin-page-header'
import { AdminPageLayout } from '@/features/admin/components/admin-page-layout'
import { RobotsPage } from '@/features/admin/robots/pages/RobotsPage.tsx'
import { robotsQueryOptions } from '@/features/admin/robots/robot.query'
import {
  DEFAULT_PAGE,
  DEFAULT_PAGE_SIZE,
} from '@/utils/constants.ts'

const defaults = {
  query: '',
  page: DEFAULT_PAGE,
  pageSize: DEFAULT_PAGE_SIZE,
}

const robotsSearchSchema = z.object({
  query: z.string().trim().optional().default(defaults.query),
  page: z.coerce.number().int().positive().optional().default(defaults.page),
  pageSize: z.coerce.number().int().positive().optional().default(defaults.pageSize),
