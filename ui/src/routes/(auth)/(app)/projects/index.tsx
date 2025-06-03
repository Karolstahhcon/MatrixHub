import { createFileRoute, stripSearchParams } from '@tanstack/react-router'
import { z } from 'zod'

import { ProjectsPage } from '@/features/projects/pages/ProjectsPage'
import { projectsQueryOptions } from '@/features/projects/projects.query'

const DEFAULT_PROJECTS_PAGE = 1

const defaults = {
  page: DEFAULT_PROJECTS_PAGE,
  query: '',
