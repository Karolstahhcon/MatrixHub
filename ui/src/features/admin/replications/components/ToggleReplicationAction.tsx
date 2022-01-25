import { Button } from '@mantine/core'
import { useMutation } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { switchReplicationMutationOptions } from '../replications.mutation'

import type { SyncPolicyItem } from '@matrixhub/api-ts/v1alpha1/sync_policy.pb'

interface ToggleReplicationActionProps {
