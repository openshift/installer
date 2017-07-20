#!/bin/bash -ex
set -o pipefail
shopt -s expand_aliases
DIR="$(cd "$(dirname "$0")" && pwd)"
WORKSPACE=${WORKSPACE:-"$(cd "$DIR"/../../.. && pwd)"}
# Alias filter for convenience
# shellcheck disable=SC2139
alias filter="$WORKSPACE"/installer/scripts/filter.sh

common() {
    # make core utils accessible to make
    export PATH=/bin:$PATH
    export PLATFORM=azure

    # Set the specified vars file
    TF_VARS_FILE=$1

    TEST_NAME=$(basename "$TF_VARS_FILE" | cut -d "." -f 1)
    MAX_NAME_LENGTH=10

    LENGTH=${#TEST_NAME}
    if [ "$LENGTH" -gt "$MAX_NAME_LENGTH" ]
    then
        TEST_NAME="${TEST_NAME:0:MAX_NAME_LENGTH}"
        echo "TEST_NAME too long. Truncated to $TEST_NAME"
    fi

    # Set required configuration
    CLUSTER="$TEST_NAME-$BRANCH_NAME-$BUILD_ID"
    MAX_LENGTH=28
    
    LENGTH=${#CLUSTER}
    if [ "$LENGTH" -gt "$MAX_LENGTH" ]
    then
        CLUSTER="${CLUSTER:0:MAX_LENGTH}"
        echo "Cluster name too long. Truncated to $CLUSTER"
    elif [ "$LENGTH" -lt "$MAX_LENGTH" ]
    then
        # Shellcheck v0.4.4 does not detect that 'APPEND' is used two lines
        # further below. This comment can be removed once we (/debian) moves to
        # v0.4.6
        # shellcheck disable=SC2034
        APPEND=$(( MAX_LENGTH - LENGTH ))
        APPEND_STR="012345678901234567890123456789"
        CLUSTER="$CLUSTER${APPEND_STR:0:APPEND}"
        echo "Cluster name too short. Appended to $CLUSTER"
    fi
    
    CLUSTER=$(echo "${CLUSTER}" | awk '{print tolower($0)}')
    export CLUSTER
    export TF_VAR_tectonic_cluster_name=$CLUSTER
    
    echo "cluster name: $CLUSTER"

    mkdir -p "$HOME/.ssh"
    echo "$AZURE_SMOKE_SSH_KEY_PUB" > "$HOME/.ssh/id_rsa.pub"
    echo "$AZURE_SMOKE_SSH_KEY" > "$HOME/.ssh/id_rsa"
    export TF_VAR_tectonic_azure_ssh_key
    TF_VAR_tectonic_azure_ssh_key=$(realpath "$HOME/.ssh/id_rsa.pub")

    export TF_VAR_tectonic_azure_client_secret
    TF_VAR_tectonic_azure_client_secret="${ARM_CLIENT_SECRET}"

    # Create local config
    make localconfig
    # Use smoke test configuration for deployment
    cp "$DIR/$TF_VARS_FILE" "$WORKSPACE/build/$CLUSTER/terraform.tfvars"
}

create() {
    common "$1"
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
    export MANIFEST_PATHS=$WORKSPACE/build/$CLUSTER/generated/
    # shellcheck disable=SC2155
    export MANIFEST_EXPERIMENTAL=$(grep tectonic_experimental "$CONFIG" | awk -F "=" '{gsub(/"/, "", $2); print $2}' | tr -d ' ')
    # shellcheck disable=SC2155
    export CALICO_NETWORK_POLICY=$(grep tectonic_calico_network_policy "$CONFIG" | awk -F "=" '{gsub(/"/, "", $2); print $2}' | tr -d ' ')
    bin/smoke -test.v -test.parallel=1 --cluster
}

usage() {
    # It's annoying to print the debug statement and the output from printf
    set +x
    printf "%s is a tool for running Tectonic smoke tests on AWS.\n\n" "$(basename "$0")"
    printf "Usage:\n\n \t %s command [arguments]\n\n" "$(basename "$0")"
    printf "The commands are:\n\n"
    printf "\t create <tfvars>                              create a Tectonic cluster parameterized by <tfvars>\n"
    printf "\t destroy <tfvars>                             destroy the Tectonic cluster parameterized by <tfvars>\n"
    printf "\t plan <tfvars>                                plan a Tectonic cluster parameterized by <tfvars>\n"
    printf "\t                                              path <policy> and trust policy at file path <trust-policy>\n"
    printf "\t test <tfvars>                                test a Tectonic cluster parameterized by <tfvars>\n"
    printf "\n"
}

main () {
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
}

pushd "$WORKSPACE" > /dev/null
main "$@"
popd > /dev/null
