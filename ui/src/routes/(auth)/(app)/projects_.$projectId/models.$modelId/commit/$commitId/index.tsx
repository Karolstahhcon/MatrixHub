import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'

import { modelCommitQueryOptions } from '@/features/models/models.query'
import { ModelCommitDetailPage } from '@/features/models/pages/ModelCommitDetailPage'

const commitDetailSearchSchema = z.object({
  branch: z.string().optional(),
