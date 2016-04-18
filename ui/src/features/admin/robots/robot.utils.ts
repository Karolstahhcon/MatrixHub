import {
  type GetRobotAccountResponse,
  RobotAccountStatus,
} from '@matrixhub/api-ts/v1alpha1/robot.pb'

export type RobotAccount = GetRobotAccountResponse

export function getRobotRowId(robot: RobotAccount) {
  return String(robot.id ?? robot.name ?? '-')
