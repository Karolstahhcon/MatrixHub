#!/bin/bash

# Copyright 2026 MatrixHub Authors
# SPDX-License-Identifier: Apache-2.0

# Generate GitHub Release notes for an official tag from CHANGELOG/ files.
#
# Usage:
#   ./scripts/release-notes.sh <output-dir> <dest-tag>
#
# If CHANGELOG is not present at the checked-out ref (for example when promoting
# v0.1.0 on the same commit as v0.1.0-rc.1), the script falls back to
# origin/main for the CHANGELOG file.

set -euo pipefail

OUTPUT_DIR=${1:-}
DEST_TAG=${2:-}

if [ -z "${OUTPUT_DIR}" ] || [ -z "${DEST_TAG}" ]; then
  echo "usage: $0 <output-dir> <dest-tag>" >&2
  exit 1
fi
