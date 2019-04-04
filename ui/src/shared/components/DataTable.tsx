import {
  ActionIcon,
  type BoxProps,
  Button,
  Center,
  type PaperProps,
  Stack,
  type TableProps as MantineTableProps,
  Text,
  noop,
  Box,
} from '@mantine/core'
import { useDebouncedCallback, useDebouncedValue } from '@mantine/hooks'
import { IconRefresh, IconTrash } from '@tabler/icons-react'
import { MantineReactTable } from 'mantine-react-table'
import 'mantine-react-table/styles.css'
import {
  useEffect,
  useMemo,
  useRef,
} from 'react'
import { useTranslation } from 'react-i18next'

import classes from './DataTable.module.css'
import { Pagination } from './Pagination'
import {
  SearchToolbar,
  type SearchToolbarProps,
} from './SearchToolbar'

import type { Pagination as PaginationData } from '@matrixhub/api-ts/v1alpha1/utils.pb'
import type {
  HTMLPropsRef,
  MRT_ColumnDef,
  MRT_Row,
  MRT_RowData,
  MRT_RowSelectionState,
  MRT_TableInstance,
  MRT_TableOptions,
} from 'mantine-react-table'
import type {
  Dispatch, MutableRefObject, ReactNode, SetStateAction,
} from 'react'

// -- Toolbar types --

export interface DataTableToolbarProps {
  /** Show search input. Provide a placeholder string or `true` to use default. */
  searchPlaceholder?: string | boolean
  /** Controlled search value. */
  searchValue?: string
  /** Called when search input changes. */
  onSearchChange?: (value: string) => void
  /** Extra props forwarded to the default SearchToolbar. */
  searchToolbarProps?: Omit<
    SearchToolbarProps,
    'searchPlaceholder' | 'searchValue' | 'onSearchChange' | 'children'
  >

  /** Show refresh button when provided. */
  onRefresh?: () => void

  /** Number of selected items. Shows batch-delete button when > 0 and `onBatchDelete` is provided. */
  selectedCount?: number
  /** Called when batch-delete button is clicked. */
  onBatchDelete?: () => void

  /** Slot: replaces the entire toolbar row. Receives default toolbar as children for composition. */
  renderToolbar?: (defaultToolbar: ReactNode) => ReactNode
  /** Slot: extra actions rendered after built-in buttons. */
  toolbarExtra?: ReactNode
}

export type DataTableRowActionsProps<TData extends MRT_RowData>
  = Parameters<NonNullable<MRT_TableOptions<TData>['renderRowActions']>>[0]

/**
 * Common props for resource table wrappers (e.g. ProjectsTable, ModelsTable).
 * Provides the standard set of pagination, search, selection, and action props
 * so each resource table doesn't have to re-declare them.
 */

export interface DataTableProps<TData extends MRT_RowData> extends DataTableToolbarProps {
  /** Row data array */
  data: TData[]
  /** Column definitions */
  columns: MRT_ColumnDef<TData>[]

  // --- Pagination ---
  pagination?: PaginationData
  page?: number
  onPageChange?: (page: number) => void

  // --- Empty state ---
  emptyTitle?: ReactNode
  emptyDescription?: ReactNode

  // --- Row selection ---
  enableRowSelection?: boolean | ((row: MRT_Row<TData>) => boolean)
  enableSelectAll?: boolean
  rowSelection?: MRT_RowSelectionState
  onRowSelectionChange?: Dispatch<SetStateAction<MRT_RowSelectionState>>
  getRowId?: (row: TData, index: number) => string

  // --- Row actions ---
  enableRowActions?: boolean
  renderRowActions?: MRT_TableOptions<TData>['renderRowActions']
  positionActionsColumn?: 'first' | 'last'
  pinRowActions?: boolean

  // --- Loading ---
  loading?: boolean
  fetching?: boolean

  // --- Empty rows fallback ---
  renderEmptyRowsFallback?: MRT_TableOptions<TData>['renderEmptyRowsFallback']

  // --- Header ---
  /** Hide the table header row. Defaults to false. */
  hideTableHead?: boolean

  // --- Display column overrides ---
  displayColumnDefOptions?: MRT_TableOptions<TData>['displayColumnDefOptions']

  // --- Escape hatch ---
  tableOptions?: Omit<
    Partial<MRT_TableOptions<TData>>,
    | 'columns'
    | 'data'
    | 'enableRowSelection'
    | 'enableSelectAll'
    | 'onRowSelectionChange'
    | 'getRowId'
    | 'enableRowActions'
    | 'renderRowActions'
    | 'positionActionsColumn'
    | 'displayColumnDefOptions'
  >
}

// -- Internal helpers --

const ROW_ACTIONS_COLUMN_ID = 'mrt-row-actions'

function hasContent(value: ReactNode) {
  return value !== null && value !== undefined && value !== false && value !== ''
}

function mergeTableOptionProps<TData extends MRT_RowData, TProps extends object>(
  defaults: TProps,
  props:
    | TProps
    | ((args: { table: MRT_TableInstance<TData> }) => TProps)
    | undefined,
) {
  if (!props) {
    return defaults
  }

  if (typeof props === 'function') {
    return (args: { table: MRT_TableInstance<TData> }) => ({
      ...defaults,
      ...props(args),
    })
  }

  return {
    ...defaults,
    ...props,
  }
}

type TableContainerProps = BoxProps & HTMLPropsRef<HTMLDivElement>

// The wrapper's own `ref` and `className` drive the pinned-shadow tracking and
// its stylesheet, so they must survive a caller-supplied
// `mantineTableContainerProps`. Merge those two explicitly instead of letting a
// plain spread drop them; everything else still overrides as usual.
//
// MRT types this `ref` as a ref object rather than a callback, so the two refs
// cannot be combined into one. The wrapper's ref stays on the element and the
// caller's is populated from it in an effect.
function mergeTableContainerProps<TData extends MRT_RowData>(
  defaults: TableContainerProps & { className: string },
  props:
    | TableContainerProps
    | ((args: { table: MRT_TableInstance<TData> }) => TableContainerProps)
    | undefined,
) {
  const combine = (overrides: TableContainerProps | undefined): TableContainerProps => {
    if (!overrides) {
      return defaults
    }

    const {
      ref: _callerRef,
      className: overrideClassName,
      ...rest
    } = overrides

    return {
      ...defaults,
      ...rest,
      ref: defaults.ref,
      className: overrideClassName
        ? `${defaults.className} ${overrideClassName}`
        : defaults.className,
    }
  }

  if (typeof props === 'function') {
    return (args: { table: MRT_TableInstance<TData> }) => combine(props(args))
  }

  return combine(props)
}

function resolveTableOptionProps<TArgs extends object, TProps extends object>(
  props: TProps | ((args: TArgs) => TProps) | undefined, args: TArgs,
) {
  if (!props) {
    return undefined
  }

  return typeof props === 'function'
    ? props(args)
    : props
}

// Build a value→label map from the static `data` of a select-filter column so
// the active-filter tooltip shows the human-readable label instead of the raw
// value. Skipped when `mantineFilterSelectProps` is a function (we can't inspect
// it without table args) or `data` is missing/grouped.
function deriveSelectFilterLabelMap<TData extends MRT_RowData>(
  column: MRT_ColumnDef<TData>,
): Map<string, string> | undefined {
  const filterProps = column.mantineFilterSelectProps

  if (!filterProps || typeof filterProps === 'function') {
    return undefined
  }

  const data = (filterProps as { data?: unknown }).data

  if (!Array.isArray(data) || data.length === 0) {
    return undefined
  }

  const labelByValue = new Map<string, string>()

  for (const item of data) {
    if (typeof item === 'string') {
      labelByValue.set(item, item)
    } else if (item && typeof item === 'object' && 'value' in item) {
      const value = String((item as { value: unknown }).value ?? '')
      const label = (item as { label?: unknown }).label

      labelByValue.set(value, typeof label === 'string' ? label : value)
    }
  }

  return labelByValue.size > 0 ? labelByValue : undefined
}

function withFilterTooltipLabels<TData extends MRT_RowData>(
  columns: MRT_ColumnDef<TData>[],
): MRT_ColumnDef<TData>[] {
  return columns.map((column) => {
    if (column.filterTooltipValueFn) {
      return column
    }
    const labelByValue = deriveSelectFilterLabelMap(column)

    if (!labelByValue) {
      return column
    }

    return {
      ...column,
      filterTooltipValueFn: (value: unknown) => {
        const key = String(value ?? '')

        return labelByValue.get(key) ?? key
      },
    }
  })
}

// MRT always paints the pinned-column divider, even when the table fits and
// there is nothing to scroll. Track the container's real scroll position and
// expose it as data attributes so the stylesheet can hide each side's shadow
// until that side actually has content scrolled out of view.
function usePinnedShadowContainerRef(callerRef?: MutableRefObject<HTMLDivElement | null> | null) {
  const containerRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    const element = containerRef.current

    // MRT only accepts one ref object on the container, so a caller-supplied
    // ref is populated from ours rather than replacing it.
    if (callerRef) {
      callerRef.current = element
    }

    if (!element) {
      return
    }

    const detachCallerRef = () => {
      if (callerRef) {
        callerRef.current = null
      }
    }

    const update = () => {
      const {
        scrollLeft,
        scrollWidth,
        clientWidth,
      } = element

      // 1px tolerance: sub-pixel layout can leave a fractional remainder.
      element.toggleAttribute('data-scrolled-left', scrollLeft > 0)
      element.toggleAttribute(
        'data-scrollable-right',
        scrollWidth - clientWidth - scrollLeft > 1,
      )
    }

    update()

    const observer = new ResizeObserver(update)

    // The container tracks viewport resizes; the inner table tracks row and
    // column changes, which move the scrollable width without resizing the
    // container.
    observer.observe(element)

    const table = element.querySelector('table')

    if (table) {
      observer.observe(table)
    }

    element.addEventListener('scroll', update, { passive: true })

    return () => {
