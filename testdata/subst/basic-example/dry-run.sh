#!/usr/bin/env bash
set -euo pipefail

# include env-vars
. vars.sh

# expand manifests
kubectl envsubst apply -f manifests/ --envsubst-allowed-prefixes=CI_,APP_,INFRA_ --dry-run=client -oyaml
