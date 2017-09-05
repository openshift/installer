#!/bin/bash
# shellcheck disable=SC2002,SC2015,SC2086,SC2091
set -exuo pipefail
shopt -s expand_aliases
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# Allow smoke.sh to work off jenkins.
WORKSPACE=${WORKSPACE:-"$(cd "$DIR"/../../.. && pwd)"}
BUILD_ID=${BUILD_ID:-1}
BRANCH_NAME=${BRANCH_NAME:-$(git rev-parse --abbrev-ref HEAD)}

# Alias filter for convenience
# shellcheck disable=SC2139
alias filter="$WORKSPACE"/installer/scripts/filter.sh

assume_role() {
    # Don't print out the credentials.
    set +x
    ROLE_NAME=$1
    # Get the actual role ARN. This allows us to invoke the script with friendly arguments.
    # shellcheck disable=SC2155
    ROLE_ARN="$(aws iam get-role --role-name="$ROLE_NAME" | jq -r '.Role.Arn')"
    # shellcheck disable=SC2155
    CREDENTIALS="$(aws sts assume-role --role-arn="$ROLE_ARN" --role-session-name=tectonic-installer | jq '.Credentials')"
    export AWS_ACCESS_KEY_ID
    AWS_ACCESS_KEY_ID=$(echo "$CREDENTIALS" | jq -r '.AccessKeyId')
    export AWS_SECRET_ACCESS_KEY
    AWS_SECRET_ACCESS_KEY=$(echo "$CREDENTIALS" | jq -r '.SecretAccessKey')
    export AWS_SESSION_TOKEN
    AWS_SESSION_TOKEN=$(echo "$CREDENTIALS" | jq -r '.SessionToken')
    set -x
}

set_role() {
    ROLE_NAME=$1
    ROLE_POLICY=$2
    TRUST_POLICY=$3
    # If the role exists, then update the trust policy. Otherwise, create a new role.
    if aws iam get-role --role-name="$ROLE_NAME" > /dev/null 2>&1 ; then
        aws iam update-assume-role-policy --role-name="$ROLE_NAME" --policy-document=file://"$TRUST_POLICY"
    else
        aws iam create-role --role-name="$ROLE_NAME" --assume-role-policy-document=file://"$TRUST_POLICY"
    fi
    aws iam put-role-policy --role-name="$ROLE_NAME" --policy-name="$ROLE_NAME" --policy-document=file://"$ROLE_POLICY"
}

random_region() {
    # randomly select region
    REGIONS=(us-east-1 us-east-2 us-west-1 us-west-2)
    export CHANGE_ID=${CHANGE_ID:-${BUILD_ID}}
    i=$(( CHANGE_ID % ${#REGIONS[@]} ))
    export TF_VAR_tectonic_aws_region="${REGIONS[$i]}"
    export AWS_REGION="${REGIONS[$i]}"
}

common() {
    # make core utils accessible to make
    export PATH=/bin:$PATH
    export PLATFORM=aws

    # Set the specified vars file
    TF_VARS_FILE=$1
    TEST_NAME=$(basename "$TF_VARS_FILE" | cut -d "." -f 1)

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
        # shellcheck disable=SC2034
        APPEND=$(( MAX_LENGTH - LENGTH ))
        APPEND_STR="012345678901234567890123456789"
        CLUSTER="$CLUSTER${APPEND_STR:0:APPEND}"
        echo "Cluster name too short. Appended to $CLUSTER"
    fi

    random_region
    CLUSTER=$(echo "${CLUSTER}" | awk '{print tolower($0)}')
    export CLUSTER
    export TF_VAR_tectonic_cluster_name=$CLUSTER

    echo "selected region: $TF_VAR_tectonic_aws_region"
    echo "cluster name: $CLUSTER"

    export CONFIG="${WORKSPACE}/build/${CLUSTER}/terraform.tfvars"
    if [ ! -f "$CONFIG" ]; then
      # Create local config
      make localconfig
      # Use smoke test configuration for deployment
      cp "$DIR/$TF_VARS_FILE" "$CONFIG"
      # Store AWS region in tfvars file
      echo -e "tectonic_aws_region = \"$TF_VAR_tectonic_aws_region\"" >> "$CONFIG"
    fi
}

create() {
    common "$1"
    make apply | filter
}

grafiti_clean() {
    common "$1"

    # Parse AWS_REGION from a tfvars file
    # shellcheck disable=SC2155
    export AWS_REGION="$(grep -oP 'tectonic_aws_region\s*=\s*"\K([^"]+)' "$CONFIG")"

    # Get CLUSTER_ID from `terraform apply` output. Because 'tectonicClusterID'
    # is created during a 'terraform apply' step, which might fail before a state
    # file is generated, we must collect and parse log output
    if [ ! -f "${WORKSPACE}/terraform.log" ]; then
        echo "Cannot find terraform log file"
        exit
    fi
    CLUSTER_ID="$(grep -m 1 -oP 'tags.tectonicClusterID:\s*""[^"]*"\K([^"]+)' "${WORKSPACE}/terraform.log")"

    GRAFITI_TMP_DIR="$(mktemp -d -p $WORKSPACE)"

    GRAFITI_CONFIG_FILE="${GRAFITI_TMP_DIR}/config.toml"
    echo "maxNumRequestRetries = 11" > "$GRAFITI_CONFIG_FILE"

    GRAFITI_TAG_FILE="${GRAFITI_TMP_DIR}/tag.json"
    echo -e "{\"TagFilters\":[{\"Key\":\"tectonicClusterID\",\"Values\":[\"$CLUSTER_ID\"]}]}"  > "$GRAFITI_TAG_FILE"

    echo "Cleaning up \"${CLUSTER_ID}\"..."

    grafiti --config "$GRAFITI_CONFIG_FILE" --ignore-errors delete --silent --all-deps --delete-file "$GRAFITI_TAG_FILE"

    rm -rf "$GRAFITI_TMP_DIR"
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
    MASTER_COUNT=$(grep tectonic_master_count "$CONFIG" | awk -F "=" '{gsub(/"/, "", $2); print $2}')
    WORKER_COUNT=$(grep tectonic_worker_count "$CONFIG" | awk -F "=" '{gsub(/"/, "", $2); print $2}')
    export SMOKE_NODE_COUNT=$(( MASTER_COUNT + WORKER_COUNT ))
    export SMOKE_KUBECONFIG=$WORKSPACE/build/$CLUSTER/generated/auth/kubeconfig
    export SMOKE_MANIFEST_PATHS=$WORKSPACE/build/$CLUSTER/generated/
    # shellcheck disable=SC2155
    export SMOKE_MANIFEST_EXPERIMENTAL=$(grep tectonic_experimental "$CONFIG" | awk -F "=" '{gsub(/"/, "", $2); print $2}' | tr -d ' ')
    # shellcheck disable=SC2155
    export SMOKE_CALICO_NETWORK_POLICY=$(grep tectonic_calico_network_policy "$CONFIG" | awk -F "=" '{gsub(/"/, "", $2); print $2}' | tr -d ' ')
    bin/smoke -test.v -test.parallel=1 --cluster
}

usage() {
    # It's annoying to print the debug statement and the output from printf
    set +x
    printf "%s is a tool for running Tectonic smoke tests on AWS.\n\n" "$(basename "$0")"
    printf "Usage:\n\n \t %s command [arguments]\n\n" "$(basename "$0")"
    printf "The commands are:\n\n"
    printf "\t assume-role <role-name>                      assume the role specified by <role-name>\n"
    printf "\t create <tfvars>                              create a Tectonic cluster parameterized by <tfvars>\n"
    printf "\t destroy <tfvars>                             destroy the Tectonic cluster parameterized by <tfvars>\n"
    printf "\t plan <tfvars>                                plan a Tectonic cluster parameterized by <tfvars>\n"
    printf "\t set-role <role-name> <policy> <trust-policy> create or update the <role-name> role with policy at file\n"
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
        common)
            common "$@";;
        assume-role)
            assume_role "$@";;
        create)
            create "$@";;
        destroy)
            destroy "$@";;
        grafiti-clean)
            grafiti_clean "$@";;
        plan)
            plan "$@";;
        set-role)
            set_role "$@";;
        test)
            test_cluster "$@";;
        *)
            usage;;
    esac
}

pushd "$WORKSPACE" > /dev/null
main "$@"
popd > /dev/null
