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

    # Create local config
    make localconfig
    # Use smoke test configuration for deployment
    cp "$DIR/$TF_VARS_FILE" "$WORKSPACE/build/$CLUSTER/terraform.tfvars"
}

create() {
    common "$1"
    make apply | filter
}

common_vpc() {
    random_region
    export TF_VAR_vpc_aws_region="$TF_VAR_tectonic_aws_region"
    # shellcheck disable=SC2155
    export TF_VAR_vpc_name="$(echo "vpc-$BRANCH_NAME-$BUILD_ID" | awk '{print tolower($0)}')"
}

create_vpc() {
    common_vpc
    pushd "$WORKSPACE/contrib/internal-cluster"
    set +x
    # shellcheck disable=SC2155
    export TF_VAR_ovpn_password="$(tr -cd '[:alnum:]' < /dev/urandom | head -c 32 ; echo)"
    export TF_VAR_base_domain="tectonic.dev.coreos.systems"
    set -x
    # Create the vpc.
    terraform apply
    # Get the VPN details.
    # shellcheck disable=SC2155
    local vpn_url="$(terraform output -json | jq -r '.ovpn_url.value')"
    until curl -k -L --silent "$vpn_url" > /dev/null; do
        echo "waiting for vpn access server to become available"
        sleep 5
    done
    set +x
    curl -k -L -u "openvpn:$TF_VAR_ovpn_password" --silent --fail "$(terraform output -json | jq -r '.ovpn_url.value')"/rest/GetUserlogin > vpn.conf
    printf "openvpn\n%s\n" "$TF_VAR_ovpn_password" > vpn_credentials
    set -x
    sed -i 's/auth-user-pass/auth-user-pass vpn_credentials/g' vpn.conf
    # Start the VPN.
    openvpn --config vpn.conf --daemon
    until ping -c 1 8.8.8.8 > /dev/null; do
        echo "waiting for vpn connection to become available"
        sleep 5
    done
    # Use AWS VPC DNS rather than host's.
    # shellcheck disable=SC2155
    local vpc_dns="$(terraform output -json | jq -r '.vpc_dns.value')"
    cp /etc/resolv.conf /etc/resolv.conf.bak
    echo "nameserver $vpc_dns" > /etc/resolv.conf
    echo "nameserver 8.8.8.8"  >> /etc/resolv.conf
    # Export all of the VPC details.
    # shellcheck disable=SC2155,SC2183,SC2046
    {
        export TF_VAR_tectonic_aws_external_private_zone="$(terraform output -json | jq -r '.private_zone_id.value')"
        export TF_VAR_tectonic_aws_external_vpc_id="$(terraform output -json | jq -r '.vpc_id.value')"
        export TF_VAR_tectonic_aws_external_master_subnet_ids="$(printf "[%s, %s]" $(terraform output -json | jq '.subnets.value[0,1]'))"
        export TF_VAR_tectonic_aws_external_worker_subnet_ids="$(printf "[%s, %s]" $(terraform output -json | jq '.subnets.value[2,3]'))"
    }
    popd
}

grafiti_clean() {
  common "$1"

  GRAFITI_VERSION="9f58d4a"

  TMP_GRAFITI_CONFIG_DIR="${WORKSPACE}/tmg_grafiti_config"
  TMP_CONFIG_FILE=$(mktemp -d -p "$TMP_GRAFITI_CONFIG_DIR")
  cat <<EOF > $TMP_CONFIG_FILE
[grafiti]
region = "$AWS_REGION"
EOF

  TMP_TAG_FILE=$(mktemp -d -p "$TMP_GRAFITI_CONFIG_DIR")
  cat <<EOF > $TMP_TAG_FILE
{"TagFilters": [{"Key": "tectonicClusterID","Values": ["$CLUSTER"]}]}
EOF

  echo "Cleaning up ${CLUSTER}..."
  docker run --rm --name "cluster-cleaner-${CLUSTER}" \
    -v "$TMP_GRAFITI_CONFIG_DIR":/tmp/config:z \
    -e CONFIG_FILE="/tmp/config/$(basename "$TMP_CONFIG_FILE")" \
    -e TAG_FILE="/tmp/config/$(basename "$TMP_TAG_FILE")" \
    quay.io/coreos/grafiti:"$GRAFITI_VERSION" \
    ash -c "grafiti --ignore-errors --config \"\$CONFIG_FILE\" delete --silent --all-deps --delete-file \"\$TAG_FILE\""

  rm -rf "$TMP_GRAFITI_CONFIG_DIR"
}

destroy_vpc() {
    common_vpc
    # Restore host DNS settings.
    [ -f /etc/resolv.conf.bak ] && cat /etc/resolv.conf.bak > /etc/resolv.conf && rm /etc/resolv.conf.bak
    pushd "$WORKSPACE/contrib/internal-cluster"
    pkill openvpn || true
    until ping -c 1 8.8.8.8 > /dev/null; do
        echo "waiting for network to become available"
        sleep 5
    done
    terraform destroy --force
    popd
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
        create-vpc)
            create_vpc "$@";;
        destroy)
            destroy "$@";;
        destroy-vpc)
            destroy_vpc "$@";;
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
