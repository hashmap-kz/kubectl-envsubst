#!/usr/bin/env bash
set -euo pipefail

kubectl envsubst apply -f manifests/ --recursive
