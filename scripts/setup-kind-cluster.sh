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
# Project root is two levels up from the script directory
PROJECT_ROOT="${SCRIPT_DIR}/.."

# Change to project root directory
cd "${PROJECT_ROOT}"
echo "Working directory: $(pwd)"

# Environment variables with defaults
E2E_CLUSTER_NAME=${E2E_CLUSTER_NAME:-"matrixhub-e2e"}
E2E_KIND_IMAGE_TAG=${E2E_KIND_IMAGE_TAG:-"v1.32.3"}

echo "================================================"
echo "MatrixHub Kind Cluster Setup"
echo "================================================"
echo "Cluster Name: ${E2E_CLUSTER_NAME}"
echo "K8s Version:   ${E2E_KIND_IMAGE_TAG}"
echo "================================================"

# Check if kind is installed
if ! command -v kind &> /dev/null; then
    echo "Error: kind is not installed"
    exit 1
fi

# Check kind version (require >= 0.27)
KIND_VERSION=$(kind version 2>/dev/null | grep -oE 'kind v[0-9]+\.[0-9]+\.[0-9]+' | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
MIN_KIND_VERSION="0.27.0"

if [ -z "$KIND_VERSION" ]; then
    echo "Error: Unable to determine kind version"
    exit 1
fi

# Compare versions using sort
if [ "$(printf '%s\n' "$MIN_KIND_VERSION" "$KIND_VERSION" | sort -V | head -n1)" != "$MIN_KIND_VERSION" ]; then
    echo "Error: kind version $KIND_VERSION is too old (require >= $MIN_KIND_VERSION)"
    exit 1
fi
