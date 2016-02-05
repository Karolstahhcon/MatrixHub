/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../fetch.pb"
import * as GoogleProtobufWrappers from "../google/protobuf/wrappers.pb"
import * as MatrixhubV1alpha1Utils from "./utils.pb"

export enum FileType {
  DIR = "DIR",
  FILE = "FILE",
}

export enum Category {
  TASK = "TASK",
  LIBRARY = "LIBRARY",
  LICENSE = "LICENSE",
  LANGUAGE = "LANGUAGE",
  OTHER = "OTHER",
}

export type ListModelTaskLabelsRequest = {
}

export type ListModelTaskLabelsResponse = {
  items?: Label[]
}

export type ListModelFrameLabelsRequest = {
}

export type ListModelFrameLabelsResponse = {
  items?: Label[]
}

export type ListModelsRequest = {
  labels?: string[]
  search?: string
  sort?: string
  project?: string
  page?: number
  pageSize?: number
  popular?: boolean
}

export type ListModelsResponse = {
  items?: Model[]
  pagination?: MatrixhubV1alpha1Utils.Pagination
}

export type GetModelRequest = {
  project?: string
  name?: string
}

export type CreateModelRequest = {
  project?: string
  name?: string
}

export type CreateModelResponse = {
}

export type DeleteModelRequest = {
  project?: string
  name?: string
}

export type DeleteModelResponse = {
}

export type ListModelRevisionsRequest = {
  project?: string
  name?: string
}

export type ListModelRevisionsResponse = {
  items?: Revisions
}

export type Revisions = {
  branches?: Revision[]
  tags?: Revision[]
}

export type ListModelCommitsRequest = {
  project?: string
  name?: string
  revision?: string
  diff?: boolean
  page?: number
  pageSize?: number
}

export type ListModelCommitsResponse = {
  items?: Commit[]
  pagination?: MatrixhubV1alpha1Utils.Pagination
}

export type GetModelCommitRequest = {
  project?: string
  name?: string
  id?: string
}

export type GetModelTreeRequest = {
  project?: string
  name?: string
  revision?: string
  path?: string
}

export type File = {
