import { createFileRoute, stripSearchParams } from '@tanstack/react-router'
import { z } from 'zod'

import { projectModelsQueryOptions } from '@/features/models/models.query'
import { ProjectModelsPage } from '@/features/projects/pages/ProjectModelsPage'

// -- URL search schema (route concern) --

const defaults = {
  q: '',
  sort: 'updatedAt' as const,
  order: 'desc' as const,
