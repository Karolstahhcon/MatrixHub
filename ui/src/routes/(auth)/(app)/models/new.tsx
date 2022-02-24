import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'

import ModelCreatePage from '@/features/models/pages/ModelCreatePage'

const modelCreateSearchSchema = z.object({
  projectId: z.string().trim().optional().catch(undefined),
