import { createFileRoute, stripSearchParams } from '@tanstack/react-router'
import { z } from 'zod'

import { catalogModelsQueryOptions } from '@/features/models/models.query.ts'
import { ModelsPage } from '@/features/models/pages/ModelsPage'

const defaults = {
  q: '',
  sort: 'updatedAt' as const,
  order: 'desc' as const,
  page: 1,
}
