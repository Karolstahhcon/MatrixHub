import {
  type ComboboxItemGroup,
  Select,
} from '@mantine/core'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import type { Revision } from '@matrixhub/api-ts/v1alpha1/model.pb.ts'

const EMPTY_REVISIONS: Revision[] = []

export interface RevisionSelectProps {
  revision?: string
  branches?: Revision[]
  tags?: Revision[]
  onRevisionChange?: (revision: string) => void
}

export function RevisionSelect({
  revision,
  branches = EMPTY_REVISIONS,
  tags = EMPTY_REVISIONS,
  onRevisionChange,
}: RevisionSelectProps) {
  const { t } = useTranslation()

  const revisionOptions: ComboboxItemGroup[] = []

  if (branches.length > 0) {
    revisionOptions.push({
      group: t('shared.revisionSelect.group.branch'),
      items: branches?.filter(Boolean)?.map(item => ({
        label: item.name as string,
        value: `branch:${item.name}`,
      })),
    })
  }

  if (tags.length > 0) {
    revisionOptions.push({
      group: t('shared.revisionSelect.group.tag'),
      items: tags?.filter(Boolean)?.map(item => ({
        label: item.name as string,
        value: `tag:${item.name}`,
      })),
