import { createFileRoute, stripSearchParams } from '@tanstack/react-router'

import {
  membersQueryOptions,
  membersSearchDefaults,
  membersSearchSchema,
} from '@/features/projects/members/members.query'
import { ProjectMembersPage } from '@/features/projects/members/pages/ProjectMembersPage'

// -- Route definition --
