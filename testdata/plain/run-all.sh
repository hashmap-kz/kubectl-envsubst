#!/usr/bin/env bash
set -euo pipefail

echo "*** recursive **"
. 01-apply-recursive.sh

echo "*** directory **"
. 02-apply-dir.sh

echo "*** globbing **"
. 03-apply-glob.sh
