/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../fetch.pb"
import * as MatrixhubV1alpha1Model from "./model.pb"
import * as MatrixhubV1alpha1Utils from "./utils.pb"
export type ListDatasetTaskLabelsRequest = {
}

export type ListDatasetTaskLabelsResponse = {
  items?: MatrixhubV1alpha1Model.Label[]
}

export type ListDatasetsRequest = {
  label?: string[]
  search?: string
  sort?: string
  project?: string
  page?: number
  pageSize?: number
  popular?: boolean
}

export type ListDatasetsResponse = {
  items?: Dataset[]
  pagination?: MatrixhubV1alpha1Utils.Pagination
}

export type GetDatasetRequest = {
  project?: string
  name?: string
}

export type CreateDatasetRequest = {
  project?: string
  name?: string
}

export type CreateDatasetResponse = {
}

export type DeleteDatasetRequest = {
  project?: string
  name?: string
}

export type DeleteDatasetResponse = {
}

export type ListDatasetRevisionsRequest = {
  project?: string
  name?: string
}

export type ListDatasetRevisionsResponse = {
  items?: MatrixhubV1alpha1Model.Revisions
}

export type ListDatasetCommitsRequest = {
  project?: string
  name?: string
  revision?: string
  diff?: boolean
  page?: number
  pageSize?: number
}

export type ListDatasetCommitsResponse = {
  items?: MatrixhubV1alpha1Model.Commit[]
  pagination?: MatrixhubV1alpha1Utils.Pagination
}

export type GetDatasetCommitRequest = {
  project?: string
  name?: string
  id?: string
}

export type GetDatasetTreeRequest = {
  project?: string
  name?: string
  revision?: string
  path?: string
}

export type GetDatasetTreeResponse = {
  items?: MatrixhubV1alpha1Model.File[]
}

export type GetDatasetBlobRequest = {
  project?: string
  name?: string
  revision?: string
