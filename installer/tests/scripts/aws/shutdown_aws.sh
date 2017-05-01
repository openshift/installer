#!/bin/bash -e
# This script destroys the current cluster.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/../../.."

source "${ROOT}/scripts/awsutil.sh"

# require AWS creds
check_aws_creds

RKT_STAGE1_IMG=${RKT_STAGE1_IMG:-"coreos.com/rkt/stage1-fly:1.17.0"}
RKT_OPTS="--interactive --stage1-name=${RKT_STAGE1_IMG} --insecure-options=image"
sudo rkt run ${RKT_OPTS} --set-env=AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} --set-env=AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} docker://python:2.7.12-slim --exec /bin/bash <<EOF
    pip install awscli
    export AWS_DEFAULT_REGION=${AWS_REGION}
    # For debugging, print info before removing the cluster
    aws cloudformation describe-stacks --stack-name=${CLUSTER_NAME}
    aws cloudformation describe-stack-resources --stack-name=${CLUSTER_NAME}
    aws cloudformation delete-stack --stack-name=${CLUSTER_NAME}
EOF
