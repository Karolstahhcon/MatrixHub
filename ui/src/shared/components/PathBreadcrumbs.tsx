import {
  Breadcrumbs,
  Group,
  Text,
} from '@mantine/core'
import { type LinkComponentProps } from '@tanstack/react-router'

import AnchorLink from '@/shared/components/AnchorLink.tsx'
import { CopyValueButton } from '@/shared/components/CopyValueButton.tsx'

export interface ModelPathBreadcrumbsProps {
  name: string
  treePath?: string
  getPathLinkProps: (path: string) => LinkComponentProps
  copyValue?: string
}

export function PathBreadcrumbs({
  name,
  treePath,
  getPathLinkProps,
  copyValue,
}: ModelPathBreadcrumbsProps) {
