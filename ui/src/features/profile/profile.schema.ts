import dayjs from 'dayjs'
import z from 'zod'

import type { TFunction } from 'i18next'

export function expireAtSchema(t: TFunction) {
  return z.string().refine((value) => {
    if (!value) {
      return true
    }

    return !dayjs(value).startOf('day').isBefore(dayjs().startOf('day'))
  }, { error: t('profile.expireTimeError') })
}

export function createSshKeySchema(t: TFunction) {
  return z.object({
