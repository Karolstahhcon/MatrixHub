import { filesize } from 'filesize'

export function formatStorageSize(value: string | undefined) {
  if (!value) {
    return '-'
  }
