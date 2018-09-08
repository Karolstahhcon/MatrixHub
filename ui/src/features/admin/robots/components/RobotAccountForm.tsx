import {
  Accordion,
  Box,
  Button,
  Checkbox,
  Group,
  LoadingOverlay,
  MultiSelect,
  NumberInput,
  Radio,
  SimpleGrid,
  Stack,
  Text,
  TextInput,
  Textarea,
  rem,
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import {
  type CreateRobotAccountRequest,
  type GetRobotAccountResponse,
  RobotAccountProjectScope,
  type UpdateRobotAccountRequest,
} from '@matrixhub/api-ts/v1alpha1/robot.pb'
import { IconArrowLeft } from '@tabler/icons-react'
import { useMutation } from '@tanstack/react-query'
import {
  Link,
  useNavigate,
} from '@tanstack/react-router'
import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { useForm } from '@/shared/hooks/useForm'
import { fieldError } from '@/shared/utils/form'

import { RobotTokenInfoModal } from './RobotTokenInfoModal.tsx'
import {
  createRobotAccountMutationOptions,
  updateRobotAccountMutationOptions,
} from '../robot.mutation'
import {
  useRobotPermissions,
  useRobotProjects,
} from '../robot.query'
import {
  ROBOT_DESCRIPTION_MAX_LENGTH,
  type RobotAccountFormValues,
  createRobotAccountFormSchema,
} from '../robot.schema'
import {
  getPlatformPermissionCategories,
  getProjectPermissionCategories,
  type RobotPermissionCategory,
} from '../robots.permissions'

interface RobotAccountFormProps {
  mode: 'create' | 'edit'
  defaultValues: RobotAccountFormValues
  robot?: GetRobotAccountResponse
}

function permissionSelectionState(
  selectedPermissions: string[],
  category: RobotPermissionCategory,
) {
  const permissionValues = category.permissions.map(item => item.value)
  const checkedCount = permissionValues.filter(value => selectedPermissions.includes(value)).length

  return {
    checked: checkedCount === permissionValues.length && checkedCount > 0,
    indeterminate: checkedCount > 0 && checkedCount < permissionValues.length,
  }
}

function nextPermissionsForCategory(
  selectedPermissions: string[],
  category: RobotPermissionCategory,
  checked: boolean,
) {
  const permissionValues = category.permissions.map(item => item.value)

  if (checked) {
    return Array.from(new Set([...selectedPermissions, ...permissionValues]))
  }

  return selectedPermissions.filter(value => !permissionValues.includes(value))
}

export function RobotAccountForm({
  mode,
  defaultValues,
  robot,
}: RobotAccountFormProps) {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const createMutation = useMutation(createRobotAccountMutationOptions())
  const updateMutation = useMutation(updateRobotAccountMutationOptions())
  const {
    data: projects,
  } = useRobotProjects()
  const {
    data: permissionsResponse,
  } = useRobotPermissions()
  const [createdRobotName, setCreatedRobotName] = useState('')
  const [createdToken, setCreatedToken] = useState('')
  const [tokenOpened, tokenHandlers] = useDisclosure(false)

  const projectOptions = useMemo(() => {
    const uniqueNames = Array.from(new Set(
      (projects ?? [])
        .map(project => project.name?.trim() ?? '')
        .filter(Boolean),
    ))

    return uniqueNames.map(name => ({
      value: name,
      label: name,
    }))
  }, [projects])

  const platformPermissionCategories = useMemo(
    () => getPlatformPermissionCategories(permissionsResponse?.systemCategories),
    [permissionsResponse?.systemCategories],
  )
  const projectPermissionCategories = useMemo(
    () => getProjectPermissionCategories(permissionsResponse?.projectCategories),
    [permissionsResponse?.projectCategories],
  )

  const robotAccountFormSchema = createRobotAccountFormSchema(t)
  const getRobotAccountFieldIssue = (
    fieldName: keyof RobotAccountFormValues,
    values: RobotAccountFormValues,
  ) => {
    const result = robotAccountFormSchema.safeParse(values)

    if (result.success) {
      return undefined
    }

    return result.error.issues.find(issue => issue.path[0] === fieldName)?.message
  }

  const form = useForm({
    defaultValues,
    validators: {
      onSubmit: robotAccountFormSchema,
    },
    onSubmit: async ({ value }) => {
      if (mode === 'create') {
        const response = await createMutation.mutateAsync({
          name: value.name.trim(),
          description: value.description.trim(),
          expireDays: value.expiryMode === 'days' ? value.expireDays : undefined,
          platformPermissions: value.platformPermissions,
          projectPermissions: value.projectPermissions,
          projects: value.projectScope === RobotAccountProjectScope.ROBOT_ACCOUNT_PROJECT_SCOPE_SELECTED ? value.projects : [],
          projectScope: value.projectScope,
        } satisfies CreateRobotAccountRequest)

        setCreatedRobotName(value.name.trim())
        setCreatedToken(response.token ?? '')
        tokenHandlers.open()

        return
      }

      await updateMutation.mutateAsync({
        id: robot?.id,
        description: value.description.trim(),
        expireDays: value.expiryMode === 'days' ? value.expireDays : undefined,
        platformPermissions: value.platformPermissions,
        projectPermissions: value.projectPermissions,
        projects: value.projectScope === RobotAccountProjectScope.ROBOT_ACCOUNT_PROJECT_SCOPE_SELECTED ? value.projects : [],
        projectScope: value.projectScope,
      } satisfies UpdateRobotAccountRequest)

      await navigate({
        to: '/admin/robots',
      })
    },
  })

  const isMutating = createMutation.isPending || updateMutation.isPending

  const handleBack = async () => {
    await navigate({
      to: '/admin/robots',
    })
  }

  const handleTokenClose = async () => {
    tokenHandlers.close()
    await navigate({
      to: '/admin/robots',
    })
  }

  return (
    <>
      <Stack gap="lg" py="md">
        <Button
          component={Link}
          to="/admin/robots"
          justify="flex-start"
          leftSection={<IconArrowLeft size={32} />}
          variant="transparent"
          color="text/title"
          fw={600}
          px={0}
          fz={rem(18)}
        >
          {t(mode === 'create'
            ? 'routes.admin.robots.createRobot'
            : 'routes.admin.robots.editRobot')}
        </Button>

        <Box
          component="form"
          pos="relative"
          onSubmit={(event) => {
            event.preventDefault()
            void form.handleSubmit()
          }}
        >
          <LoadingOverlay visible={isMutating} />

          <Stack gap="xl">
            <Stack gap="md" w={350}>
              <Text fw={700}>
                {t('routes.admin.robots.sections.basic')}
              </Text>

              <form.Field
                name="name"
                validators={{
                  onChange: robotAccountFormSchema.shape.name,
                }}
              >
                {field => (
                  <TextInput
                    label={t('routes.admin.robots.fields.name')}
                    required
                    w={350}
                    disabled={mode === 'edit'}
                    value={field.state.value}
                    onChange={event => field.handleChange(event.currentTarget.value)}
                    onBlur={field.handleBlur}
                    error={fieldError(field)}
                  />
                )}
              </form.Field>

              <form.Field
                name="description"
                validators={{
                  onChange: robotAccountFormSchema.shape.description,
                }}
              >
                {field => (
                  <Textarea
                    label={t('routes.admin.robots.fields.description')}
                    autosize
                    minRows={3}
                    maxLength={ROBOT_DESCRIPTION_MAX_LENGTH}
                    value={field.state.value}
                    onChange={event => field.handleChange(event.currentTarget.value)}
                    onBlur={field.handleBlur}
                    error={fieldError(field)}
                    description={t('routes.admin.robots.fields.descriptionCount', {
                      count: field.state.value.length,
                      max: ROBOT_DESCRIPTION_MAX_LENGTH,
                    })}
                  />
                )}
              </form.Field>

              <form.Field name="expiryMode">
                {field => (
                  <Radio.Group
                    label={t('routes.admin.robots.fields.expiry')}
