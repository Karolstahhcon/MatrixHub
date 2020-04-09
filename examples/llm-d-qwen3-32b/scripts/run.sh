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
    kill "$port_forward_pid" >/dev/null 2>&1 || true
    wait "$port_forward_pid" 2>/dev/null || true
  fi
  if (( status != 0 )); then
    if is_experiment_namespace "$namespace"; then
      collect_diagnostics "$namespace" "$run_dir"
      printf 'Run failed; diagnostics were saved in %s. The Router namespace was retained; run make cleanup when finished.\n' \
        "$run_dir" >&2
    else
      printf 'Run failed before the Router prerequisite could be verified.\n' >&2
    fi
  fi
  exit "$status"
}

require_router_prerequisite "$namespace"
trap collect_on_failure EXIT
remove_model_workload "$namespace"

if [[ "$mode" == matrixhub ]]; then
  kubectl create configmap matrixhub-hf-endpoint \
    -n "$namespace" \
    --from-literal="HF_ENDPOINT=${MATRIXHUB_ENDPOINT}" \
    --dry-run=client -o yaml | kubectl apply -f - >/dev/null
fi

rendered_manifest="${run_dir}/rendered.yaml"
if [[ -n "$MODEL_NODE_SELECTOR" ]]; then
  node_selector_overlay_dir="${run_dir}/node-selector-overlay"
  mkdir -p "$node_selector_overlay_dir"
  cat >"${node_selector_overlay_dir}/kustomization.yaml" <<EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../../../overlays/${mode}
patches:
