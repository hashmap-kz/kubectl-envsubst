#!/usr/bin/env bash
set -euo pipefail

# include env-vars
. vars.sh

# prepare namespace and context
kubectl create ns "${APP_NAMESPACE}" --dry-run=client -oyaml | kubectl apply -f -
kubectl config set-context --current --namespace="${APP_NAMESPACE}"

# configure CLI
export ENVSUBST_ALLOWED_PREFIXES='CI_,APP_,INFRA_'

# expand and apply manifests
kubectl envsubst apply -f manifests/
