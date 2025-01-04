#!/usr/bin/env bash
set -euo pipefail

# include env-vars
export APP_NAMESPACE='kubectl-envsubst-tests-6bf80e57-f7ba-4cab-9502-105cf669820b'
export APP_NAME=my-app
export IMAGE_NAME=nginx
export IMAGE_TAG=latest

# prepare namespace and context
kubectl create ns "${APP_NAMESPACE}" --dry-run=client -oyaml | kubectl apply -f -
kubectl config set-context --current --namespace="${APP_NAMESPACE}"

# configure CLI
export ENVSUBST_ALLOWED_PREFIXES='APP_,IMAGE_'

# expand and apply manifests
kubectl envsubst apply -f mix-single/

# restore context
kubectl config set-context --current --namespace="default"
