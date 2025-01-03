#!/usr/bin/env bash
set -euo pipefail

# gitlab-specific variables, that defined in each pipeline
# we simulate that pipeline by setting them by hand
export CI_REGISTRY=mirror-0.company.org:5000
export CI_PROJECT_ROOT_NAMESPACE=banking-system-envsubst-test
export CI_PROJECT_NAMESPACE=backend/auth-svc
export CI_PROJECT_PATH=banking-system-envsubst-test/backend/auth-svc
export CI_PROJECT_NAME=auth-svc
export CI_COMMIT_REF_NAME=dev
export APP_IMAGE=nginx
export APP_NAMESPACE="${CI_PROJECT_ROOT_NAMESPACE}-${CI_COMMIT_REF_NAME}"
export INFRA_DOMAIN_NAME=company.org
