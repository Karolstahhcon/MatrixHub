import dayjs from 'dayjs'

import '@/i18n/dayjs'

type TimestampValue = number | string

export interface ProtobufTimestampLike {
  seconds?: TimestampValue
  nanos?: number
}

export type DateValue = TimestampValue | Date | ProtobufTimestampLike | null | undefined

function isNil(value: unknown): value is null | undefined {
  return value === null || value === undefined
}

function isProtobufTimestamp(value: DateValue): value is ProtobufTimestampLike {
  return typeof value === 'object' && value !== null && ('seconds' in value || 'nanos' in value)
}

export function toMilliseconds(value: TimestampValue) {
  const dateNum = Number(value)

  if (Number.isNaN(dateNum)) {
