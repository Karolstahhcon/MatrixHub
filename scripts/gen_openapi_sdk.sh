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

    mkdir -p "$client_path"

    local versions
    versions=$(ls "$swagger_path")

    for version in $versions; do
        local files
        files=$(ls "$swagger_path/$version")
        for file in $files; do
            if [[ $file == *.swagger.json ]]; then
                local name="${file%.swagger.json}"
                local target="$client_path/$version/$name"
                rm -rf "$target" && mkdir -p "$target"
                if docker info >/dev/null 2>&1; then
                    docker run --rm -u "$(id -u):$(id -g)" \
                        -v "$swagger_path/$version":/swagger \
