import {
  Button, Center, Stack, Text,
} from '@mantine/core'
import { Link, type LinkComponentProps } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'

interface PathNotFoundProps {
  path: string
  branch: string
  rootLinkProps: LinkComponentProps
}

export function PathNotFound({
  path,
  branch,
  rootLinkProps,
}: PathNotFoundProps) {
  const { t } = useTranslation()
