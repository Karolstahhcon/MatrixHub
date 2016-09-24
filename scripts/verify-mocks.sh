#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DIR="$(dirname "${BASH_SOURCE[0]}")"
ROOT_DIR="$(realpath "${DIR}/..")"

cd "${ROOT_DIR}"

go generate ./...

if ! git diff --quiet -- '**/mocks/**'; then
  echo "mock files differ from committed versions"
  git --no-pager diff -- '**/mocks/**'
  exit 1
fi
