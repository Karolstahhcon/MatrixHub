import { z } from 'zod'

import i18n from '@/i18n'

const NAME_REGEX = /^[A-Za-z0-9](?:[A-Za-z0-9-]*[A-Za-z0-9])$/

// ---------------------------------------------------------------------------
// Field-level schemas — used on each <form.Field validators={{ onBlur }}>
// ---------------------------------------------------------------------------

export const projectNameSchema = z
  .string()
  .trim()
  .superRefine((val, ctx) => {
    if (!val) {
      ctx.addIssue({
        code: 'custom',
        message: i18n.t('projects.validation.nameRequired'),
      })

      return
    }

    if (val.length < 2) {
      ctx.addIssue({
        code: 'custom',
        message: i18n.t('projects.validation.nameMinLength'),
      })

