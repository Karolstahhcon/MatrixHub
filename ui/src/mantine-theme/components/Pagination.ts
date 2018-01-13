import { Pagination } from '@mantine/core'

export const paginationTheme = Pagination.extend({
  defaultProps: {
    boundaries: 1,
    siblings: 2,
    color: 'cyan',
    size: 20,
    radius: 4,
    gap: 8,
