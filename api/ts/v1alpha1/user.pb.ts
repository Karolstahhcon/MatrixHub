/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../fetch.pb"
import * as GoogleProtobufTimestamp from "../google/protobuf/timestamp.pb"
import * as MatrixhubV1alpha1Utils from "./utils.pb"

export enum UserSource {
  USER_SOURCE_UNSPECIFIED = "USER_SOURCE_UNSPECIFIED",
  USER_SOURCE_LOCAL = "USER_SOURCE_LOCAL",
}

export type User = {
  id?: number
  username?: string
  email?: string
  isAdmin?: boolean
  source?: UserSource
  createdAt?: GoogleProtobufTimestamp.Timestamp
}

export type CreateUserRequest = {
