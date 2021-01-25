import fs from 'fs/promises'
import path from 'path'
import { ESLint } from 'eslint'
import ts from 'typescript'

const projectRoot = process.cwd()

const prepareTsConfigForLint = (config) => {
  if (!config.compilerOptions) {
    config.compilerOptions = {}
  }
  delete config.compilerOptions.composite
  delete config.compilerOptions.tsBuildInfoFile
  const currentBaseUrl = config.compilerOptions.baseUrl ?? '.'

  config.compilerOptions.baseUrl = path.resolve(projectRoot, currentBaseUrl)

  return config
}

// Try to use TypeScript's file expansion ability to parse include patterns into specific file path sets
const expandTsIncludesToSet = (config) => {
  const emptySet = new Set()

  if (!config?.include?.length) {
    return emptySet
  }
  if (!ts || !ts.sys || !ts.sys.readDirectory) {
    return emptySet
  }

  try {
    // ts.sys.readDirectory(baseDir, extensions, excludes, includes)
    const files = ts.sys.readDirectory(projectRoot, ['ts', 'tsx', 'js', 'jsx'], config?.exclude, config.include)

    return new Set(files.map(f => path.resolve(f)))
  } catch {
    return emptySet
