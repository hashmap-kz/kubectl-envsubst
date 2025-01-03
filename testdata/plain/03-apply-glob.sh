#!/usr/bin/env bash
set -euo pipefail

(
  cd manifests
  kubectl envsubst apply -f '*.yaml'
)
