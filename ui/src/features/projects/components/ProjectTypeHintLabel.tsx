import {
  Anchor, Group, List, Popover, Text,
} from '@mantine/core'
import { IconExternalLink } from '@tabler/icons-react'
import { useState } from 'react'
import { Trans, useTranslation } from 'react-i18next'

import IconQuestion from '@/assets/svgs/question.svg?react'

import type { KeyboardEvent, ReactNode } from 'react'

// TODO: point this at the permission guide once that doc is published. It
// currently falls back to the project members doc, which covers role
// permissions but not the full ruleset.
const PROJECT_PERMISSION_DOC_URL = '/docs/operations/project-management/members/'

const DROPDOWN_WIDTH = 360

/**
 * Anchor that appends an external-link icon after the translated label.
 *
 * The icon cannot be an empty placeholder tag in the translation string:
 * Trans clones the mapped element with the tag's own children, so a
 * `<1></1>` placeholder renders as nothing.
 */
function PermissionDocLink({ children }: { children?: ReactNode }) {
  const {
    t,
    i18n,
