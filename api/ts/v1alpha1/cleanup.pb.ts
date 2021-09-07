/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../fetch.pb"
export type PreviewCleanupRequest = {
  includeOrphanedRepos?: boolean
  includeOrphanedLfs?: boolean
}

export type CleanupPreview = {
  orphanedRepos?: OrphanedRepo[]
  orphanedLfsObjects?: OrphanedLFS[]
  totalReclaimableBytes?: string
}

export type OrphanedRepo = {
  path?: string
  type?: string
