import { Anchor } from '@mantine/core'
import { createLink } from '@tanstack/react-router'
import * as React from 'react'

import type { AnchorProps } from '@mantine/core'
import type { LinkComponent } from '@tanstack/react-router'

type TableCellAnchorProps = Omit<AnchorProps, 'href'>
