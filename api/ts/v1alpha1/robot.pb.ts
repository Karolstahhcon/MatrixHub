/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../fetch.pb"
import * as MatrixhubV1alpha1Utils from "./utils.pb"

export enum RobotAccountProjectScope {
  ROBOT_ACCOUNT_PROJECT_SCOPE_SELECTED = "ROBOT_ACCOUNT_PROJECT_SCOPE_SELECTED",
  ROBOT_ACCOUNT_PROJECT_SCOPE_ALL = "ROBOT_ACCOUNT_PROJECT_SCOPE_ALL",
}

export enum RobotAccountStatus {
  ROBOT_ACCOUNT_STATUS_UNSPECIFIED = "ROBOT_ACCOUNT_STATUS_UNSPECIFIED",
  ROBOT_ACCOUNT_STATUS_DISABLED = "ROBOT_ACCOUNT_STATUS_DISABLED",
  ROBOT_ACCOUNT_STATUS_ENABLED = "ROBOT_ACCOUNT_STATUS_ENABLED",
}

export enum RobotAccountExpireStatus {
  ROBOT_ACCOUNT_EXPIRE_STATUS_EXPIRED = "ROBOT_ACCOUNT_EXPIRE_STATUS_EXPIRED",
  ROBOT_ACCOUNT_EXPIRE_STATUS_VALID = "ROBOT_ACCOUNT_EXPIRE_STATUS_VALID",
  ROBOT_ACCOUNT_EXPIRE_STATUS_NEVER = "ROBOT_ACCOUNT_EXPIRE_STATUS_NEVER",
}

export type CreateRobotAccountRequest = {
  name?: string
  description?: string
  expireDays?: number
