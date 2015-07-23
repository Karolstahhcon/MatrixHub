import { createFileRoute, stripSearchParams } from '@tanstack/react-router'
import { z } from 'zod'

import { modelCommitsQueryOptions } from '@/features/models/models.query'
import { ModelCommitsPage } from '@/features/models/pages/ModelCommitsPage'

const defaults = { page: 1 }

const commitsSearchSchema = z.object({
  page: z.coerce.number().int().positive().optional().default(defaults.page).catch(defaults.page),
})
