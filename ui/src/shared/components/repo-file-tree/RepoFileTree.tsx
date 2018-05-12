import { Text } from '@mantine/core'
import { FileType } from '@matrixhub/api-ts/v1alpha1/model.pb'
import {
  IconFile,
  IconFileDescription,
  IconFileSettings,
  IconFolder,
} from '@tabler/icons-react'
import { Link, type LinkComponentProps } from '@tanstack/react-router'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import AnchorLink from '@/shared/components/AnchorLink'
import { DataTable, type DataTableRowActionsProps } from '@/shared/components/DataTable'
import { formatRelativeTime } from '@/shared/utils/date'
import { formatStorageSize } from '@/shared/utils/format'

import classes from './RepoFileTree.module.css'

import type { File as RepoFile } from '@matrixhub/api-ts/v1alpha1/model.pb'
import type { MRT_ColumnDef } from 'mantine-react-table'

type LinkPropsBuilder<T> = (value: T) => LinkComponentProps

interface RepoFileTreeMeta {
  buildTreeLink: LinkPropsBuilder<string>
  buildBlobLink: LinkPropsBuilder<string>
}

interface RepoFileTreeProps {
  files: RepoFile[]
  latestCommit?: RepoFile['commit']
  buildTreeLink: LinkPropsBuilder<string>
  buildBlobLink: LinkPropsBuilder<string>
  buildCommitLink?: LinkPropsBuilder<string>
  isLoading?: boolean
}

type FileCellProps = Parameters<NonNullable<MRT_ColumnDef<RepoFile>['Cell']>>[0]

function FileNameCell({
  row, table,
}: FileCellProps) {
  const file = row.original
  const isDir = file.type === FileType.DIR
  const filePath = file.path ?? file.name ?? ''
  const fileName = file.name ?? ''
  const {
    buildTreeLink, buildBlobLink,
  } = table.options.meta as RepoFileTreeMeta
  const linkProps = isDir ? buildTreeLink(filePath) : buildBlobLink(filePath)

  return (
    <div className={classes.nameCell}>
      <span className={classes.nameIcon}>
        <FileIcon file={file} />
      </span>
      <Link className={classes.fileLink} {...linkProps}>
        {fileName}
      </Link>
    </div>
  )
}

