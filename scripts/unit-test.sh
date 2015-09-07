#!/usr/bin/env bash

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

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="${SCRIPT_DIR}/.."

cd "${PROJECT_ROOT}"

UNIT_TEST_PKGS="${UNIT_TEST_PKGS:-./cmd/... ./internal/...}"
UNIT_TEST_EXCLUDE_PKGS="${UNIT_TEST_EXCLUDE_PKGS-}"
COUNT="${COUNT:-1}"
VERBOSE="${VERBOSE:-false}"
RACE="${RACE:-false}"
COVERAGE="${COVERAGE:-false}"
COVERAGE_PROFILE="${COVERAGE_PROFILE:-coverage.out}"
COVERAGE_HTML="${COVERAGE_HTML:-false}"
COVERAGE_HTML_FILE="${COVERAGE_HTML_FILE:-${COVERAGE_PROFILE%.*}.html}"

args=(-short "-count=${COUNT}")

if [[ "${VERBOSE}" == "true" ]]; then
	args+=(-v)
fi

if [[ "${RACE}" == "true" ]]; then
	args+=(-race)
fi

if [[ "${COVERAGE}" == "true" ]]; then
	args+=("-coverprofile=${COVERAGE_PROFILE}" -covermode=atomic)
fi

# Intentional word splitting: package variables are space-separated lists of Go
# package patterns, e.g. "./cmd/... ./internal/...".
read -r -a package_patterns <<< "${UNIT_TEST_PKGS}"

if [[ ${#package_patterns[@]} -eq 0 ]]; then
	echo "UNIT_TEST_PKGS is empty" >&2
	exit 1
fi

packages=()
while IFS= read -r pkg; do
	# Generated mock packages contain test scaffolding rather than production
	# logic. Tests that use them still compile them as dependencies, while
