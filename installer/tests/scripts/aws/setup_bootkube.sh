#!/bin/bash -e
# This script transfers bootkube assets to a cluster and starts it.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/../../.."

CLUSTER_DIR=${CLUSTER_DIR:-"${ROOT}/tests/scripts/aws/output/${CLUSTER_NAME}"}
assets_dir="${CLUSTER_DIR}/assets"

no_hostcheck="-o stricthostkeychecking=no"
ssh_args="${no_hostcheck}"
target="core@${CLUSTER_DOMAIN}"

until $(curl --silent --fail -m 1 "http://${CLUSTER_DOMAIN}:10255/healthz" > /dev/null); do
  echo "Waiting for Kubelets to start..."
  sleep 15
done

echo "Starting bootkube..."
ssh ${ssh_args} ${target} -- 'sudo systemctl start bootkube'

until $(curl --silent --fail -k https://${TECTONIC_DOMAIN} > /dev/null); do
  echo "Waiting for Tectonic Console..."
  sleep 10
done

ssh ${ssh_args} ${target} -- 'sudo journalctl -u bootkube'
