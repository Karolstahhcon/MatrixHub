import { useClipboard } from '@/shared/hooks/useClipboard'

import type { ReactNode } from 'react'

export interface CopyButtonProps {
  /** Children callback that receives the current status and copy function. */
  children: (payload: {
    copied: boolean
    copy: () => void
