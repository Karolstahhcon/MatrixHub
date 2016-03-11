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
