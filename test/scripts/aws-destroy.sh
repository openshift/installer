#!/bin/bash -ex
set -o pipefail
shopt -s expand_aliases

DIR="$( cd "$( dirname "$0" )" && pwd )"
# shellcheck disable=SC1090
source "$DIR/aws-common.sh"

echo "Destroying ${CLUSTER}..."
make destroy
