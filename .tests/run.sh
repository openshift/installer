#!/usr/bin/env bash
#shellcheck disable=SC2155

set -e

# This should be executed from top-level directory not from `tests` directory
# Script needs one variable to be set before execution
# 1) PULL_SECRET_PATH - path to pull secret file

set -eo pipefail

BACKEND="${1}"
LEAVE_RUNNING="${LEAVE_RUNNING:-n}"  # do not teardown after successful initialization
SMOKE_TEST_OUTPUT="Never executed. Problem with one of previous stages"
[ -z ${PULL_SECRET_PATH+x} ] && (echo "Please set PULL_SECRET_PATH"; exit 1)
[ -z ${DOMAIN+x} ] && DOMAIN="tectonic-ci.de"
[ -z ${JOB_NAME+x} ] && PREFIX="${USER:-test}" || PREFIX="ci-${JOB_NAME#*/}"
CLUSTER_NAME=$(echo "${PREFIX}-$(uuidgen -r | cut -c1-5)" | tr '[:upper:]' '[:lower:]')
TECTONIC="${PWD}/tectonic-dev/installer/tectonic"
exec &> >(tee -ai "$CLUSTER_NAME.log")

function destroy() {
  echo -e "\\e[34m Exiting... Destroying Tectonic...\\e[0m"
  "${TECTONIC}" destroy --dir="${CLUSTER_NAME}" --continue-on-error
  echo -e "\\e[36m Finished! Smoke test output:\\e[0m ${SMOKE_TEST_OUTPUT}"
  echo -e "\\e[34m So Long, and Thanks for All the Fish\\e[0m"
}

echo -e "\\e[36m Starting build process...\\e[0m"
bazel build tarball smoke_tests
# In future bazel build could be extracted to another job which could be running in docker container like this:
# docker run --rm -v $PWD:$PWD:Z -w $PWD quay.io/coreos/tectonic-builder:bazel-v0.3 bazel build tarball smoke_tests

echo -e "\\e[36m Unpacking artifacts...\\e[0m"
tar -zxf bazel-bin/tectonic-dev.tar.gz
cp bazel-bin/tests/smoke/linux_amd64_pure_stripped/go_default_test tectonic-dev/smoke
chmod 755 tectonic-dev/smoke
cd tectonic-dev

### HANDLE SSH KEY ###
if [ ! -f ~/.ssh/id_rsa.pub ]; then
  echo -e "\\e[36m Generating SSH key-pair...\\e[0m"
  ssh-keygen -qb 2048 -t rsa -f ~/.ssh/id_rsa -N "" </dev/zero
fi

case "${BACKEND}" in
aws)
  if test -z "${AWS_REGION+x}"
  then
    echo -e "\\e[36m Calculating the AWS region...\\e[0m"
    AWS_REGION="$(aws configure get region)" ||
    AWS_REGION="${AWS_REGION:-us-east-1}"
  fi
  export AWS_DEFAULT_REGION="${AWS_REGION}"
  unset AWS_SESSION_TOKEN

  ### ASSUME ROLE ###
  echo -e "\\e[36m Setting up AWS credentials...\\e[0m"
  ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
  ROLE_ARN="arn:aws:iam::${ACCOUNT_ID}:role/tf-tectonic-installer"
  RES="$(aws sts assume-role --role-arn="${ROLE_ARN}" --role-session-name="jenkins-${CLUSTER_NAME}" --query Credentials --output json)" &&
  export AWS_SECRET_ACCESS_KEY="$(echo "${RES}" | jq --raw-output '.SecretAccessKey')" &&
  export AWS_ACCESS_KEY_ID="$(echo  "${RES}" | jq --raw-output '.AccessKeyId')" &&
  export AWS_SESSION_TOKEN="$(echo "${RES}" | jq --raw-output '.SessionToken')" &&
  CONFIGURE_AWS_ROLES=True ||
  CONFIGURE_AWS_ROLES=False
  ;;
libvirt)
  ;;
*)
  echo "unrecognized backend: ${BACKEND}" >&2
  echo "Use ${0} BACKEND, where BACKEND is aws or libvirt" >&2
  exit 1
esac

echo -e "\\e[36m Creating Tectonic configuration...\\e[0m"
python <<-EOF >"${CLUSTER_NAME}.yaml"
	import datetime
	import os.path
	import sys

	import yaml

	with open('examples/${BACKEND}.yaml') as f:
	    config = yaml.load(f)
	config['name'] = '${CLUSTER_NAME}'
	with open(os.path.expanduser(os.path.join('~', '.ssh', 'id_rsa.pub'))) as f:
	    config['admin']['sshKey'] = f.read()
	config['baseDomain'] = '${DOMAIN}'
	with open('${PULL_SECRET_PATH}') as f:
	    config['pullSecret'] = f.read()
	if '${BACKEND}' == 'aws':
	    config['aws']['region'] = '${AWS_REGION}'
	    config['aws']['extraTags'] = {
	        'expirationDate': (
	            datetime.datetime.utcnow() + datetime.timedelta(hours=4)
	        ).strftime('%Y-%m-%dT%H:%M+0000'),
	    }
	    if ${CONFIGURE_AWS_ROLES:-False}:
	        config['aws']['master']['iamRoleName'] = 'tf-tectonic-master-node'
	        config['aws']['worker']['iamRoleName'] = 'tf-tectonic-worker-node'
	elif '${BACKEND}' == 'libvirt' and '${IMAGE_URL}':
	    config['libvirt']['image'] = '${IMAGE_URL}'
	yaml.safe_dump(config, sys.stdout)
	EOF

echo -e "\\e[36m Initializing Tectonic...\\e[0m"
"${TECTONIC}" init --config="${CLUSTER_NAME}".yaml

trap destroy EXIT

echo -e "\\e[36m Deploying Tectonic...\\e[0m"
"${TECTONIC}" install --dir="${CLUSTER_NAME}"
echo -e "\\e[36m Running smoke test...\\e[0m"
export SMOKE_KUBECONFIG="${PWD}/${CLUSTER_NAME}/generated/auth/kubeconfig"
export SMOKE_MANIFEST_PATHS="${PWD}/${CLUSTER_NAME}/generated"
case "${BACKEND}" in
aws)
  export SMOKE_NODE_COUNT=5  # Sum of all nodes (boostrap + master + worker)
  ;;
libvirt)
  export SMOKE_NODE_COUNT=4
  ;;
esac
exec 5>&1
if test "${LEAVE_RUNNING}" = y; then
  echo "leaving running; tear down manually with: cd ${PWD} && installer/tectonic destroy --dir=${CLUSTER_NAME}"
  trap - EXIT
fi
SMOKE_TEST_OUTPUT=$(./smoke -test.v --cluster | tee -i >(cat - >&5))
