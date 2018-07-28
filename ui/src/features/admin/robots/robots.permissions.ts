import type { RoleCategory } from '@matrixhub/api-ts/v1alpha1/role.pb'

export interface RobotPermissionItem {
  value: string
  name: string
}

export interface RobotPermissionCategory {
  key: string
  name: string
  permissions: RobotPermissionItem[]
}

function transformPermissionCategories(
  categories: RoleCategory[] | undefined,
): RobotPermissionCategory[] {
  return (categories ?? []).map((category) => {
    const key = category.name ?? ''

    return {
