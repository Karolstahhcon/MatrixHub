import {
  Alert,
  PasswordInput,
  Stack,
  Switch,
  rem,
} from '@mantine/core'
import { IconInfoCircle } from '@tabler/icons-react'
import { useMutation } from '@tanstack/react-query'
import { useEffect, useEffectEvent } from 'react'
import { useTranslation } from 'react-i18next'

import { ModalWrapper } from '@/shared/components/ModalWrapper.tsx'
import { useForm } from '@/shared/hooks/useForm.ts'
import { fieldError } from '@/shared/utils/form.ts'

import { refreshRobotAccountTokenMutationOptions } from '../robot.mutation.ts'
import {
  ROBOT_TOKEN_MAX_LENGTH,
  ROBOT_TOKEN_MIN_LENGTH,
  refreshRobotTokenFormDefaults,
  refreshRobotTokenSchema,
  type RefreshRobotTokenFormValues,
} from '../robot.schema.ts'

import type { RobotAccount } from '../robot.utils.ts'
import type { StandardSchemaV1Issue } from '@tanstack/react-form'

interface RefreshRobotTokenModalProps {
  robot: RobotAccount
  opened: boolean
  onClose: () => void
  onSuccess: (token: string) => void
}

function getRefreshRobotTokenFieldIssues(
  fieldName: keyof Pick<RefreshRobotTokenFormValues, 'token' | 'confirmToken'>,
  values: RefreshRobotTokenFormValues,
) {
  const result = refreshRobotTokenSchema.safeParse(values)

