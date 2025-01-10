#!/bin/bash
set -euo pipefail

# setup envs
export PROJECT_ROOT_NAMESPACE=kubectl-envubst-examples
export PROJECT_ENV=prod
export PROJECT_NAME=nginx-gateway
export PROJECT_NAMESPACE="${PROJECT_ROOT_NAMESPACE}-${PROJECT_ENV}"
export IMAGE_NAME=nginx
export IMAGE_TAG=latest

# setup namespace and context
kubectl create ns "${PROJECT_NAMESPACE}" --dry-run=client -oyaml | kubectl apply -f -
kubectl config set-context --current --namespace="${PROJECT_NAMESPACE}"

# substitute and apply resources, according to the environment (dev, stage, prod)
export ENVSUBST_ALLOWED_PREFIXES='PROJECT_,IMAGE_'
kubectl envsubst apply -f "${PROJECT_ENV}"
