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
  projectName?: string
  resourceName?: string
  sizeBytes?: string
}

export type OrphanedLFS = {
  oid?: string
  sizeBytes?: string
}

export type ExecuteCleanupRequest = {
  cleanOrphanedRepos?: boolean
  cleanOrphanedLfs?: boolean
  dryRun?: boolean
}

export type CleanupResult = {
  reposDeleted?: number
  lfsObjectsDeleted?: number
  spaceReclaimedBytes?: string
  errors?: string[]
}
