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

