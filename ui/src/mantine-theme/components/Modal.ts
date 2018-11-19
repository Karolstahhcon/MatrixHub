import { Modal, rem } from '@mantine/core'

const modalSizeMap = {
  xs: rem(300),
  sm: rem(390),
  md: rem(480),
} as const

const isMappedModalSize = (size: unknown): size is keyof typeof modalSizeMap => (
  typeof size === 'string' && Object.prototype.hasOwnProperty.call(modalSizeMap, size)
)

