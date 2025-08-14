import { useSuspenseQuery } from '@tanstack/react-query'
import { useNavigate } from '@tanstack/react-router'

import { modelRevisionsQueryOptions } from '@/features/models/models.query.ts'
import { RevisionSelect } from '@/shared/components/RevisionSelect.tsx'

interface ModelRevisionSelectProps {
  projectId: string
