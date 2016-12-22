import {
  rem,
  Stack,
  Text,
} from '@mantine/core'

import { parseConventionalCommit } from '@/shared/utils'

interface CommitMessageSectionProps {
  message?: string
}

export function CommitMessageSection({ message }: CommitMessageSectionProps) {
