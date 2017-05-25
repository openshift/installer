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

echo "Destroying ${CLUSTER}..."
make destroy
