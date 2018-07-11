#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DIR="$(dirname "${BASH_SOURCE[0]}")"
ROOT_DIR="$(realpath "${DIR}/..")"
FETCH_FILE="${ROOT_DIR}/api/ts/fetch.pb.ts"

if [[ ! -f "${FETCH_FILE}" ]]; then
  echo "missing ${FETCH_FILE}" >&2
  exit 1
fi

if grep -q "export type FetchFn = typeof fetch" "${FETCH_FILE}"; then
  exit 0
fi

grep -q "export interface InitReq extends RequestInit" "${FETCH_FILE}"
grep -q "return fetch(url, req)" "${FETCH_FILE}"
