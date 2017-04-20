#!/bin/bash -e
# This script brings a cluster down.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/../../.."

echo "stopping cluster"
${ROOT}/tests/scripts/aws/shutdown_aws.sh
