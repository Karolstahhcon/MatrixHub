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
E2E_MATRIXHUB_IMAGE=${E2E_MATRIXHUB_IMAGE:-"ghcr.io/matrixhub-ai/matrixhub:latest"}
E2E_JOBSERVER_POLL_INTERVAL=${E2E_JOBSERVER_POLL_INTERVAL:-"10s"}
E2E_MATRIXHUB_HOST_URL=${E2E_MATRIXHUB_HOST_URL:-"http://localhost:30001"}

echo "================================================"
echo "MatrixHub Deployment"
echo "================================================"
echo "Cluster Name:     ${E2E_CLUSTER_NAME}"
echo "MatrixHub Image:  ${E2E_MATRIXHUB_IMAGE}"
echo "JobServer Poll:   ${E2E_JOBSERVER_POLL_INTERVAL}"
echo "Git LFS Host URL: ${E2E_MATRIXHUB_HOST_URL}"
echo "================================================"

# Get the image registry and tag from the E2E_MATRIXHUB_IMAGE
# Format: ghcr.io/matrixhub-ai/matrixhub-ci:tag
IMAGE_REGISTRY=""
IMAGE_REPOSITORY=""
IMAGE_TAG=""

# Check if image contains a registry (has at least 2 slashes)
SLASH_COUNT=$(echo "${E2E_MATRIXHUB_IMAGE}" | tr -cd '/' | wc -c)
if [ "${SLASH_COUNT}" -ge 2 ]; then
    # Full format: registry/org/repo:tag
    IMAGE_REGISTRY=$(echo "${E2E_MATRIXHUB_IMAGE}" | cut -d'/' -f1)
    IMAGE_REPOSITORY=$(echo "${E2E_MATRIXHUB_IMAGE}" | cut -d'/' -f2- | cut -d':' -f1)
    IMAGE_TAG=$(echo "${E2E_MATRIXHUB_IMAGE}" | cut -d':' -f2)
else
    # Short format: org/repo:tag or repo:tag
    PARTS=$(echo "${E2E_MATRIXHUB_IMAGE}" | tr '/' '\n' | wc -l)
    if [ "${PARTS}" -eq 2 ]; then
        # org/repo:tag
        IMAGE_REPOSITORY=$(echo "${E2E_MATRIXHUB_IMAGE}" | cut -d':' -f1)
        IMAGE_TAG=$(echo "${E2E_MATRIXHUB_IMAGE}" | cut -d':' -f2)
    else
        # repo:tag
        IMAGE_REPOSITORY=$(echo "${E2E_MATRIXHUB_IMAGE}" | cut -d':' -f1)
        IMAGE_TAG=$(echo "${E2E_MATRIXHUB_IMAGE}" | cut -d':' -f2)
    fi
fi

echo ""
echo "Parsed image configuration:"
echo "  Full Image:  ${E2E_MATRIXHUB_IMAGE}"
echo "  Registry:    ${IMAGE_REGISTRY:-<default>}"
echo "  Repository:  ${IMAGE_REPOSITORY}"
