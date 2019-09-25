#!/usr/bin/env node

import { mkdir, writeFile } from 'node:fs/promises'
import path from 'node:path'
import process from 'node:process'

function printUsage() {
  process.stdout.write(`Usage:
  Single:
    node path/to/import-figma-svg.mjs --url <url> --name <file-name> [--dir <output-dir>]

  Batch:
    node path/to/import-figma-svg.mjs --items '<json-array>' [--dir <output-dir>] [--concurrency <n>]

  The JSON array format: [{"url":"...","name":"..."},...]

Options:
  --dir          Output directory (default: current working directory)
  --concurrency  Max parallel downloads for batch mode (default: 4)

Examples:
  node path/to/import-figma-svg.mjs --url "https://figma.com/api/mcp/asset/..." --name admin-users-nav --dir ./svgs
  node path/to/import-figma-svg.mjs --items '[{"url":"https://...","name":"icon-a"},{"url":"https://...","name":"icon-b"}]' --dir ./svgs --concurrency 6
`)
}

function parseArgs(argv) {
  const args = {
    concurrency: 4,
    dir: '.',
    items: '',
    name: '',
    url: '',
  }

  for (let index = 0; index < argv.length; index += 1) {
    const arg = argv[index]

    if (arg === '--help' || arg === '-h') {
      printUsage()
      process.exit(0)
    }
