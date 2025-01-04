#!/usr/bin/env bash
set -euo pipefail

# include env-vars
export APP_NAMESPACE='kubectl-envsubst-tmpns-testmix'
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
