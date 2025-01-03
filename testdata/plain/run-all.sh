#!/usr/bin/env bash
set -euo pipefail

echo "*****"
. 01-apply-recursive.sh

echo "*****"
. 02-apply-dir.sh

echo "*****"
. 03-apply-glob.sh