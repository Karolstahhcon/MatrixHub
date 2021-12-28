#!/bin/bash

set -e
set -o pipefail
set -u
set -x

bash_source_path="$(dirname "${BASH_SOURCE[0]}")"
pro_path=$(cd "$bash_source_path/.."; pwd)

SWAGGER_CODEGEN_VERSION="${SWAGGER_CODEGEN_VERSION:-2.4.41}"

function gen_openapi_sdk() {
    swagger_path=${1:-"$pro_path/api/openapiv2"}
    client_path="$pro_path/test/client"

