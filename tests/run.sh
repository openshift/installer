#!/bin/bash -e
#shellcheck disable=SC2155

# This should be executed from top-level directory not from `tests` directory
# Script needs two variables to be set before execution
# 1) LICENSE_PATH - path to tectonic license file
# 2) PULL_SECRET_PATH - path to pull secret file

set -eo pipefail

SMOKE_TEST_OUTPUT="Never executed. Problem with one of previous stages"
[ -z ${LICENSE_PATH+x} ] && (echo "Please set LICENSE_PATH"; exit 1)
[ -z ${PULL_SECRET_PATH+x} ] && (echo "Please set PULL_SECRET_PATH"; exit 1)
[ -z ${DOMAIN+x} ] && DOMAIN="tectonic-ci.de"
[ -z ${AWS_REGION+x} ] && AWS_REGION="eu-west-1"
[ -z ${JOB_NAME+x} ] && PREFIX="${USER:-test}" || PREFIX="ci-${JOB_NAME#*/}"
CLUSTER_NAME=$(echo "${PREFIX}-$(uuidgen -r | cut -c1-5)" | tr '[:upper:]' '[:lower:]')
exec &> >(tee -a "$CLUSTER_NAME.log")

function destroy() {
  echo -e "\\e[34m Exiting... Destroying Tectonic and cleaning SSH keys...\\e[0m"
  tectonic destroy --dir="${CLUSTER_NAME}"
  aws ec2 delete-key-pair --key-name "${CLUSTER_NAME}"
  echo -e "\\e[36m Finished! Smoke test output:\\e[0m ${SMOKE_TEST_OUTPUT}"
  echo -e "\\e[34m So Long, and Thanks for All the Fish\\e[0m"
}

trap destroy EXIT

echo -e "\\e[36m Starting build process...\\e[0m"
bazel build tarball tests/smoke
# In future bazel build could be extracted to another job which could be running in docker container like this:
# docker run --rm -v $PWD:$PWD:Z -w $PWD quay.io/coreos/tectonic-builder:bazel-v0.3 bazel build tarball tests/smoke

echo -e "\\e[36m Unpacking artifacts...\\e[0m"
tar -zxf bazel-bin/tectonic-dev.tar.gz
cp bazel-bin/tests/smoke/linux_amd64_stripped/smoke tectonic-dev/smoke
export PATH="$(pwd)/tectonic-dev/installer:${PATH}"
cd tectonic-dev

echo -e "\\e[36m Creating Tectonic configuration...\\e[0m"
CONFIG=$(python -c 'import sys, yaml, json; json.dump(yaml.load(sys.stdin), sys.stdout)' < examples/tectonic.aws.yaml)
CONFIG=$(echo "${CONFIG}" | jq ".name = \"${CLUSTER_NAME}\"" |\
                            jq ".baseDomain = \"${DOMAIN}\"" |\
                            jq ".licensePath = \"${LICENSE_PATH}\"" |\
                            jq ".pullSecretPath = \"${PULL_SECRET_PATH}\"" |\
                            jq ".aws.region = \"${AWS_REGION}\"" |\
                            jq ".aws.master.iamRoleName = \"tf-tectonic-master-node\"" |\
                            jq ".aws.worker.iamRoleName = \"tf-tectonic-worker-node\"" |\
                            jq ".aws.etcd.iamRoleName = \"tf-tectonic-etcd-node\""
)
echo "${CONFIG}" | python -c 'import sys, yaml, json; yaml.safe_dump(json.load(sys.stdin), sys.stdout)' > "${CLUSTER_NAME}.yaml"

echo -e "\\e[36m Initializing Tectonic...\\e[0m"
tectonic init --config="${CLUSTER_NAME}".yaml

### ASSUME ROLE ###
echo -e "\\e[36m Setting up AWS credentials...\\e[0m"
export AWS_DEFAULT_REGION="${AWS_REGION}"
unset AWS_SESSION_TOKEN
ACCOUNT_ID=$(aws sts get-caller-identity | jq --raw-output '.Account')
ROLE_ARN="arn:aws:iam::${ACCOUNT_ID}:role/tf-tectonic-installer"
RES=$(aws sts assume-role --role-arn="${ROLE_ARN}" --role-session-name="jenkins-${CLUSTER_NAME}")
export AWS_SECRET_ACCESS_KEY=$(echo "${RES}" | jq --raw-output '.Credentials.SecretAccessKey')
export AWS_ACCESS_KEY_ID=$(echo  "${RES}" | jq --raw-output '.Credentials.AccessKeyId')
export AWS_SESSION_TOKEN=$(echo "${RES}" | jq --raw-output '.Credentials.SessionToken')

### HANDLE SSH KEY ###
echo -e "\\e[36m Uploading SSH key-pair to AWS...\\e[0m"
if [ ! -f "$HOME/.ssh/id_rsa.pub" ]; then
  #shellcheck disable=SC2034
  SSH=$(ssh-keygen -b 2048 -t rsa -f "${HOME}/.ssh/id_rsa" -N "" < /dev/zero)
fi
aws ec2 import-key-pair --key-name "${CLUSTER_NAME}" --public-key-material "file://$HOME/.ssh/id_rsa.pub"
export TF_VAR_tectonic_aws_ssh_key="${CLUSTER_NAME}"

echo -e "\\e[36m Deploying Tectonic...\\e[0m"
tectonic install --dir="${CLUSTER_NAME}"
echo -e "\\e[36m Running smoke test...\\e[0m"
export SMOKE_KUBECONFIG="$(pwd)/$CLUSTER_NAME/generated/auth/kubeconfig"
export SMOKE_NETWORKING="canal"
export SMOKE_NODE_COUNT="7"  # Sum of all nodes (etcd + master + worker)
export SMOKE_MANIFEST_PATHS="$(pwd)/$CLUSTER_NAME/generated"
exec 5>&1
SMOKE_TEST_OUTPUT=$(./smoke -test.v --cluster | tee >(cat - >&5))
