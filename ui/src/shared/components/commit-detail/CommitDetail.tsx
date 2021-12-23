import {
  Box,
  Group,
  Loader,
  Paper,
  Stack,
  Text,
  Tooltip,
} from '@mantine/core'
import dayjs from 'dayjs'
import {
  type MouseEvent,
  useEffect,
  useState,
} from 'react'
import { useTranslation } from 'react-i18next'

import {
  MAX_FILES,
  renderCommitDiff,
  splitCommitDiff,
} from '@/shared/components/commit-detail/commitDiff'
import { CommitHashCard } from '@/shared/components/commit-detail/CommitHashCard'
import { CommitMessageSection } from '@/shared/components/commit-detail/CommitMessageSection'
import { DiffLimitAlert } from '@/shared/components/commit-detail/DiffLimitAlert'
import { formatDateTime } from '@/shared/utils/date'

import 'diff2html/bundles/css/diff2html.min.css'
import './CommitDetail.css'

import type { Commit } from '@matrixhub/api-ts/v1alpha1/model.pb'

interface CommitDetailProps {
  commit?: Commit
}

const DIFF_PARSE_DELAY_MS = 30
const RAW_DIFF_URL_REVOKE_DELAY_MS = 60_000

interface DiffRenderState {
  html: string
  hasTooManyFiles: boolean
  loading: boolean
}

const EMPTY_DIFF_RENDER_STATE: DiffRenderState = {
  html: '',
  hasTooManyFiles: false,
  loading: false,
}

function EmptyDiffPrompt() {
  const { t } = useTranslation()

  return (
    <Text size="sm" c="dimmed">
      {t('shared.commitDetail.emptyDiff')}
    </Text>
  )
}

export function CommitDetail({ commit }: CommitDetailProps) {
  const { t } = useTranslation()
  const hasRawDiff = Boolean(commit?.diff)
  const authorDate = formatDateTime(commit?.authorDate, 'YYYY-MM-DD HH:mm:ss')

  const [diffRender, setDiffRender] = useState<DiffRenderState>(EMPTY_DIFF_RENDER_STATE)

  useEffect(() => {
    const rawDiff = commit?.diff ?? ''
    let cancelled = false

    if (!rawDiff) {
      const resetTimerId = window.setTimeout(() => {
        if (!cancelled) {
          setDiffRender(EMPTY_DIFF_RENDER_STATE)
        }
      }, 0)

      return () => {
        cancelled = true
        window.clearTimeout(resetTimerId)
      }
    }

    const loadingTimerId = window.setTimeout(() => {
      if (cancelled) {
        return
      }

      setDiffRender(prev => ({
        ...prev,
        loading: true,
      }))
    }, 0)

    const parseTimerId = window.setTimeout(() => {
      if (cancelled) {
        return
      }

      const nextDiffRender = renderCommitDiff(
        rawDiff,
        t('shared.commitDetail.diffTooBig'),
        t('shared.commitDetail.rawDiffLink'),
      )

      if (!cancelled) {
        setDiffRender({
          ...nextDiffRender,
          loading: false,
        })
      }
    }, DIFF_PARSE_DELAY_MS)

