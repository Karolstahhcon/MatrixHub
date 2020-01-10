#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DIR="$(dirname "${BASH_SOURCE[0]}")"
ROOT_DIR="$(realpath "${DIR}/..")"

cd "${ROOT_DIR}"

if ! command -v gh >/dev/null 2>&1; then
  echo "gh is required to validate GitHub action commit and release refs"
  exit 1
fi

status=0

while IFS= read -r line; do
  file="${line%%:*}"
  rest="${line#*:}"
  lineno="${rest%%:*}"
  value="${line#*uses:}"
  value="${value%%#*}"
  value="$(echo "${value}" | xargs)"
  value="${value#\"}"
  value="${value%\"}"
  value="${value#\'}"
  value="${value%\'}"

  if [[ -z "${value}" || "${value}" == ./* || "${value}" == docker://* ]]; then
    continue
