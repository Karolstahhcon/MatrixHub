#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXAMPLE_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
ENV_FILE="${EXAMPLE_ROOT}/env.sh"
ARTIFACTS_DIR="${EXAMPLE_ROOT}/artifacts"

MODEL_DEPLOYMENT="optimized-baseline-nvidia-gpu-vllm-decode"
MODEL_SERVICE="qwen3-32b"
MODEL_ID="Qwen/Qwen3-32B"
EXPERIMENT_LABEL_KEY="matrixhub.ai/experiment"
EXPERIMENT_LABEL_VALUE="llm-d-qwen3-32b"

die() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || die "required command not found: $1"
}

validate_namespace() {
  local value="$1"
  local variable="$2"
  [[ "$value" =~ ^[a-z0-9]([-a-z0-9]*[a-z0-9])?$ ]] || \
    die "${variable} must be a valid Kubernetes namespace name: ${value}"
}

load_env() {
  [[ -f "$ENV_FILE" ]] || die "${ENV_FILE} not found; run: cp env.example env.sh"

  # shellcheck disable=SC1090
  source "$ENV_FILE"

  : "${MATRIXHUB_ENDPOINT:?MATRIXHUB_ENDPOINT is required in env.sh}"
  : "${MATRIXHUB_NAMESPACE:?MATRIXHUB_NAMESPACE is required in env.sh}"
  : "${MATRIXHUB_DEPLOYMENT:?MATRIXHUB_DEPLOYMENT is required in env.sh}"
  : "${LLMD_NAMESPACE:?LLMD_NAMESPACE is required in env.sh}"
  : "${MODEL_NODE_SELECTOR:=}"
  : "${LLMD_ROUTER_CHART:=oci://ghcr.io/llm-d/charts/llm-d-router-standalone}"
  : "${LLMD_ROUTER_CHART_VERSION:=v0.9.0}"
  : "${LLMD_ROUTER_RELEASE:=optimized-baseline}"
  ROUTER_SERVICE="${LLMD_ROUTER_RELEASE}-epp"

  [[ "$MATRIXHUB_ENDPOINT" =~ ^https?:// ]] || \
