import { Anchor } from '@mantine/core'
import { createLink } from '@tanstack/react-router'
import * as React from 'react'

import type { AnchorProps } from '@mantine/core'
import type { LinkComponent } from '@tanstack/react-router'

type MartineAnchorProps = Omit<AnchorProps, 'href'>

const MantineLinkComponent = ({
  ref, ...props
}: MartineAnchorProps & { ref?: React.RefObject<HTMLAnchorElement | null> }) => {
  return <Anchor ref={ref} {...props} />
}
