#!/bin/bash -e
# This script brings a cluster up.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/../../.."

echo "setup environment"
source "${ROOT}/tests/scripts/aws/default.env.sh"

echo "start cluster"
${ROOT}/tests/scripts/aws/launch_cluster_aws.sh
trap shutdown EXIT

echo "wait for machines to come up"
sleep 200
${ROOT}/tests/scripts/aws/wait_for_dns.sh

echo "transfer config and start bootkube"
${ROOT}/tests/scripts/aws/setup_bootkube.sh

echo "Bootkube running!"