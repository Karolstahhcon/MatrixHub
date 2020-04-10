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
  }
}

// Simple function to strip comments from JSON (JSONC)
const parseJSONC = (str) => {
  // Remove single line // and multi-line /* */ comments
  // Note: This regex is simple and works for standard cases but might catch URLs etc. in strings.
  // For tsconfig structure it is generally safe.
  const cleanJson = str.replace(/\/\*[\s\S]*?\*\/|([^\\:]|^)\/\/.*$/gm, '$1')

  return JSON.parse(cleanJson)
}

const generateTSConfig = async (stagedFiles) => {
  if (!stagedFiles.length) {
    return []
  }

  // Read app and node configurations
  const [tsAppConfigRaw, tsNodeConfigRaw] = await Promise.all([
    fs.readFile(path.resolve('tsconfig.app.json'), 'utf8'),
    fs.readFile(path.resolve('tsconfig.node.json'), 'utf8'),
  ])

  const tsAppConfig = prepareTsConfigForLint(parseJSONC(tsAppConfigRaw))
  const tsNodeConfig = prepareTsConfigForLint(parseJSONC(tsNodeConfigRaw))

  const nodeFilesSet = expandTsIncludesToSet(tsNodeConfig)
  const appFilesSet = expandTsIncludesToSet(tsAppConfig)

  const resolvedStagedFiles = stagedFiles.map(f => path.resolve(f))

  const nodeFilesToCheck = resolvedStagedFiles.filter(f => nodeFilesSet.has(f))
  const appFilesToCheck = resolvedStagedFiles.filter(f => appFilesSet.has(f))

  const commands = []

  if (appFilesToCheck.length > 0) {
    // dynamic include for app + exclude node include to avoid duplicate checks
