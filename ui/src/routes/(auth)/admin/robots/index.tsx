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
