# OpenShift development notes | BYO VPC with Multi CDIR

This guides describes how to create the network requirements
to deploy OpenShift on AWS with BYO VPC with multi CIDR blocks,
and the Machine CIDR is using non-default from VPC.

Templates used in this document:

## Prerequisites

- AWS permissions to deploy using UPI
- Additional permissions
- Considering you are under the root directory of installer repository
- [`aws cli`](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) installed

## Steps

- Clone the installer repository

- Declare the variables used in the following scripts

```sh
CLUSTER_NAME=byvpc-cidr

INSTALL_DIR=/tmp/${CLUSTER_NAME}

PULL_SECRET_FILE="${HOME}/pull-secret-latest.json"
CLUSTER_BASE_DOMAIN=devcluster.mydomain.com
SSH_PUB_KEY_FILE=$HOME/.ssh/id_rsa.pub

AWS_REGION=us-east-1

VPC_CIDR_SECONDARY="10.100.0.0/16"
```

- Create the CloudFormation stack, setting the parameter `VpcCidr2` with custom CIDR Block (or leave the default):
    - The file must be available under your path: `${PWD}/upi/aws/cloudformation-templates/byo-vpc-multi-cidr.yaml`
    - To track the installation visit the [AWS CloudFormation Console](https://console.aws.amazon.com/cloudformation)

```sh
STACK_VPC="${CLUSTER_NAME}-vpc"
aws cloudformation create-stack --region ${AWS_REGION}  --stack-name ${STACK_VPC} \
  --template-body file://${PWD}/upi/aws/cloudformation-templates/byo-vpc-multi-cidr.yaml \
  --parameters ParameterKey=VpcCidr2,ParameterValue=${VPC_CIDR_SECONDARY}

aws --region ${AWS_REGION} cloudformation wait stack-create-complete --stack-name ${STACK_VPC}
aws --region ${AWS_REGION} cloudformation describe-stacks --stack-name "${STACK_VPC}"
```

- Export the subnet IDs:

```sh
# Extract subnet IDs
mapfile -t SUBNET_IDS < <(aws --region ${AWS_REGION} \
    cloudformation describe-stacks --stack-name "${STACK_VPC}" \
    --query "Stacks[0].Outputs[?OutputKey=='SubnetsIdsForCidr2'].OutputValue" \
    --output text | tr ',' '\n')

echo ${SUBNET_IDS[@]}
```

- Create the `install-config.yaml` setting the:
    - `networking.machineNetwork.[0].cidr` with secondary VPC CIDR block
    - `platform.aws.subnets`: with subnet IDs created by CloudFormation Stack

```sh
mkdir -p ${INSTALL_DIR}
cat <<EOF | envsubst > ${INSTALL_DIR}/install-config.yaml
apiVersion: v1
baseDomain: ${CLUSTER_BASE_DOMAIN}
metadata:
  name: "${CLUSTER_NAME}"
platform:
  aws:
    region: ${AWS_REGION}
    subnets:
$(for SB in ${SUBNET_IDS[*]}; do echo "    - $SB"; done)
    userTags:
      x-red-hat-clustertype: installer
      x-red-hat-managed: "true"
networking:
  machineNetwork:
  - cidr: ${VPC_CIDR_SECONDARY}
publish: External
pullSecret: '$(cat ${PULL_SECRET_FILE} |awk -v ORS= -v OFS= '{$1=$1}1')'
sshKey: |
  $(cat ${SSH_PUB_KEY_FILE})
EOF
```

- Ensure the config `${INSTALL_DIR}/install-config.yaml` is correct and create the cluster:

```sh
openshift-install create cluster --dir "${INSTALL_DIR}" --log-level=debug
```
