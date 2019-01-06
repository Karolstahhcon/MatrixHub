#!/bin/bash

# Copyright 2026 The MatrixHub Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ---------------------------------------------------------------------------
# E2E API coverage harness (ported from leopard's hack/utils.sh).
#
# Runs a mitmproxy-based recorder container that captures every API call the
# e2e suite makes, then diffs the recorded traffic against the OpenAPI specs
# with swagger-coverage-commandline and prints a per-module coverage report.
#
# The e2e HTTP client already honors http_proxy/https_proxy
# (test/tools/http_client.go uses http.ProxyFromEnvironment), so pointing those
# at the recorder is enough to capture traffic -- with one caveat handled by
# the caller: Go's ProxyFromEnvironment bypasses "localhost" and loopback IPs,
# so MATRIXHUB_BASE_URL must use a non-loopback hostname (e.g. matrixhub.local)
# for anything to be recorded.
#
# Usage (source, then call the lifecycle functions around the test run):
#   source scripts/api-coverage.sh
#   api_coverage::start
#   <run ginkgo>
#   api_coverage::report
#   api_coverage::stop
# ---------------------------------------------------------------------------

# Public, network-reachable mirror of the internal api-test-coverage image
# (mitmproxy recorder + swagger-coverage-commandline).
API_COVERAGE_IMAGE="${API_COVERAGE_IMAGE:-release.daocloud.io/matrixhub/api-test-coverage:1.0.0}"

# Directory holding the generated OpenAPI v2 specs.
API_COVERAGE_OPENAPI_DIR="${API_COVERAGE_OPENAPI_DIR:-./api/openapiv2}"

# Merged swagger spec passed to the recorder for the "total" coverage figure.
API_COVERAGE_MERGED_SWAGGER="${API_COVERAGE_MERGED_SWAGGER:-./merged_swagger.json}"

# Randomized container name so concurrent runs on one host do not collide.
API_COVERAGE_CONTAINER="${API_COVERAGE_CONTAINER:-matrixhub_api_coverage_${RANDOM}${RANDOM}}"

# Merge every *.swagger.json under $1 into a single swagger doc at $2.
# Combines paths + definitions; keeps the first file's top-level fields.
function api_coverage::merge_swagger() {
    local source_dir="${1:-${API_COVERAGE_OPENAPI_DIR}}"
    local output_file="${2:-${API_COVERAGE_MERGED_SWAGGER}}"

    if [ ! -d "${source_dir}" ]; then
        echo "api-coverage: openapi dir not found: ${source_dir}" >&2
        return 1
    fi

    local files=()
    while IFS= read -r f; do
        files+=("${f}")
    done < <(find "${source_dir}" -type f -name "*.swagger.json" | sort)
    if [ "${#files[@]}" -eq 0 ]; then
        echo "api-coverage: no *.swagger.json found under ${source_dir}" >&2
        return 1
    fi

    local main="${files[0]}"

    local tmp
    tmp=$(mktemp -d)
    jq -s 'map(.paths) | add' "${files[@]}" > "${tmp}/paths.json"
    jq -s 'map(.definitions) | add' "${files[@]}" > "${tmp}/definitions.json"
    jq 'del(.paths, .definitions)' "${main}" > "${tmp}/base.json"
    jq -s '
      .[0] as $base |
      .[1] as $paths |
      .[2] as $defs |
      $base + {paths: $paths} + {definitions: $defs}
    ' "${tmp}/base.json" "${tmp}/paths.json" "${tmp}/definitions.json" > "${output_file}"
    rm -rf "${tmp}"

    echo "api-coverage: merged ${#files[@]} specs into ${output_file}"
}

# Start the recorder container and export proxy env for the test process.
# Picks a random high port; the container runs on the host network so it can
# both listen on that port and forward to the service (via matrixhub.local).
function api_coverage::start() {
    local port=$(( (RANDOM % 30000) + 30000 ))

    export http_proxy="http://127.0.0.1:${port}"
