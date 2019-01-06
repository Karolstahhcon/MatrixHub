#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "${SCRIPT_DIR}/common.sh"

preflight_local

direct_output="${ARTIFACTS_DIR}/direct-rendered.yaml"
matrixhub_output="${ARTIFACTS_DIR}/matrixhub-rendered.yaml"

kubectl kustomize "${EXAMPLE_ROOT}/overlays/direct" >"$direct_output"
