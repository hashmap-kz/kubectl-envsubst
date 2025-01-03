#!/usr/bin/env bash
set -euo pipefail

# include env-vars
. vars.sh

# prepare namespace and context
kubectl create ns "${APP_NAMESPACE}" --dry-run=client -oyaml | kubectl apply -f -
kubectl config set-context --current --namespace="${APP_NAMESPACE}"

# expand and apply manifests (handle both: stdin, filenames)
kubectl envsubst apply \
  -f https://raw.githubusercontent.com/hashmap-kz/kubectl-envsubst/refs/heads/master/test/manual/subst/basic/manifests/manifests.yaml \
  --envsubst-allowed-prefixes=CI_,APP_,INFRA_

# restore context
kubectl config set-context --current --namespace="default"
