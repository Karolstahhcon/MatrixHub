/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../fetch.pb"
import * as GoogleProtobufTimestamp from "../google/protobuf/timestamp.pb"
import * as GoogleProtobufWrappers from "../google/protobuf/wrappers.pb"
import * as MatrixhubV1alpha1Role from "./role.pb"
import * as MatrixhubV1alpha1Utils from "./utils.pb"

export enum ProjectType {
  PROJECT_TYPE_UNSPECIFIED = "PROJECT_TYPE_UNSPECIFIED",
  PROJECT_TYPE_PRIVATE = "PROJECT_TYPE_PRIVATE",
  PROJECT_TYPE_PUBLIC = "PROJECT_TYPE_PUBLIC",
}

export enum ProjectPermissionFilter {
  PERMISSION_FILTER_UNSPECIFIED = "PERMISSION_FILTER_UNSPECIFIED",
  PERMISSION_FILTER_MANAGED_ONLY = "PERMISSION_FILTER_MANAGED_ONLY",
  PERMISSION_FILTER_CAN_WRITE = "PERMISSION_FILTER_CAN_WRITE",
  PERMISSION_FILTER_CAN_READ = "PERMISSION_FILTER_CAN_READ",
}

export enum MemberType {
  MEMBER_TYPE_USER = "MEMBER_TYPE_USER",
  MEMBER_TYPE_GROUP = "MEMBER_TYPE_GROUP",
}

export type CreateProjectRequest = {
  name?: string
  type?: ProjectType
  registryId?: GoogleProtobufWrappers.Int32Value
  organization?: string
}

export type CreateProjectResponse = {
}

export type ListProjectsRequest = {
  name?: string
  type?: ProjectType
  permissionFilter?: ProjectPermissionFilter
  page?: number
  pageSize?: number
}

export type ListProjectsResponse = {
  projects?: Project[]
  pagination?: MatrixhubV1alpha1Utils.Pagination
}

export type GetProjectRequest = {
  name?: string
}

export type GetProjectResponse = {
  name?: string
  type?: ProjectType
  registryUrl?: string
  organization?: string
  modelCount?: number
  datasetCount?: number
  updatedAt?: GoogleProtobufTimestamp.Timestamp
}

export type DeleteProjectRequest = {
  name?: string
}

export type DeleteProjectResponse = {
}

export type UpdateProjectRequest = {
  name?: string
  type?: ProjectType
}

export type UpdateProjectResponse = {
}

export type Project = {
  name?: string
  type?: ProjectType
  enabledRegistry?: boolean
