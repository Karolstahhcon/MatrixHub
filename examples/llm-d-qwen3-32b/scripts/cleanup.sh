#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=common.sh
source "${SCRIPT_DIR}/common.sh"

preflight_local

mode="${1:-model}"
assert_experiment_namespace "$LLMD_NAMESPACE"

case "$mode" in
  model)
