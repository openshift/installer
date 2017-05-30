#!/bin/bash -ex
set -o pipefail
shopt -s expand_aliases

DIR="$( cd "$( dirname "$0" )" && pwd )"
# shellcheck disable=SC1090
source "$DIR/aws-common.sh"

# make core utils accessible to make
export PATH=/bin:${PATH}

# Create local config
make localconfig

# Use smoke test configuration for deployment
ln -sf "${WORKSPACE}/test/${TF_VARS_FILE}" "${WORKSPACE}/build/${CLUSTER}/terraform.tfvars"

# shellcheck disable=SC2139
alias filter="${WORKSPACE}/installer/scripts/filter.sh"

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
