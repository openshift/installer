#!/bin/bash -e
# This script uses the installer to launch a cluster on AWS. It requires that AWS credentials be provided as environment variables.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/../../.."

CLUSTER_DIR=${CLUSTER_DIR:-"${ROOT}/tests/scripts/aws/output/${CLUSTER_NAME}"}
CLUSTER_CREATE_URL="http://127.0.0.1:4444/cluster/create"
OS=$(uname | tr A-Z a-z)
INSTALLER_BIN=${INSTALLER_BIN:-"${ROOT}/bin/${OS}/installer"}
PAYLOAD=$(source ${ROOT}/tests/scripts/aws/aws_payload.tmpl.sh)

source "${ROOT}/scripts/awsutil.sh"
# set environment variables in case they are not yet set
source "${ROOT}/tests/scripts/aws/default.env.sh"
# check for AWS credentials
check_aws_creds

cleanup() {
  echo "Killing Tectonic Installer"
  kill ${INSTALLER_PID} || true
}

trap cleanup EXIT

# check that cluster dir does not already exist
if [ -d "${CLUSTER_DIR}" ]; then
  echo "Cluster already exists at ${CLUSTER_DIR}. Change \$CLUSTER_NAME, delete the existing cluster, or set \$CLUSTER_DIR."
  exit 1
fi
mkdir -p ${CLUSTER_DIR}
env_file="${CLUSTER_DIR}/env"

# write cluster configuration
cat << EOF > ${env_file}
export CLUSTER_NAME=${CLUSTER_NAME}
export AWS_REGION=${AWS_REGION}
export AWS_HOSTEDZONE=${AWS_HOSTEDZONE}
export CLUSTER_DOMAIN=${CLUSTER_DOMAIN}
export TECTONIC_DOMAIN=${TECTONIC_DOMAIN}
export UPDATER_ENABLED=${UPDATER_ENABLED}
export UPDATER_SERVER=${UPDATER_SERVER}
export UPDATER_CHANNEL=${UPDATER_CHANNEL}
export UPDATER_APPID=${UPDATER_APPID}
export TECTONIC_LICENSE=${TECTONIC_LICENSE}
export TECTONIC_PULL_SECRET=$(echo ${TECTONIC_PULL_SECRET} | sed 's/\\/\\\\/g')
EOF

echo "A cluster configuration has been written to ${CLUSTER_DIR}."
echo "In order to use this cluster, you may need to run 'source ${env_file}'"

echo "Starting Tectonic Installer"
${INSTALLER_BIN} -log-level=debug -open-browser=false -platforms=aws  &
INSTALLER_PID=$!
sleep 2

echo "Submitting to Tectonic Installer"
assets_zip="${CLUSTER_DIR}/assets.zip"
echo "${PAYLOAD}" | >${assets_zip} curl -v -X POST --data-binary @- ${CLUSTER_CREATE_URL}
echo "Provisioning requested"

# Killing quickly because multiple tests are using this port
echo "Killing Tectonic Installer"
kill ${INSTALLER_PID} || true

echo "Unzipping self-hosting provisioning assets"
unzip ${assets_zip} -d ${CLUSTER_DIR} > /dev/null

