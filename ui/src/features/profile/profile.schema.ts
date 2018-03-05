import dayjs from 'dayjs'
import z from 'zod'

import type { TFunction } from 'i18next'

export function expireAtSchema(t: TFunction) {
  return z.string().refine((value) => {
    if (!value) {
