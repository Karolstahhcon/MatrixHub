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

set -e

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# Project root is one level up from the script directory
PROJECT_ROOT="${SCRIPT_DIR}/.."

# Change to project root directory
cd "${PROJECT_ROOT}"

# Ginkgo label filter expression (e.g. "smoke", or comma-separated labels).
# Empty or "all" runs every test (no filter). Passed to `ginkgo --label-filter`.
LABEL_FILTER=${1:-}

# Base URL from environment or default local API server.
MATRIXHUB_BASE_URL="${MATRIXHUB_BASE_URL:-http://localhost:3001}"

# URL by which MatrixHub reaches itself, resolved inside the deployment rather
