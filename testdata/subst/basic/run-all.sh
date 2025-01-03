#!/usr/bin/env bash
set -euo pipefail

echo "*****"
. 01-apply-cli-args.sh

echo "*****"
. 02-apply-env-args.sh

echo "*****"
. 03-apply-stdin.sh

echo "*****"
. 04-apply-url.sh
