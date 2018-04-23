#!/usr/bin/env bash
set -euo pipefail

example_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

test -f "${example_root}/router/base.values.yaml"
test -f "${example_root}/router/optimized-baseline.values.yaml"
grep -Fqx '    protocol: ""' "${example_root}/router/base.values.yaml"
test -f "${example_root}/overlays/direct/patch-hf-endpoint.yaml"
grep -Fqx '              value: https://hf.m.daocloud.io' \
  "${example_root}/overlays/direct/patch-hf-endpoint.yaml"
grep -Fqx '  - path: patch-hf-endpoint.yaml' \
  "${example_root}/overlays/direct/kustomization.yaml"
grep -Fqx 'LLMD_ROUTER_CHART="oci://ghcr.io/llm-d/charts/llm-d-router-standalone"' \
  "${example_root}/env.example"
grep -Fqx 'LLMD_ROUTER_CHART_VERSION="v0.9.0"' "${example_root}/env.example"
grep -Fqx 'LLMD_NAMESPACE="matrixhub-llmd-qwen3-32b"' "${example_root}/env.example"
grep -Fqx 'MODEL_NODE_SELECTOR=""' "${example_root}/env.example"
grep -Fq 'helm upgrade --install' "${example_root}/scripts/common.sh"
grep -Fq 'grep -qx inferencepools.inference.networking.k8s.io' \
  "${example_root}/scripts/common.sh"
grep -Fq 'kubernetes.io/hostname: ${MODEL_NODE_SELECTOR}' "${example_root}/scripts/run.sh"
grep -Fq 'service/"$ROUTER_SERVICE"' "${example_root}/scripts/run.sh"
grep -Fq 'https://hf.m.daocloud.io' "${example_root}/scripts/run.sh"
test -f "${example_root}/scripts/prepare-router.sh"
grep -Fq 'prepare-router' "${example_root}/Makefile"
