/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../fetch.pb"
import * as MatrixhubV1alpha1Role from "./role.pb"

export enum AccessTokenStatus {
  ACCESS_TOKEN_STATUS_UNKNOWN = "ACCESS_TOKEN_STATUS_UNKNOWN",
  ACCESS_TOKEN_STATUS_VALID = "ACCESS_TOKEN_STATUS_VALID",
  ACCESS_TOKEN_STATUS_EXPIRED = "ACCESS_TOKEN_STATUS_EXPIRED",
}

export enum SSHKeyStatus {
  SSH_KEY_STATUS_UNKNOWN = "SSH_KEY_STATUS_UNKNOWN",
  SSH_KEY_STATUS_VALID = "SSH_KEY_STATUS_VALID",
  SSH_KEY_STATUS_EXPIRED = "SSH_KEY_STATUS_EXPIRED",
}

export type ResetPasswordRequest = {
  oldPassword?: string
  newPassword?: string
}

export type ResetPasswordResponse = {
}

export type ListAccessTokensRequest = {
}

export type ListAccessTokensResponse = {
  items?: AccessToken[]
}

export type AccessToken = {
  id?: number
  name?: string
  status?: AccessTokenStatus
  createdAt?: string
  expiredAt?: string
}

export type CreateAccessTokenRequest = {
  name?: string
  expireAt?: string
}

export type CreateAccessTokenResponse = {
  token?: string
}

export type DeleteAccessTokenRequest = {
  id?: number
}

export type DeleteAccessTokenResponse = {
}

export type GetProjectRolesRequest = {
}

export type GetProjectRolesResponse = {
  projectRoles?: {[key: string]: MatrixhubV1alpha1Role.ProjectRoleType}
}

export type GetCurrentUserRequest = {
}

export type GetCurrentUserResponse = {
  id?: number
  username?: string
  isAdmin?: boolean
}

export type CreateSSHKeyRequest = {
  name?: string
  publicKey?: string
  expireAt?: string
}

export type CreateSSHKeyResponse = {
}

export type DeleteSSHKeyRequest = {
  id?: number
}

export type DeleteSSHKeyResponse = {
