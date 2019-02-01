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
      key,
      name: category.name ?? '',
      permissions: (category.permissions ?? []).map((permission) => {
        const value = permission.permission ?? permission.name ?? ''

        return {
          value,
          name: permission.name ?? value,
        }
      }).filter(item => item.value),
