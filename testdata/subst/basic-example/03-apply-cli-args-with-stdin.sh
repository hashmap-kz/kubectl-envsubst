#!/usr/bin/env bash
set -euo pipefail

# include env-vars
. vars.sh

# prepare namespace and context
kubectl create ns "${APP_NAMESPACE}" --dry-run=client -oyaml | kubectl apply -f -
kubectl config set-context --current --namespace="${APP_NAMESPACE}"

# expand and apply manifests (handle both: stdin, filenames)
cat separated/service.yaml | kubectl envsubst apply \
  -f - \
  --filename=separated/deployment.yaml \
  -f separated/secret.yaml \
  --envsubst-allowed-prefixes=CI_,APP_,INFRA_
