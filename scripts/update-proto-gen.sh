#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DIR="$(dirname "${BASH_SOURCE[0]}")"

ROOT_DIR="$(realpath "${DIR}/..")"

PROTOVENDOR_DIR="${ROOT_DIR}/api/proto/vendor"
PROTO_DIR="${ROOT_DIR}/api/proto"
SWAGGER_DIR="${ROOT_DIR}/api/openapiv2"
TS_DIR="${ROOT_DIR}/api/ts"
GO_DIR="${ROOT_DIR}/api/go"
GOBIN="${ROOT_DIR}/bin"
PATH="${GOBIN}:${PATH}"

# Install protoc
