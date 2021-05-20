/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../fetch.pb"
import * as MatrixhubV1alpha1Registry from "./registry.pb"
import * as MatrixhubV1alpha1Utils from "./utils.pb"

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };
type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (
    keyof T extends infer K ?
      (K extends string & keyof T ? { [k in K]: T[K] } & Absent<T, K>
        : never)
    : never);

export enum SyncPolicyType {
  SYNC_POLICY_TYPE_UNSPECIFIED = "SYNC_POLICY_TYPE_UNSPECIFIED",
  SYNC_POLICY_TYPE_PULL_BASE = "SYNC_POLICY_TYPE_PULL_BASE",
  SYNC_POLICY_TYPE_PUSH_BASE = "SYNC_POLICY_TYPE_PUSH_BASE",
}

export enum ResourceType {
  RESOURCE_TYPE_UNSPECIFIED = "RESOURCE_TYPE_UNSPECIFIED",
  RESOURCE_TYPE_MODEL = "RESOURCE_TYPE_MODEL",
  RESOURCE_TYPE_DATASET = "RESOURCE_TYPE_DATASET",
}

export enum TriggerType {
  TRIGGER_TYPE_UNSPECIFIED = "TRIGGER_TYPE_UNSPECIFIED",
  TRIGGER_TYPE_MANUAL = "TRIGGER_TYPE_MANUAL",
  TRIGGER_TYPE_SCHEDULED = "TRIGGER_TYPE_SCHEDULED",
}

export enum SyncTaskStatus {
  SYNC_TASK_STATUS_UNSPECIFIED = "SYNC_TASK_STATUS_UNSPECIFIED",
  SYNC_TASK_STATUS_RUNNING = "SYNC_TASK_STATUS_RUNNING",
  SYNC_TASK_STATUS_SUCCEEDED = "SYNC_TASK_STATUS_SUCCEEDED",
  SYNC_TASK_STATUS_FAILED = "SYNC_TASK_STATUS_FAILED",
  SYNC_TASK_STATUS_STOPPED = "SYNC_TASK_STATUS_STOPPED",
  SYNC_TASK_STATUS_PENDING = "SYNC_TASK_STATUS_PENDING",
}

export enum SyncJobStatus {
  SYNC_JOB_STATUS_UNSPECIFIED = "SYNC_JOB_STATUS_UNSPECIFIED",
  SYNC_JOB_STATUS_RUNNING = "SYNC_JOB_STATUS_RUNNING",
  SYNC_JOB_STATUS_SUCCEEDED = "SYNC_JOB_STATUS_SUCCEEDED",
  SYNC_JOB_STATUS_FAILED = "SYNC_JOB_STATUS_FAILED",
  SYNC_JOB_STATUS_STOPPED = "SYNC_JOB_STATUS_STOPPED",
  SYNC_JOB_STATUS_PENDING = "SYNC_JOB_STATUS_PENDING",
}

export type PullBasePolicy = {
  sourceRegistryId?: number
  resourceName?: string
  resourceTypes?: ResourceType[]
  sourceRegistry?: MatrixhubV1alpha1Registry.Registry
  targetProjectName?: string
}

export type PushBasePolicy = {
  resourceName?: string
  resourceTypes?: ResourceType[]
  targetRegistryId?: number
  targetRegistry?: MatrixhubV1alpha1Registry.Registry
  targetProjectName?: string
}


type BaseSyncPolicyItem = {
  id?: number
  name?: string
  description?: string
  policyType?: SyncPolicyType
  triggerType?: TriggerType
  bandwidth?: string
  isOverwrite?: boolean
  isDisabled?: boolean
  triggerTypeSchedule?: TriggerTypeSchedule
}

export type SyncPolicyItem = BaseSyncPolicyItem
  & OneOf<{ pullBasePolicy: PullBasePolicy; pushBasePolicy: PushBasePolicy }>


type BaseCreateSyncPolicyRequest = {
  name?: string
  description?: string
  policyType?: SyncPolicyType
  triggerType?: TriggerType
  bandwidth?: string
  isOverwrite?: boolean
  triggerTypeSchedule?: TriggerTypeSchedule
}

export type CreateSyncPolicyRequest = BaseCreateSyncPolicyRequest
  & OneOf<{ pullBasePolicy: PullBasePolicy; pushBasePolicy: PushBasePolicy }>

export type TriggerTypeSchedule = {
  cron?: string
}

export type CreateSyncPolicyResponse = {
  syncPolicy?: SyncPolicyItem
}


type BaseUpdateSyncPolicyRequest = {
  syncPolicyId?: number
  name?: string
  description?: string
  triggerType?: TriggerType
  bandwidth?: string
  isOverwrite?: boolean
  isDisabled?: boolean
  triggerTypeSchedule?: TriggerTypeSchedule
}

export type UpdateSyncPolicyRequest = BaseUpdateSyncPolicyRequest
  & OneOf<{ pullBasePolicy: PullBasePolicy; pushBasePolicy: PushBasePolicy }>

export type UpdateSyncPolicyResponse = {
  syncPolicy?: SyncPolicyItem
}

export type DeleteSyncPolicyRequest = {
  syncPolicyId?: number
}

export type DeleteSyncPolicyResponse = {
  syncPolicy?: SyncPolicyItem
}

export type ListSyncPoliciesRequest = {
  page?: number
  pageSize?: number
  search?: string
}

export type ListSyncPoliciesResponse = {
  syncPolicies?: SyncPolicyItem[]
  pagination?: MatrixhubV1alpha1Utils.Pagination
}

export type GetSyncPolicyRequest = {
  syncPolicyId?: number
}

export type GetSyncPolicyResponse = {
  syncPolicy?: SyncPolicyItem
}

export type CreateSyncTaskRequest = {
  syncPolicyId?: number
}

export type CreateSyncTaskResponse = {
  id?: number
}

export type SyncTask = {
  id?: number
  syncPolicyId?: number
  triggerType?: TriggerType
  status?: SyncTaskStatus
  startedTimestamp?: string
  completedTimestamp?: string
  totalItems?: string
  successfulItems?: string
  stoppedItems?: string
  failedItems?: string
}

export type ListSyncTasksRequest = {
  syncPolicyId?: number
  page?: number
  pageSize?: number
  status?: SyncTaskStatus
}

export type ListSyncTasksResponse = {
  syncTasks?: SyncTask[]
  pagination?: MatrixhubV1alpha1Utils.Pagination
}

export type StopSyncTaskRequest = {
  syncPolicyId?: number
  syncTaskId?: number
}

export type StopSyncTaskResponse = {
  syncTask?: SyncTask
}

export type SyncJob = {
  id?: number
  syncTaskId?: number
  resourceType?: ResourceType
  resourceName?: string
  targetResourceName?: string
  action?: string
  status?: SyncJobStatus
  completedTimestamp?: string
  createdTimestamp?: string
}

export type ListSyncJobsRequest = {
  syncPolicyId?: number
  syncTaskId?: number
  page?: number
  pageSize?: number
  status?: SyncJobStatus
  resourceType?: ResourceType
}

export type ListSyncJobsResponse = {
  syncJobs?: SyncJob[]
  pagination?: MatrixhubV1alpha1Utils.Pagination
}

export type GetSyncJobLogRequest = {
  syncPolicyId?: number
  syncTaskId?: number
  syncJobId?: number
}

export type GetSyncJobLogResponse = {
  log?: string
