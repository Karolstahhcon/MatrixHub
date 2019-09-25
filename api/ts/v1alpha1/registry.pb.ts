/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../fetch.pb"
import * as GoogleProtobufTimestamp from "../google/protobuf/timestamp.pb"
import * as MatrixhubV1alpha1Utils from "./utils.pb"

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };
type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (
    keyof T extends infer K ?
      (K extends string & keyof T ? { [k in K]: T[K] } & Absent<T, K>
        : never)
    : never);

export enum RegistryType {
  REGISTRY_TYPE_UNSPECIFIED = "REGISTRY_TYPE_UNSPECIFIED",
  REGISTRY_TYPE_HUGGINGFACE = "REGISTRY_TYPE_HUGGINGFACE",
  REGISTRY_TYPE_MATRIXHUB = "REGISTRY_TYPE_MATRIXHUB",
}

export enum RegistryStatus {
  REGISTRY_STATUS_UNSPECIFIED = "REGISTRY_STATUS_UNSPECIFIED",
  REGISTRY_STATUS_HEALTHY = "REGISTRY_STATUS_HEALTHY",
  REGISTRY_STATUS_UNHEALTHY = "REGISTRY_STATUS_UNHEALTHY",
}


type BaseRegistry = {
  id?: number
  name?: string
  description?: string
  type?: RegistryType
  url?: string
  insecure?: boolean
  status?: RegistryStatus
  createdAt?: GoogleProtobufTimestamp.Timestamp
  updatedAt?: GoogleProtobufTimestamp.Timestamp
}

export type Registry = BaseRegistry
  & OneOf<{ basic: RegistryBasicCredential }>

export type RegistryBasicCredential = {
  username?: string
  password?: string
}

export type ListRegistriesRequest = {
  page?: number
  pageSize?: number
  search?: string
}

export type ListRegistriesResponse = {
  registries?: Registry[]
  pagination?: MatrixhubV1alpha1Utils.Pagination
}

export type GetRegistryRequest = {
  id?: number
}

export type GetRegistryResponse = {
  registry?: Registry
}


type BaseCreateRegistryRequest = {
  name?: string
  description?: string
  type?: RegistryType
  url?: string
  insecure?: boolean
}

export type CreateRegistryRequest = BaseCreateRegistryRequest
  & OneOf<{ basic: RegistryBasicCredential }>

export type CreateRegistryResponse = {
  registry?: Registry
}


type BaseUpdateRegistryRequest = {
  id?: number
  name?: string
  description?: string
  url?: string
  insecure?: boolean
}
