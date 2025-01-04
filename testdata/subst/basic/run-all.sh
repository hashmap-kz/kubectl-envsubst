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
