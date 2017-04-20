#!/bin/bash -e
# This script brings a cluster up then brings the cluster down.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/../../.."

function shutdown() {
  echo "Bringing cluster down"
  ${ROOT}/tests/scripts/aws/down.sh
}

echo "Bringing cluster up"
source ${ROOT}/tests/scripts/aws/up.sh
echo "Cluster running."
trap shutdown EXIT

CLUSTER_DIR=${CLUSTER_DIR:-"tests/scripts/aws/output/${CLUSTER_NAME}"}
kubeconfig="${CLUSTER_DIR}/assets/auth/kubeconfig"
KUBECONFIG="${ROOT}/${kubeconfig}"
echo "Printing debug information:"
OS=$(uname | tr A-Z a-z)
curl -O https://storage.googleapis.com/kubernetes-release/release/v1.5.5/bin/${OS}/amd64/kubectl && chmod +x ./kubectl
./kubectl --kubeconfig=${KUBECONFIG} get pods --all-namespaces
./kubectl --kubeconfig=${KUBECONFIG} get nodes

# the current tectonic configuration is num workers + 1 controller
export NODE_COUNT=$(expr ${WORKER_COUNT} + 1)
project="github.com/coreos/tectonic-installer"
export TEST_KUBECONFIG=${KUBECONFIG}

${ROOT}/bin/sanity
