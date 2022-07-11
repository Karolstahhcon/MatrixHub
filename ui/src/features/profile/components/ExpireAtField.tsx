import {
  Group,
  Radio,
} from '@mantine/core'
import { DatePickerInput } from '@mantine/dates'
import dayjs from 'dayjs'
import { useTranslation } from 'react-i18next'

const defaultExpireAt = () => dayjs().add(1, 'day').format('YYYY-MM-DD')

interface ExpireAtFieldProps {
  value: string
  onChange: (value: string) => void
  error?: string
}

