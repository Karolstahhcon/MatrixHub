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

if ! [[ "${DEST_TAG}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "error: dest tag must be an official release (vX.Y.Z), got ${DEST_TAG}" >&2
  exit 1
fi

mkdir -p "${OUTPUT_DIR}"
OUTPUT_DIR=$(cd "${OUTPUT_DIR}" && pwd)

VERSION="${DEST_TAG#v}"
MAJOR_MINOR=$(echo "${VERSION}" | cut -d. -f1,2)
CHANGELOG_PATH="CHANGELOG/CHANGELOG-${MAJOR_MINOR}.md"
OUTPUT_FILE="${OUTPUT_DIR}/release_notes_${DEST_TAG}.md"

extract_section() {
  local file=$1
  awk -v tag="${DEST_TAG}" '
    BEGIN { found=0 }
    $0 ~ "^## " tag "$" { found=1; print; next }
    found && $0 ~ "^## v[0-9]+\\.[0-9]+\\.[0-9]+" { exit }
    found { print }
  ' "${file}"
}
