import {
  Breadcrumbs,
  Group,
  Text,
  Title,
  rem,
} from '@mantine/core'

import AnchorLink from '@/shared/components/AnchorLink'

import type { TablerIcon } from '@tabler/icons-react'
import type { LinkComponentProps } from '@tanstack/react-router'
import type { ReactNode } from 'react'

export interface AdminPageHeaderItem {
  label: ReactNode
  linkProps?: LinkComponentProps
}

interface AdminPageHeaderProps {
  actions?: ReactNode
  icon: TablerIcon
  items: [AdminPageHeaderItem, ...AdminPageHeaderItem[]]
}

const headerItemStyles = [
  {
    fw: 600,
    fz: rem(18),
    c: 'inherit',
  },
  {
    fw: 600,
    fz: rem(14),
    c: 'inherit',
  },
  {
    fw: 400,
    fz: rem(14),
    c: 'inherit',
  },
] as const

function HeaderItem({
  item,
  index,
