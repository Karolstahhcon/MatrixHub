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
