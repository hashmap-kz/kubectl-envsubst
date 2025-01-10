#!/bin/bash
set -euo pipefail

# setup envs
export IMAGE_NAME='monopole/hello'
export IMAGE_TAG='1'

export ENVSUBST_ALLOWED_PREFIXES='IMAGE_'
kubectl kustomize manifests/ | kubectl envsubst apply -f -
