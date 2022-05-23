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
