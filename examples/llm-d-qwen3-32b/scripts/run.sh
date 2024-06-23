#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "${SCRIPT_DIR}/common.sh"

mode="${1:-}"
case "$mode" in
  direct)
    overlay="${EXAMPLE_ROOT}/overlays/direct"
    ;;
  matrixhub)
    overlay="${EXAMPLE_ROOT}/overlays/matrixhub"
    ;;
  *)
    die "usage: $0 direct|matrixhub"
    ;;
esac

preflight_cluster
require_command python3
require_command curl

namespace="$LLMD_NAMESPACE"

run_id="$(date -u +%Y%m%dT%H%M%SZ)"
run_dir="${ARTIFACTS_DIR}/${mode}/${run_id}"
mkdir -p "$run_dir"
port_forward_pid=""

collect_on_failure() {
  local status="$?"
  trap - EXIT
  if [[ -n "$port_forward_pid" ]]; then
