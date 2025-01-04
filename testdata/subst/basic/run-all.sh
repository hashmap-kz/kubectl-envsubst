#!/usr/bin/env bash
set -euo pipefail

echo "*** cli-args **"
. 01-apply-cli-args.sh

echo "*** env-args **"
. 02-apply-env-args.sh

echo "*** stdin **"
. 03-apply-stdin.sh

echo "*** url **"
. 04-apply-url.sh

echo "*** mixed-single **"
. 05-apply-mix.sh

echo "*** cleanup **"
kubectl delete ns "kubectl-envsubst-tests-6bf80e57-f7ba-4cab-9502-105cf669820b" --ignore-not-found
