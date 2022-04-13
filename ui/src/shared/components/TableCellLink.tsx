import { Anchor } from '@mantine/core'
import { createLink } from '@tanstack/react-router'
import * as React from 'react'

import type { AnchorProps } from '@mantine/core'
import type { LinkComponent } from '@tanstack/react-router'

type TableCellAnchorProps = Omit<AnchorProps, 'href'>

const TableCellAnchor = ({
  ref, ...props
}: TableCellAnchorProps & { ref?: React.RefObject<HTMLAnchorElement | null> }) => {
  return (
    <Anchor
      ref={ref}
      size="sm"
      fw={600}
