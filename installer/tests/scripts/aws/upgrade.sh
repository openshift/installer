#!/bin/bash -e

function shutdown() {
  echo "Shutting down the cluster"

  ${ROOT}/tests/scripts/aws/down.sh
}

function upload_logs() {
  echo "Uploading ${1} log"

  podname=`${KUBECTL} get pods -n tectonic-system | grep ${1}-* | cut -d " " -f1`

  echo "Pod name for ${1}: ${podname}"

  logname="${CLUSTER_NAME}-${1}.log"
  bucketname="tectonic-upgrade-test-logs"

  ${KUBECTL} logs ${podname} -n tectonic-system > ${ROOT}/${logname}
  aws_upload_file ${ROOT}/${logname} ${logname} ${bucketname}
  rm "${ROOT}/${logname}"

  echo "${1} log uploaded at https://s3-us-west-2.amazonaws.com/tectonic-upgrade-test-logs/${logname}"
}

function cleanup_coreupdate() {
  echo "Deleting package ${VERSION}"

  ${UPDATESERVICECTL} --server ${SERVER} \
                      --key ${COREUPDATE_KEY} \
                      --user ${COREUPDATE_USR} \
                      package delete \
                      --app-id ${APPID} \
                      --version ${VERSION}

  echo "Deleting channel ${CHANNEL}"
  ${UPDATESERVICECTL} --server ${SERVER} \
                      --key ${COREUPDATE_KEY} \
                      --user ${COREUPDATE_USR} \
                      channel delete \
                      --app-id ${APPID} \
                      --channel ${CHANNEL}

  echo "Deleting group ${CHANNEL}"
  ${UPDATESERVICECTL} --server ${SERVER} \
                      --key ${COREUPDATE_KEY} \
                      --user ${COREUPDATE_USR} \
                      group delete \
                      --app-id ${APPID} \
                      --group-id ${CHANNEL}
  echo "Coreupdate server cleaned up"
}

function finish() {
  echo "Finishing the test"

  trap shutdown EXIT

  cleanup_coreupdate
  upload_logs tectonic-channel-operator
  upload_logs kube-version-operator

  shutdown
}

function upload_payloads() {
  echo "Uploading payload, assuming the payload (update-payload/payload.json) is up-to-date"

  ${ROOT}/update-payload/upload-payload.sh

  echo "Creating channel ${CHANNEL}"
  ${UPDATESERVICECTL} --server ${SERVER} \
                      --key ${COREUPDATE_KEY} \
                      --user ${COREUPDATE_USR} \
                      channel create \
                      --app-id ${APPID} \
                      --channel ${CHANNEL} \
                      --version ${VERSION}

  echo "Creating group ${CHANNEL}"
  ${UPDATESERVICECTL} --server ${SERVER} \
                      --key ${COREUPDATE_KEY} \
                      --user ${COREUPDATE_USR} \
                      group create \
                      --app-id ${APPID} \
                      --channel ${CHANNEL} \
                      --group-id ${CHANNEL}

  echo "Publishing payload to channel ${CHANNEL}"
  ${ROOT}/update-payload/publish-payload.sh ${CHANNEL}
  echo "Payload successfully published to channel ${CHANNEL}"

  trap finish EXIT
}

function bringup_cluster() {
  echo "Downloading kubectl"

  ${ROOT}/installer/scripts/bare-metal/get-kubectl $(pwd)
  curl -O https://storage.googleapis.com/kubernetes-release/release/v1.5.2/bin/linux/amd64/kubectl && chmod +x ${ROOT}/kubectl
  export KUBECTL=$(pwd)/kubectl

  echo "Fetching the release tarball"

  INSTALLER_DIR=${ROOT}/bin
  URL="https://releases.tectonic.com/tectonic-${ORIGINAL_VERSION}.tar.gz"
  curl -L -O ${URL}
  mkdir -p ${INSTALLER_DIR}
  tar -C ${INSTALLER_DIR} --strip-components=4 -xzf tectonic-${ORIGINAL_VERSION}.tar.gz ./tectonic/tectonic-installer/linux/installer
  export INSTALLER_BIN=${INSTALLER_DIR}/installer
  chmod +x ${INSTALLER_BIN}

  echo "Launching cluster"

  result=$(git checkout -b current)
  if [[ $? != 0 ]]; then
      git branch -D current
      git checkout -b current
  fi
  git checkout tags/${ORIGINAL_VERSION}

  # ---- Start of up.sh ----

  echo "setup environment"
  source "${ROOT}/tests/scripts/aws/default.env.sh"

  echo "Publishing payload to staging channel"
  ${ROOT}/scripts/update-payload/publish-payload.sh "tectonic-1.5-staging"
  echo "Payload successfully published to tectonic-1.5-staging channel"

  echo "wait for machines to come up"
  sleep 200
  ${ROOT}/tests/scripts/aws/wait_for_dns.sh

  echo "transfer config and start bootkube"
  ${ROOT}/tests/scripts/aws/setup_bootkube.sh

  echo "Bootkube running!"

  # ---- End of up.sh ----

  source ${ROOT}/scripts/awsutil.sh

  git checkout current

  echo "Configure kubectl"
  CLUSTER_DIR=${CLUSTER_DIR:-"${ROOT}/tests/scripts/aws/output/${CLUSTER_NAME}"}
  export KUBECONFIG="${CLUSTER_DIR}/assets/auth/kubeconfig"

  echo "Kubeconfig:"
  cat ${KUBECONFIG}

  echo "Printing debug information:"
  ${KUBECTL} get pods --all-namespaces
  ${KUBECTL} get nodes
}

# main:
# The upgrade test will be triggered by commenting "test upgrade from xxx",
# where "xxx" is the start version, e.g. "1.5.2-tectonic.2".
# The jenkins job can be found at https://jenkins-tectonic.prod.coreos.systems/job/tectonic-upgrade-test/

# ${ghprbTriggerCommentBody} should be in the form of:
# "test upgrade from 1.5.2-tectonic.1
export DOCKER_CONFIG=$(dirname ${DOCKER_CONFIG_FILE})
export ROOT=$(git rev-parse --show-toplevel)
export CLUSTER_NAME=${CLUSTER_NAME:-"upgrade-test-${BUILD_NUMBER}"}

arr=(${ghprbCommentBody})
payload=${ROOT}/update-payload/payload.json
export ORIGINAL_VERSION=${ORIGINAL_VERSION:-${arr[3]}}
export TARGET_VERSION=$(cat ${payload} | jq -r .version)
export KUBERNETES_TARGET_VERSION=$(cat ${payload} | jq -r '.desiredVersions[] | select(.name | contains("kubernetes")) |.version')

export AWS_REGION="us-west-2"
export AWS_DEFAULT_REGION=${AWS_REGION}
export SERVER="https://public-update.staging.core-os.net"
export APPID="1f650cdd-878d-4550-9013-6786be951696"
export CHANNEL="upgrade-test-${BUILD_NUMBER}"
# The "version" on the coreupdate channel needs to be big enough
# so the TCO can receive it with the ${ORIGINAL_VERSION}.
export VERSION="${TARGET_VERSION}${BUILD_NUMBER}"

export BUCKET="tectonic-update-payload-testing"
export DESTINATION="${CLUSTER_NAME}.json"

export UPDATESERVICECTL="docker run -v ${ROOT}:${ROOT} quay.io/coreos/updateservicectl@sha256:471ef9637c16df765c4fd928e4c93f011e471bae7e4191022cc6a074f9b4b628"

echo "Updating from ${ORIGINAL_VERSION} to ${TARGET_VERSION}"

upload_payloads

bringup_cluster

sudo rkt run \
    --volume bk,kind=host,source=${ROOT} \
    --mount volume=bk,target=/go/src/github.com/coreos/tectonic-installer \
    --insecure-options=image docker://golang:1.7.5 \
    --user=$(id -u) \
    --set-env=KUBECONFIG_PATH="${ROOT}/tests/scripts/aws/output/${CLUSTER_NAME}/assets/auth/kubeconfig" \
    --set-env=ORIGINAL_VERSION=${ORIGINAL_VERSION} \
    --set-env=CHANNEL=${CHANNEL} \
    --set-env=TECTONIC_TARGET_VERSION=${TARGET_VERSION} \
    --set-env=KUBERNETES_TARGET_VERSION=${KUBERNETES_TARGET_VERSION} \
    --exec /bin/bash -- -c \
    "cd /go/src/github.com/coreos/tectonic-installer && go test -timeout 30m -v ./installer/tests/upgrade"
