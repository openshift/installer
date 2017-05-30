#!/bin/bash -ex
set -o pipefail
shopt -s expand_aliases

DIR="$( cd "$( dirname "$0" )" && pwd )"
# make core utils accessible to make
export PATH=/bin:$PATH
 # Alias filter for convenience
# shellcheck disable=SC2139
alias filter="$WORKSPACE"/installer/scripts/filter.sh
export PLATFORM=aws

common() {
    # Set the specified vars file
    TF_VARS_FILE=$1
    TEST_NAME=$(basename "$TF_VARS_FILE" | cut -d "." -f 1)
    
    # Set required configuration
    CLUSTER="$TEST_NAME-$BRANCH_NAME-$BUILD_ID"
    MAX_LENGTH=19
    
    LENGTH=${#CLUSTER}
    if [ "$LENGTH" -gt "$MAX_LENGTH" ]
    then
        CLUSTER="${CLUSTER:0:MAX_LENGTH}"
        echo "Cluster name too long. Truncated to $CLUSTER"
    elif [ "$LENGTH" -lt "$MAX_LENGTH" ]
    then
        APPEND=$(( MAX_LENGTH - LENGTH ))
        APPEND_STR="012345678901234567890123456789"
        CLUSTER="$CLUSTER${APPEND_STR:0:APPEND}"
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
    echo "selected region: $TF_VAR_tectonic_aws_region"
    echo "cluster name: $CLUSTER"

    # Create local config
    make localconfig
    # Use smoke test configuration for deployment
    ln -sf "$DIR/$TF_VARS_FILE" "$WORKSPACE/build/$CLUSTER/terraform.tfvars"
}

create() {
    common "$1"
    make plan | filter
    make apply | filter
}

destroy() {
    common "$1"
    echo "Destroying ${CLUSTER}..."
    make destroy
}

plan() {
    common "$1"
    make plan | filter
}

test_cluster() {
    common "$1"
    # TODO: replace in Go
    CONFIG=$WORKSPACE/build/$CLUSTER/terraform.tfvars
    MASTER_COUNT=$(grep tectonic_master_count "$CONFIG" | awk -F "=" '{gsub(/"/, "", $2); print $2}')
    WORKER_COUNT=$(grep tectonic_worker_count "$CONFIG" | awk -F "=" '{gsub(/"/, "", $2); print $2}')
    export NODE_COUNT=$(( MASTER_COUNT + WORKER_COUNT ))
    export TEST_KUBECONFIG=$WORKSPACE/build/$CLUSTER/generated/auth/kubeconfig
    installer/bin/sanity -test.v -test.parallel=1
}

usage() {
    set +x
    printf "%s is a tool for running Tectonic smoke tests on AWS.\n\n" "$(basename "$0")"
    printf "Usage:\n\n \t %s command [arguments]\n\n" "$(basename "$0")"
    printf "The commands are:\n\n"
    printf "\t create <tfvars>  \tcreate a Tectonic cluster parameterized by <tfvars>\n"
    printf "\t destroy <tfvars> \tdestroy the Tectonic cluster parameterized by <tfvars>\n"
    printf "\t plan <tfvars>    \tplan a Tectonic cluster parameterized by <tfvars>\n"
    printf "\t test <tfvars>    \ttest a Tectonic cluster parameterized by <tfvars>\n"
    printf "\n"
}

COMMAND=$1
if [ $# -eq 0 ]; then
    usage
    exit 1
fi

shift
case $COMMAND in
    create)
        create "$@";;
    destroy)
        destroy "$@";;
    plan)
        plan "$@";;
    test)
        test_cluster "$@";;
    *)
        usage;;
esac
