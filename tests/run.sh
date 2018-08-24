#!/usr/bin/env bash
#shellcheck disable=SC2155

set -e

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
  echo -e "\\e[34m Exiting... Destroying Tectonic...\\e[0m"
  tectonic destroy --dir="${CLUSTER_NAME}"
  echo -e "\\e[36m Finished! Smoke test output:\\e[0m ${SMOKE_TEST_OUTPUT}"
  echo -e "\\e[34m So Long, and Thanks for All the Fish\\e[0m"
}

trap destroy EXIT

echo -e "\\e[36m Starting build process...\\e[0m"
bazel build tarball smoke_tests
# In future bazel build could be extracted to another job which could be running in docker container like this:
# docker run --rm -v $PWD:$PWD:Z -w $PWD quay.io/coreos/tectonic-builder:bazel-v0.3 bazel build tarball smoke_tests

echo -e "\\e[36m Unpacking artifacts...\\e[0m"
tar -zxf bazel-bin/tectonic-dev.tar.gz
cp bazel-bin/tests/smoke/linux_amd64_stripped/go_default_test tectonic-dev/smoke
export PATH="$(pwd)/tectonic-dev/installer:${PATH}"
cd tectonic-dev

echo -e "\\e[36m Creating Tectonic configuration...\\e[0m"
python <<-EOF >"${CLUSTER_NAME}.yaml"
	import sys
	import yaml

	with open('examples/tectonic.aws.yaml') as f:
	    config = yaml.load(f)
	config['name'] = '${CLUSTER_NAME}'
	config['baseDomain'] = '${DOMAIN}'
	config['licensePath'] = '${LICENSE_PATH}'
	config['pullSecretPath'] = '${PULL_SECRET_PATH}'
	config['aws']['region'] = '${AWS_REGION}'
	config['aws']['master']['iamRoleName'] = 'tf-tectonic-master-node'
	config['aws']['worker']['iamRoleName'] = 'tf-tectonic-worker-node'
	yaml.safe_dump(config, sys.stdout)
	EOF

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
if [ ! -f ~/.ssh/id_rsa.pub ]; then
  echo -e "\\e[36m Generating SSH key-pair...\\e[0m"
  ssh-keygen -qb 2048 -t rsa -f ~/.ssh/id_rsa -N "" </dev/zero
fi
export TF_VAR_tectonic_admin_ssh_key="$(cat ~/.ssh/id_rsa.pub)"

echo -e "\\e[36m Deploying Tectonic...\\e[0m"
tectonic install --dir="${CLUSTER_NAME}"
echo -e "\\e[36m Running smoke test...\\e[0m"
export SMOKE_KUBECONFIG="$(pwd)/$CLUSTER_NAME/generated/auth/kubeconfig"
export SMOKE_NODE_COUNT="5"  # Sum of all nodes (master + worker)
export SMOKE_MANIFEST_PATHS="$(pwd)/$CLUSTER_NAME/generated"
exec 5>&1
SMOKE_TEST_OUTPUT=$(./smoke -test.v --cluster | tee >(cat - >&5))
