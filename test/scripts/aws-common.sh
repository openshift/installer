#!/bin/bash -ex
set -o pipefail
shopt -s expand_aliases

# Set the specified vars file
export TF_VARS_FILE=$1
TEST_NAME=$(echo "${TF_VARS_FILE}" | cut -d "." -f 1)

# Set required configuration
export PLATFORM=aws
CLUSTER="${TEST_NAME}-${BRANCH_NAME}-${BUILD_ID}"
MAX_LENGTH=19

LENGTH=${#CLUSTER}
if [ "$LENGTH" -gt "$MAX_LENGTH" ]
then
  CLUSTER="${CLUSTER:0:MAX_LENGTH}"
  echo "Cluster name too long. Truncated to $CLUSTER"
elif [ "$LENGTH" -lt "$MAX_LENGTH" ]
then
  APPEND=$(( MAX_LENGTH - LENGTH ))
  APPEND_STR="01234567890123456789123456789"
  CLUSTER="${CLUSTER}${APPEND_STR:0:APPEND}"
  echo "Cluster name too short. Appended to $CLUSTER"
fi

CLUSTER=$(echo "${CLUSTER}" | awk '{print tolower($0)}')
export CLUSTER
export TF_VAR_tectonic_cluster_name=$CLUSTER

# randomly select region
REGIONS=(us-east-1 us-east-2 us-west-1 us-west-2)
export CHANGE_ID=${CHANGE_ID:-${BUILD_ID}}
i=$(( CHANGE_ID % ${#REGIONS[@]} ))
export TF_VAR_tectonic_aws_region="${REGIONS[$i]}"
export AWS_REGION="${REGIONS[$i]}"
echo "selected region: ${TF_VAR_tectonic_aws_region}"
echo "cluster name: ${CLUSTER}"
