import {
  type AddProjectMemberWithRoleRequest,
  Projects,
  type RemoveProjectMembersRequest,
  type UpdateProjectMemberRoleRequest,
} from '@matrixhub/api-ts/v1alpha1/project.pb'
import { mutationOptions } from '@tanstack/react-query'

import { memberKeys } from './members.query'

export function addMemberMutationOptions() {
  return mutationOptions({
    mutationFn: (input: AddProjectMemberWithRoleRequest) =>
      Projects.AddProjectMemberWithRole(input),
    meta: {
      invalidates: [memberKeys.lists()],
    },
  })
}
