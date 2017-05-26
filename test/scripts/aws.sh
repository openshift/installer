#!/bin/bash -ex
set -o pipefail
shopt -s expand_aliases

# Set the specified vars file
TF_VARS_FILE=$1
TEST_NAME=$(echo ${TF_VARS_FILE} | cut -d "." -f 1)

# Set required configuration
export PLATFORM=aws
export CLUSTER="${TEST_NAME}-${BRANCH_NAME}-${BUILD_ID}"

# s3 buckets require lowercase names
export TF_VAR_tectonic_cluster_name=$(echo ${CLUSTER} | awk '{print tolower($0)}')

# randomly select region
REGIONS=(us-east-1 us-east-2 us-west-1 us-west-2)
export CHANGE_ID=${CHANGE_ID:-${BUILD_ID}}
i=$(( ${CHANGE_ID} % ${#REGIONS[@]} ))
export TF_VAR_tectonic_aws_region="${REGIONS[$i]}"
export AWS_REGION="${REGIONS[$i]}"
echo "selected region: ${TF_VAR_tectonic_aws_region}"
# make core utils accessible to make
export PATH=/bin:${PATH}

# Create local config
make localconfig

# Use smoke test configuration for deployment
ln -sf ${WORKSPACE}/test/${TF_VARS_FILE} ${WORKSPACE}/build/${CLUSTER}/terraform.tfvars

alias filter=${WORKSPACE}/installer/scripts/filter.sh

make plan | filter
if [ "$ONLY_PLAN" = true ]; then
    exit 0
fi
make apply | filter

# TODO: replace in Go
CONFIG=${WORKSPACE}/build/${CLUSTER}/terraform.tfvars
MASTER_COUNT=$(grep tectonic_master_count ${CONFIG} | awk -F "=" '{gsub(/"/, "", $2); print $2}')
WORKER_COUNT=$(grep tectonic_worker_count ${CONFIG} | awk -F "=" '{gsub(/"/, "", $2); print $2}')

export NODE_COUNT=$(( ${MASTER_COUNT} + ${WORKER_COUNT} ))

export TEST_KUBECONFIG=${WORKSPACE}/build/${CLUSTER}/generated/auth/kubeconfig
installer/bin/sanity -test.v -test.parallel=1
