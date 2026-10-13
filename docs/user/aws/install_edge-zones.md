# Install a cluster in AWS extending nodes to AWS Local Zones

The steps below describe how to install a cluster in AWS extending worker nodes to Local Zones.

This document is split into the following sections:

- Prerequisites
- [Local Zones](#local-zones)
  - [Install a cluster extending nodes to the Local Zone [new VPC]](#ipi-localzones) (4.14+)
  - [Install a cluster into existing VPC with Local Zone subnets](#ipi-localzones-existing-vpc) (4.13+)
  - [Extend worker nodes to AWS Local Zones in existing clusters [Day 2]](#day2-localzones)
- [Wavelength Zones](#wavelength-zones)
  - [Install a cluster extending nodes to the Wavelength Zone [new VPC]](#ipi-wavelength-auto)
  - [Install a cluster on AWS in existing VPC with subnets in Wavelength Zone](#ipi-wavelength-byovpc)
- [Use Cases](#use-cases)

## Prerequisites for edge zones

### Additional IAM permissions <a name="pre-iam-opt-in"></a>

The AWS Local Zone deployment described in this document requires additional permission from the user creating the cluster allowing Local Zone group modification: `ec2:ModifyAvailabilityZoneGroup`

Example of the permissive IAM Policy that can be attached to the User or Role:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:ModifyAvailabilityZoneGroup"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```

___
___

# Local Zones

## Install a cluster extending nodes to Local Zone <a name="ipi-localzones"></a>

Starting on 4.14 you can install an OCP cluster in AWS extending nodes to the AWS Local Zones,
letting the installation process automate all the steps from the subnet creation to
node running through MachineSet manifests.

There are some design considerations when using the fully automated process:

- Read the [AWS Local Zones limitations](ocp-aws-localzone-limitations)
- Cluster-wide network MTU: the Maximum Transmission Unit for the overlay network will automatically be adjusted when the edge pool configuration is set
- Machine Network CIDR block allocation: the Machine CIDR blocks used to create the cluster will be sharded to smaller blocks depending on the number of zones provided on install-config.yaml to create the public and private subnets.
- Internet egress traffic for private subnets: When using the installer automation to create subnets in Local Zones, the egress traffic for private subnets in AWS Local Zones will use the Nat Gateway from the parent zone, when the parent zone's route table is present, otherwise it will use the first route table for private subnets found in the region.

The sections below describe how to create a cluster using a basic example with single-zone local, and a full example of retrieving all zones in the region.

### Prerequisites

The prerequisite for installing a cluster using AWS Local Zones is to opt-in to every Local Zone group.

For Local Zones, the group name must be the zone name without the letter (zone identifier). Example: for Local Zone `us-east-1-bos-1a` the zone group will be `us-east-1-bos-1`.

It's also possible to query the group name reading the zone attribute:

```bash
$ aws --region us-east-1 ec2 describe-availability-zones \
  --all-availability-zones \
  --filters Name=zone-name,Values=us-east-1-bos-1a \
  --query "AvailabilityZones[].GroupName" --output text
us-east-1-bos-1
```

### Option 1. Steps to create a cluster with a single Local Zone

<!-- > Note: this example preferably goes to the product documentation. -->

Create a cluster in the region `us-east-1` extending worker nodes to AWS Local Zone `us-east-1-bos-1a`:

- opt-in to the Zone Group

```bash
aws ec2 modify-availability-zone-group \
    --region us-east-1 \
    --group-name us-east-1-bos-1a \
    --opt-in-status opted-in
```

AWS will process the request in the background, it could take a few minutes. Check the field `OptInStatus` before proceeding:

```bash
aws --region us-east-1 ec2 describe-availability-zones \
  --all-availability-zones \
  --filters Name=zone-name,Values=us-east-1-bos-1a \
  --query "AvailabilityZones[]"
```

- Create the `install-config.yaml`:

```yaml
apiVersion: v1
publish: External
baseDomain: devcluster.openshift.com
metadata:
  name: "cluster-name"
pullSecret: ...
sshKey: ...
platform:
  aws:
    region: us-east-1
compute:
- name: edge
  platform:
    aws:
      zones:
      - us-east-1-bos-1a
```

- Create the cluster

```bash
./openshift-install create cluster
```

### Option 2. Steps to create a cluster with many zones

Steps to create a cluster using the AWS Region `us-east-1` as a reference, selecting all Local Zones in the Region.

- Build the lists for zone groups and names:

```bash
mapfile -t local_zone_names < <(aws --region us-east-1 ec2 describe-availability-zones   --all-availability-zones   --filters Name=zone-type,Values=local-zone   --query "AvailabilityZones[].ZoneName" | jq -r .[])
mapfile -t local_zone_groups < <(aws --region us-east-1 ec2 describe-availability-zones   --all-availability-zones   --filters Name=zone-type,Values=local-zone   --query "AvailabilityZones[].GroupName" | jq -r .[])
```

- Opt-in the zone group:

```bash
for zone_group in ${local_zone_groups[@]}; do
  aws ec2 modify-availability-zone-group \
    --region us-east-1 \
    --group-name ${zone_group} \
    --opt-in-status opted-in
done
```

- Export the zone list:

```bash
$ for zone_name in ${local_zone_names[@]}; do echo "      - $zone_name"; done
      - us-east-1-atl-1a
      - us-east-1-bos-1a
      - us-east-1-bue-1a
      - us-east-1-chi-1a
      - us-east-1-dfw-1a
      - us-east-1-iah-1a
      - us-east-1-lim-1a
      - us-east-1-mci-1a
      - us-east-1-mia-1a
      - us-east-1-msp-1a
      - us-east-1-nyc-1a
      - us-east-1-phl-1a
      - us-east-1-qro-1a
      - us-east-1-scl-1a
```

- Create the `install-config.yaml` with the local zone list:

```yaml
apiVersion: v1
publish: External
baseDomain: devcluster.openshift.com
metadata:
  name: "cluster-name"
pullSecret: ...
sshKey: ...
platform:
  aws:
    region: us-east-1
compute:
- name: edge
  platform:
    aws:
      zones:
      - us-east-1-atl-1a
      - us-east-1-bos-1a
      - us-east-1-bue-1a
      - us-east-1-chi-1a
      - us-east-1-dfw-1a
      - us-east-1-iah-1a
      - us-east-1-lim-1a
      - us-east-1-mci-1a
      - us-east-1-mia-1a
      - us-east-1-msp-1a
      - us-east-1-nyc-1a
      - us-east-1-phl-1a
      - us-east-1-qro-1a
      - us-east-1-scl-1a
```

For each specified zone, a CIDR block range will be allocated, and subnets created.

- Create the cluster

```bash
./openshift-install create cluster
```

## Install a cluster into existing VPC with Local Zone subnets <a name="ipi-localzones-existing-vpc"></a>

The steps below describe how to install a cluster in existing VPC with AWS Local Zones subnets using Edge Machine Pool, introduced in 4.12.

The Edge Machine Pool was created to create a pool of workers running in the AWS Local Zones locations. This pool differs from the default compute pool on these items - Edge workers was not designed to run regular cluster workloads:
- The resources in AWS Local Zones are more expensive than the normal availability zones
- The latency between the application and end-users is lower in Local Zones and may vary for each location. So it will impact if some workloads like routers are mixed in the normal availability zones due to the unbalanced latency
- Network Load Balancers do not support subnets in the Local Zones
- The total time to connect to the applications running in Local Zones from the end-users close to the metropolitan region running the workload, is almost 10x faster than the parent region.

Table of Contents:

- [Prerequisites](#prerequisites)
    - [Additional IAM permissions](#prerequisites-iam)
- [Create the Network stack](#create-network)
    - [Create the VPC](#create-network-vpc)
    - [Create the Local Zone subnet](#create-network-subnet)
        - [Opt-in zone group](#create-network-subnet-optin)
        - [Creating the Subnet using AWS CloudFormation](#create-network-subnet-cfn)
- [Install](#install-cluster)
    - [Create the install-config.yaml](#create-config)
    - [Setting up the Edge Machine Pool](#create-config-edge-pool")
        - [Example edge pool created without customization](#create-config-edge-pool)
        - [Example edge pool with custom Instance type](#reate-config-edge-pool-example-ec2)
        - [Example edge pool with custom EBS type](#create-config-edge-pool-example-ebs)
    - [Create the cluster](#create-cluster-run)
- [Uninstall](#uninstall)
    - [Destroy the cluster](#uninstall-destroy-cluster)
    - [Destroy the Local Zone subnet](#uninstall-destroy-subnet)
    - [Destroy the VPC](#uninstall-destroy-vpc)
- [Use Cases](#use-cases)
    - [Example of a sample application deployment](#uc-deployment)
    - [User-workload ingress traffic](#uc-exposing-ingress)

To install a cluster in an existing VPC with Local Zone subnets, you should provision the network resources and then add the subnet IDs to the `install-config.yaml`.

## Prerequisites <a name="prerequisites"></a>

- [AWS Command Line Interface](aws-cli)
- [openshift-install >= 4.12](openshift-install)
- environment variables exported:

```bash
export CLUSTER_NAME="ipi-localzones"

# AWS Region and extra Local Zone group Information
export AWS_REGION="us-west-2"
export ZONE_GROUP_NAME="us-west-2-lax-1"
export ZONE_NAME="us-west-2-lax-1a"

# VPC Information
export VPC_CIDR="10.0.0.0/16"
export VPC_SUBNETS_BITS="10"
export VPC_SUBNETS_COUNT="3"

# Local Zone Subnet information
export SUBNET_CIDR="10.0.192.0/22"
export SUBNET_NAME="${CLUSTER_NAME}-public-usw2-lax-1a"
```

### Additional IAM permissions

The AWS Local Zone deployment described in this document, requires the additional permission for the user creating the cluster to modify the Local Zone group: `ec2:ModifyAvailabilityZoneGroup`

Example of the permissive IAM Policy that can be attached to the User or Role:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:ModifyAvailabilityZoneGroup"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```

## Create the Network Stack <a name="create-network"></a>

### Create the VPC <a name="create-network-vpc"></a>

The steps to install a cluster in an existing VPC are [detailed in the official documentation](aws-install-vpc). You can alternatively use [the CloudFormation templates to create the Network resources](aws-install-cloudformation), which will be used in this document.

- Create the Stack

```bash
INSTALLER_URL="https://raw.githubusercontent.com/openshift/installer/master"
TPL_URL="${INSTALLER_URL}/upi/aws/cloudformation/01_vpc.yaml"

aws cloudformation create-stack \
    --region ${AWS_REGION} \
    --stack-name ${CLUSTER_NAME}-vpc \
    --template-body ${TPL_URL} \
    --parameters \
        ParameterKey=VpcCidr,ParameterValue=${VPC_CIDR} \
        ParameterKey=SubnetBits,ParameterValue=${VPC_SUBNETS_BITS} \
        ParameterKey=AvailabilityZoneCount,ParameterValue=${VPC_SUBNETS_COUNT}
```

- Wait for the stack to be created: `StackStatus=CREATE_COMPLETE`

```bash
aws cloudformation wait stack-create-complete \
    --region ${AWS_REGION} \
    --stack-name ${CLUSTER_NAME}-vpc 
```

- Export the VPC ID:

```bash
export VPC_ID=$(aws cloudformation describe-stacks \
  --region ${AWS_REGION} --stack-name ${CLUSTER_NAME}-vpc \
  --query 'Stacks[0].Outputs[?OutputKey==`VpcId`].OutputValue' --output text)
```

- Extract the subnets IDs to the environment variable list `SUBNETS`:

```bash
mapfile -t SUBNETS < <(aws cloudformation describe-stacks \
  --region ${AWS_REGION} \
  --stack-name ${CLUSTER_NAME}-vpc \
  --query 'Stacks[0].Outputs[?OutputKey==`PrivateSubnetIds`].OutputValue' \
  --output text | tr ',' '\n')
mapfile -t -O "${#SUBNETS[@]}" SUBNETS < <(aws cloudformation describe-stacks \
  --region ${AWS_REGION} \
  --stack-name ${CLUSTER_NAME}-vpc  \
  --query 'Stacks[0].Outputs[?OutputKey==`PublicSubnetIds`].OutputValue' \
  --output text | tr ',' '\n')
```

- Export the Public Route Table ID:

```bash
export PUBLIC_RTB_ID=$(aws cloudformation describe-stacks \
  --region us-west-2 \
  --stack-name ${CLUSTER_NAME}-vpc \
  --query 'Stacks[0].Outputs[?OutputKey==`PublicRouteTableId`].OutputValue' --output text)
```

- Make sure all variables have been correctly set:

```bash
echo "SUBNETS=${SUBNETS[*]}
VPC_ID=${VPC_ID}
PUBLIC_RTB_ID=${PUBLIC_RTB_ID}"
```

### Create the Local Zone subnet <a name="create-network-subnet"></a>

The following actions are required to create subnets in Local Zones:
- choose the zone group to be enabled
- opt-in the zone group

#### Opt-in Zone groups <a name="create-network-subnet-optin"></a>

Opt-in the zone group:

```bash
aws ec2 modify-availability-zone-group \
    --region ${AWS_REGION} \
    --group-name ${ZONE_GROUP_NAME} \
    --opt-in-status opted-in
```

#### Creating the Subnet using AWS CloudFormation <a name="create-network-subnet-cfn"></a>

- Create the Stack for Local Zone subnet `us-west-2-lax-1a`

```bash
INSTALLER_URL="https://raw.githubusercontent.com/openshift/installer/master"
TPL_URL="${INSTALLER_URL}/upi/aws/cloudformation/01.99_net_local-zone.yaml"

aws cloudformation create-stack \
    --region ${AWS_REGION} \
    --stack-name ${SUBNET_NAME} \
    --template-body ${TPL_URL} \
    --parameters \
        ParameterKey=VpcId,ParameterValue=${VPC_ID} \
        ParameterKey=ZoneName,ParameterValue=${ZONE_NAME} \
        ParameterKey=SubnetName,ParameterValue=${SUBNET_NAME} \
        ParameterKey=PublicSubnetCidr,ParameterValue=${SUBNET_CIDR} \
        ParameterKey=PublicRouteTableId,ParameterValue=${PUBLIC_RTB_ID}
```

- Wait for the stack to be created `StackStatus=CREATE_COMPLETE`

```bash
aws cloudformation wait stack-create-complete \
  --region ${AWS_REGION} \
  --stack-name ${SUBNET_NAME}
```

- Export the Local Zone subnet ID

```bash
export SUBNET_ID=$(aws cloudformation describe-stacks \
  --region ${AWS_REGION} \
  --stack-name ${SUBNET_NAME} \
  --query 'Stacks[0].Outputs[?OutputKey==`PublicSubnetIds`].OutputValue' --output text)

# Append the Local Zone Subnet ID to the Subnet List
SUBNETS+=(${SUBNET_ID})
```

- Check the total of subnets. If you choose 3 AZs to be created on the VPC stack, you should have 7 subnets on this list:

```bash
$ echo ${#SUBNETS[*]}
7
```

## Install the cluster <a name="install-cluster"></a>

To install the cluster in existing VPC with subnets in Local Zones, you should:
- generate the `install-config.yaml`, or provide yours
- add the subnet IDs by setting the option `platform.aws.subnets`
- (optional) customize the `edge` compute pool

### Create the install-config.yaml <a name="create-config"></a>

Create the `install-config.yaml` providing the subnet IDs recently created:

- create the `install-config`

```bash
$ ./openshift-install create install-config --dir ${CLUSTER_NAME}
? SSH Public Key /home/user/.ssh/id_rsa.pub
? Platform aws
? Region us-west-2
? Base Domain devcluster.openshift.com
? Cluster Name ipi-localzone
? Pull Secret [? for help] **
INFO Install-Config created in: ipi-localzone     
```

- Append the subnets to the `platform.aws.subnets`:

```bash
$ echo "    subnets:"; for SB in ${SUBNETS[*]}; do echo "    - $SB"; done
    subnets:
    - subnet-0fc845d8e30fdb431
    - subnet-0a2675b7cbac2e537
    - subnet-01c0ac400e1920b47
    - subnet-0fee60966b7a93da6
    - subnet-002b48c0a91c8c641
    - subnet-093f00deb44ce81f4
    - subnet-0f85ae65796e8d107
```

### Setting up the Edge Machine Pool <a name="create-config-edge-pool"></a>

Version 4.12 or later introduces a new compute pool named `edge` designed for
the remote zones. The `edge` compute pool configuration is common between
AWS Local Zone locations, but due to the limitation of resources (Instance Types
and Sizes) of the Local Zone, the default instance type created may vary
from the traditional worker pool.

The default EBS for Local Zone locations is `gp2`, different than the default worker pool.

The preferred list of instance types follows the same order of worker pools, depending
on the availability of the location, one of those instances will be chosen*:
> Note: This list can be updated over time
- `m6i.xlarge`
- `m5.xlarge`
- `c5d.2xlarge`

The `edge` compute pool will also create new labels to help developers
deploy their applications onto those locations. The new labels introduced are:
    - `node-role.kubernetes.io/edge=''`
    - `zone_type=local-zone`
    - `zone_group=<Local Zone Group>`

Finally, the Machine Sets created by the `edge` compute pool have `NoSchedule` taint to avoid the
regular workloads spread out on those machines, and only user workloads will be allowed to run
when the tolerations are defined on the pod spec (you can see the example in the following sections).

By default, the `edge` compute pool will be created only when AWS Local Zone subnet IDs are added
to the list of `platform.aws.subnets`.

See below some examples of `install-config.yaml` with `edge` compute pool.

#### Example edge pool created without customization <a name="create-config-edge-pool-example-def"></a>

```yaml
apiVersion: v1
baseDomain: devcluster.openshift.com
metadata:
  name: ipi-localzone
platform:
  aws:
    region: us-west-2
    subnets:
    - subnet-0fc845d8e30fdb431
    - subnet-0a2675b7cbac2e537
    - subnet-01c0ac400e1920b47
    - subnet-0fee60966b7a93da6
    - subnet-002b48c0a91c8c641
    - subnet-093f00deb44ce81f4
    - subnet-0f85ae65796e8d107
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

#### Example edge pool with custom Instance type <a name="create-config-edge-pool-example-ec2"></a>

The Instance Type may differ between locations. You should check the AWS Documentation to check availability in the Local Zone that the cluster will run.

`install-config.yaml` example customizing the Instance Type for the Edge Machine Pool:

```yaml
apiVersion: v1
baseDomain: devcluster.openshift.com
metadata:
  name: ipi-localzone
compute:
- name: edge
  platform:
    aws:
      type: m5.4xlarge
platform:
  aws:
    region: us-west-2
    subnets:
    - subnet-0fc845d8e30fdb431
    - subnet-0a2675b7cbac2e537
    - subnet-01c0ac400e1920b47
    - subnet-0fee60966b7a93da6
    - subnet-002b48c0a91c8c641
    - subnet-093f00deb44ce81f4
    - subnet-0f85ae65796e8d107
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

#### Example edge pool with custom EBS type <a name="create-config-edge-pool-example-ebs"></a>

The EBS Type may differ between locations. You should check the AWS Documentation to check availability in the Local Zone that the cluster will run.

`install-config.yaml` example customizing the EBS Type for the Edge Machine Pool:

```yaml
apiVersion: v1
baseDomain: devcluster.openshift.com
metadata:
  name: ipi-localzone
compute:
- name: edge
  platform:
    aws:
      rootVolume:
        type: gp3
        size: 120
platform:
  aws:
    region: us-west-2
    subnets:
    - subnet-0fc845d8e30fdb431
    - subnet-0a2675b7cbac2e537
    - subnet-01c0ac400e1920b47
    - subnet-0fee60966b7a93da6
    - subnet-002b48c0a91c8c641
    - subnet-093f00deb44ce81f4
    - subnet-0f85ae65796e8d107
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Create the cluster <a name="create-cluster-run"></a>

```bash
./openshift-install create cluster --dir ${CLUSTER_NAME}
```

### Uninstall the cluster <a name="uninstall"></a>

#### Destroy the cluster <a name="uninstall-destroy-cluster"></a>

```bash
./openshift-install destroy cluster --dir ${CLUSTER_NAME}
```

#### Destroy the Local Zone subnets <a name="uninstall-destroy-subnet"></a>

```bash
aws cloudformation delete-stack \
    --region ${AWS_REGION} \
    --stack-name ${SUBNET_NAME}
```

#### Destroy the VPC <a name="uninstall-destroy-vpc"></a>

```bash
aws cloudformation delete-stack \
    --region ${AWS_REGION} \
    --stack-name ${CLUSTER_NAME}-vpc
```

## Extend worker nodes to AWS Local Zones in existing clusters [Day 2] <a name="#day2-localzones"></a>

The following steps are required to create worker nodes in AWS Local Zones:

- Make sure the overlay network MTU is set correctly to support the AWS Local Zone limitations
- Create subnets in AWS Local Zones, and dependencies (subnet association)
- Create MachineSet to deploy compute nodes in Local Zone subnets

When the cluster is installed using the edge compute pool, the MTU for the overlay network is automatically adjusted depending on the network plugin used.

When the cluster was already installed without the edge compute pool, and without Local Zone support, the required dependencies must be satisfied. The steps below cover both scenarios.

### Adjust the MTU of the overlay network

> You can skip this section if the cluster is already installed with Local Zone support.

The [KCS](https://access.redhat.com/solutions/6996487) covers the required step to change the MTU from the overlay network.

***Example changing the default MTU (9001) to the maximum allowed for network plugin OVN-Kubernetes***:

```bash

$ CLUSTER_MTU_CUR=$(oc get network.config.openshift.io/cluster --output=jsonpath={.status.clusterNetworkMTU})
$ CLUSTER_MTU_NEW=1200

$ oc patch Network.operator.openshift.io cluster --type=merge \
  --patch "{
    \"spec\":{
      \"migration\":{
        \"mtu\":{
          \"network\":{
            \"from\":${CLUSTER_MTU_CUR},
            \"to\":${CLUSTER_MTU_NEW}
          },
          \"machine\":{\"to\":9001}
        }}}}"
```

Wait for the deployment to be finished, then remove the migration config:

```bash
$ oc patch network.operator.openshift.io/cluster --type=merge \
  --patch "{
    \"spec\":{
      \"migration\":null,
      \"defaultNetwork\":{
        \"ovnKubernetesConfig\":{\"mtu\":${CLUSTER_MTU_NEW}}
        }}}"
```

### Setup subnet for Local Zone

Prerequisites:

- You must check the free CIDR Blocks available on the VPC
- Only CloudFormation Templates for public subnets are provided, you must adapt them if need more advanced configuration

Steps:

- [Opt-in the Zone group](https://docs.openshift.com/container-platform/4.12/installing/installing_aws/installing-aws-localzone.html#installation-aws-add-local-zone-locations_installing-aws-localzone)
- [Create the Local Zone subnet](https://docs.openshift.com/container-platform/4.12/installing/installing_aws/installing-aws-localzone.html#installation-creating-aws-vpc-localzone_installing-aws-localzone)


### Create the MachineSet

The steps below describe how to create the MachineSet manifests for the AWS Local Zone node:

- [Create the MachineSet manifest: Step 3](https://docs.openshift.com/container-platform/4.12/installing/installing_aws/installing-aws-localzone.html#installation-localzone-generate-k8s-manifest_installing-aws-localzone)

Once it is created you can apply the configuration to the cluster:

***Example:***

```bash
oc create -f <installation_directory>/openshift/99_openshift-cluster-api_worker-machineset-nyc1.yaml
```

___
___

# Wavelength Zones

## Prerequisites

### Review Wavelength Zone limitations

There are some design considerations when using the fully automated process in OpenShift:

- Review the AWS Wavelength Zones documentation for [Overview](https://docs.aws.amazon.com/wavelength/latest/developerguide/what-is-wavelength.html) and [Quotas and considerations](https://docs.aws.amazon.com/wavelength/latest/developerguide/wavelength-quotas.html)
- Cluster-wide network MTU: the Maximum Transmission Unit for the overlay network will automatically be adjusted when the edge pool configuration is set
- Machine Network CIDR block allocation: the Machine CIDR blocks used to create the cluster will be sharded to smaller blocks depending on the number of zones provided on install-config.yaml to create the public and private subnets.
- Internet egress traffic for private subnets: When using the installer automation to create subnets in Wavelength Zones, the egress traffic for private subnets in AWS Wavelength Zones will use the Nat Gateway from the parent zone, when the parent zone's route table is present, otherwise it will use the first route table for private subnets found in the region.

### Opt-into AWS Wavelength Zone

Opt into AWS Wavelength Zones.

Check the zone group name for the target zone (`us-east-1-wl1-bos-wlz-1`):

```sh
$ aws --region us-east-1 ec2 describe-availability-zones \
  --all-availability-zones \
  --filters Name=zone-name,Values=us-east-1-wl1-bos-wlz-1 \
  --query "AvailabilityZones[].GroupName" --output text
us-east-1-wl1
```

Opt-in to the Zone Group

```bash
aws ec2 modify-availability-zone-group \
    --region us-east-1 \
    --group-name us-east-1-wl1 \
    --opt-in-status opted-in
```

The request will be processed in the background, it could take a few minutes. Check if the field `OptInStatus` has the value `opted-in` before proceeding:

```bash
aws --region us-east-1 ec2 describe-availability-zones \
  --all-availability-zones \
  --filters Name=zone-name,Values=us-east-1-wl1-bos-wlz-1 \
  --query "AvailabilityZones[].OptInStatus"
```

## Install a cluster extending nodes to the Wavelength Zone [new VPC] <a name="#ipi-wavelength-auto"></a>

### Prerequisites

#### Additional AWS Permissions

IAM Permissions when the installer fully automates the creation and deletion of subnets in Wavelength zones.

- [Opt-int permissions](#pre-iam-opt-in)

- Permissions to create and delete the Carrier Gateway:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DeleteCarrierGateway",
        "ec2:CreateCarrierGateway"
      ],
      "Resource": "*"
    }
  ]
}
```

### Create cluster

Create a cluster in the region `us-east-1` extending worker nodes to AWS Local Zone `us-east-1-wl1-bos-wlz-1`:

- Create the `install-config.yaml`:

```sh
CLUSTER_NAME=aws-wlz
INSTALL_DIR=${PWD}/installdir-${CLUSTER_NAME}
mkdir $INSTALL_DIR
cat << EOF > $INSTALL_DIR/install-config.yaml
apiVersion: v1
metadata:
  name: $CLUSTER_NAME
publish: External
pullSecret: '$(cat ~/.openshift/pull-secret-latest.json)'
sshKey: |
  $(cat ~/.ssh/id_rsa.pub)
baseDomain: devcluster.openshift.com
platform:
  aws:
    region: us-east-1
compute:
- name: edge
  platform:
    aws:
      zones:
      - us-east-1-wl1-bos-wlz-1
EOF
```

- Create the cluster

```bash
./openshift-install create cluster --dir ${$INSTALL_DIR}
```

 Create the cluster

```bash
./openshift-install destroy cluster --dir ${$INSTALL_DIR}
```

## Install a cluster on AWS in existing VPC with subnets in Wavelength Zone <a name="#ipi-wavelength-byovpc"></a>

This section describes how to create the CloudFormation stack to provision VPC and subnets in Wavelength Zones, and then install an OpenShift cluster into an existing network.

### Prerequisites

- [Opt-into AWS Wavelength Zone](#opt-into-aws-wavelength-zone)

### Create the Network Stack (VPC and subnets)

Steps:

- Export the general variables for the cluster, and adapt them according to your environment:

```sh
export CLUSTER_REGION=us-east-1
export CLUSTER_NAME=wlz-byovpc
export PULL_SECRET_FILE=${HOME}/path/to/pull-secret.json
export BASE_DOMAIN=example.com
export SSH_PUB_KEY_FILE=$HOME/.ssh/id_rsa.pub

export CIDR_VPC="10.0.0.0/16"

# Set the Wavelength Zone to create subnets
export ZONE_NAME="us-east-1-wl1-nyc-wlz-1"
export SUBNET_CIDR_PUB="10.0.128.0/24"
export SUBNET_CIDR_PVT="10.0.129.0/24"
```

- Export the CloudFormation template path (assuming you are in the root of the installer repository):

```sh
TEMPLATE_NAME_VPC="upi/aws/cloudformation/01_vpc.yaml"
TEMPLATE_NAME_CARRIER_GW="upi/aws/cloudformation/01_vpc_01_carrier_gateway.yaml"
TEMPLATE_NAME_SUBNET="upi/aws/cloudformation/01_vpc_99_subnet.yaml"
```

- Create the CloudFormation stack for VPC:

```sh
export STACK_VPC=${CLUSTER_NAME}-vpc
aws cloudformation create-stack \
  --region ${CLUSTER_REGION} \
  --stack-name ${STACK_VPC} \
  --template-body file://$TEMPLATE_NAME_VPC \
  --parameters \
    ParameterKey=VpcCidr,ParameterValue="${CIDR_VPC}" \
    ParameterKey=AvailabilityZoneCount,ParameterValue=3 \
    ParameterKey=SubnetBits,ParameterValue=12

aws --region $CLUSTER_REGION cloudformation wait stack-create-complete --stack-name ${STACK_VPC}
aws --region $CLUSTER_REGION cloudformation describe-stacks --stack-name ${STACK_VPC}

export VPC_ID=$(aws --region $CLUSTER_REGION cloudformation describe-stacks \
  --stack-name ${STACK_VPC} \
  | jq -r '.Stacks[0].Outputs[] | select(.OutputKey=="VpcId").OutputValue' )
```

- Create the Carrier Gateway:

```sh
export STACK_CAGW=${CLUSTER_NAME}-cagw
aws cloudformation create-stack \
  --region ${CLUSTER_REGION} \
  --stack-name ${STACK_CAGW} \
  --template-body file://$TEMPLATE_NAME_CARRIER_GW \
  --parameters \
    ParameterKey=VpcId,ParameterValue="${VPC_ID}" \
    ParameterKey=ClusterName,ParameterValue="${CLUSTER_NAME}"

aws --region $CLUSTER_REGION cloudformation wait stack-create-complete --stack-name ${STACK_CAGW}
aws --region $CLUSTER_REGION cloudformation describe-stacks --stack-name ${STACK_CAGW}
```

- Extract the variables to create the subnets

```sh
export ZONE_SUFFIX=$(echo ${ZONE_NAME/${CLUSTER_REGION}-/})

export ROUTE_TABLE_PUB=$(aws --region $CLUSTER_REGION cloudformation describe-stacks \
  --stack-name ${STACK_CAGW} \
  | jq -r '.Stacks[0].Outputs[] | select(.OutputKey=="PublicRouteTableId").OutputValue' )

export ROUTE_TABLE_PVT=$(aws --region $CLUSTER_REGION cloudformation describe-stacks \
  --stack-name ${STACK_VPC} \
  | jq -r '.Stacks[0].Outputs[]
    | select(.OutputKey=="PrivateRouteTableIds").OutputValue
    | split(",")[0] | split("=")[1]' \
)

# Review the variables (optional)
cat <<EOF
CLUSTER_REGION=$CLUSTER_REGION
VPC_ID=$VPC_ID
AZ_NAME=$AZ_NAME
AZ_SUFFIX=$AZ_SUFFIX
ZONE_GROUP_NAME=$ZONE_GROUP_NAME
ROUTE_TABLE_PUB=$ROUTE_TABLE_PUB
ROUTE_TABLE_PVT=$ROUTE_TABLE_PVT
SUBNET_CIDR_PUB=$SUBNET_CIDR_PUB
SUBNET_CIDR_PVT=$SUBNET_CIDR_PVT
EOF
```

- Create the CloudFormation stack to provision the public and private subnets:

```sh
export STACK_SUBNET=${CLUSTER_NAME}-subnets-${AZ_SUFFIX}
aws cloudformation create-stack \
  --region ${CLUSTER_REGION} \
  --stack-name ${STACK_SUBNET} \
  --template-body file://$TEMPLATE_NAME_SUBNET \
  --parameters \
    ParameterKey=VpcId,ParameterValue="${VPC_ID}" \
    ParameterKey=ClusterName,ParameterValue="${CLUSTER_NAME}" \
    ParameterKey=ZoneName,ParameterValue="${AZ_NAME}" \
    ParameterKey=PublicRouteTableId,ParameterValue="${ROUTE_TABLE_PUB}" \
    ParameterKey=PublicSubnetCidr,ParameterValue="${SUBNET_CIDR_PUB}" \
    ParameterKey=PrivateRouteTableId,ParameterValue="${ROUTE_TABLE_PVT}" \
    ParameterKey=PrivateSubnetCidr,ParameterValue="${SUBNET_CIDR_PVT}"

aws --region $CLUSTER_REGION cloudformation wait stack-create-complete --stack-name ${STACK_SUBNET}
aws --region $CLUSTER_REGION cloudformation describe-stacks --stack-name ${STACK_SUBNET}
```

### Create the cluster

- Extract the subnets to be used in the install-config.yaml:

```sh
# Regular Availability Zones (public and private) from VPC CloudFormation Stack
mapfile -t SUBNETS < <(aws --region $CLUSTER_REGION cloudformation describe-stacks   --stack-name "${STACK_VPC}" --query "Stacks[0].Outputs[?OutputKey=='PrivateSubnetIds'].OutputValue" --output text | tr ',' '\n')

mapfile -t -O "${#SUBNETS[@]}" SUBNETS < <(aws --region $CLUSTER_REGION cloudformation describe-stacks   --stack-name "${STACK_VPC}" --query "Stacks[0].Outputs[?OutputKey=='PublicSubnetIds'].OutputValue" --output text | tr ',' '\n')

# Private subnet for Wavelenth Zones from subnets CloudFormation Stack
mapfile -t -O "${#SUBNETS[@]}" SUBNETS < <(aws --region $CLUSTER_REGION cloudformation describe-stacks   --stack-name "${STACK_SUBNET}" --query "Stacks[0].Outputs[?OutputKey=='PrivateSubnetIds'].OutputValue" --output text | tr ',' '\n')
```

- Create install-config.yaml:

```sh
cat <<EOF > ./install-config.yaml
apiVersion: v1
publish: External
baseDomain: ${BASE_DOMAIN}
metadata:
  name: "${CLUSTER_NAME}"
platform:
  aws:
    region: ${CLUSTER_REGION}
    subnets:
$(for SB in ${SUBNETS[*]}; do echo "    - $SB"; done)
pullSecret: '$(cat ${PULL_SECRET_FILE} | awk -v ORS= -v OFS= '{$1=$1}1')'
sshKey: |
  $(cat ${SSH_PUB_KEY_FILE})
EOF
```

- Create the cluster:

```sh
./openshift-install create cluster
```

### Destroy the cluster and network dependencies

- Destroy the cluster:

```sh
./openshift-install destroy cluster
```

- Destroy the subnet stack:

```sh
aws cloudformation delete-stack \
    --region ${AWS_REGION} \
    --stack-name ${STACK_SUBNET}
```

- Destroy the Carrier Gateway stack:

```sh
aws cloudformation delete-stack \
    --region ${AWS_REGION} \
    --stack-name ${STACK_CAGW}
```

- Destroy the VPC Stack:

```sh
aws cloudformation delete-stack \
    --region ${AWS_REGION} \
    --stack-name ${STACK_VPC}
```

___
___

# Use Cases <a name="use-cases"></a>

> Note: part of this document was added to the official documentation: [Post-installation configuration / Cluster tasks / Creating user workloads in AWS Local Zones](ocp-aws-localzones-day2-user-workloads)

## Example of a sample application deployment <a name="uc-deployment"></a>

The example below creates one sample application on the node running in the Local zone, setting the tolerations needed to pin the pod on the correct node:

```bash
cat << EOF | oc create -f -
apiVersion: v1
kind: Namespace
metadata:
  name: local-zone-demo
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: local-zone-demo-app-nyc-1
  namespace: local-zone-demo
spec:
  selector:
    matchLabels:
      app: local-zone-demo-app-nyc-1
  replicas: 1
  template:
    metadata:
      labels:
        app: local-zone-demo-app-nyc-1
        machine.openshift.io/zone-group: ${ZONE_GROUP_NAME}
    spec:
      nodeSelector:
        machine.openshift.io/zone-group: ${ZONE_GROUP_NAME}
      tolerations:
      - key: "node-role.kubernetes.io/edge"
        operator: "Equal"
        value: ""
        effect: "NoSchedule"
      containers:
        - image: openshift/origin-node
          command:
           - "/bin/socat"
          args:
            - TCP4-LISTEN:8080,reuseaddr,fork
            - EXEC:'/bin/bash -c \"printf \\\"HTTP/1.0 200 OK\r\n\r\n\\\"; sed -e \\\"/^\r/q\\\"\"'
          imagePullPolicy: Always
          name: echoserver
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service 
metadata:
  name:  local-zone-demo-app-nyc-1 
  namespace: local-zone-demo
spec:
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
  type: NodePort
  selector:
    app: local-zone-demo-app-nyc-1
EOF
```

## User-workload ingress traffic <a name="uc-exposing-ingress"></a>

To expose the applications to the internet on AWS Local Zones, application developers
must expose the applications using an external Load Balancer, for example, AWS Application Load Balancers (ALB). The
[ALB Operator](https://docs.openshift.com/container-platform/4.11/networking/aws_load_balancer_operator/install-aws-load-balancer-operator.html) is available through OLM on 4.11+.

To explore the best of deploying applications on the AWS Local Zone locations, at least one new
ALB `Ingress` must be provisioned by location to expose the services deployed on the
zones.

If the cluster-admin decides to share the ALB `Ingress` subnets between different locations,
it will impact drastically the latency for the end-users when the traffic is routed to
backends (compute nodes) placed in different zones that the traffic entered by the Ingress/Load Balancer.

The ALB deployment is not covered by this documentation.

___

[openshift-install]: https://docs.openshift.com/container-platform/4.11/installing/index.html
[aws-cli]: https://aws.amazon.com/cli/
[aws-install-vpc]: https://docs.openshift.com/container-platform/4.11/installing/installing_aws/installing-aws-vpc.html
[aws-install-cloudformation]: https://docs.openshift.com/container-platform/4.11/installing/installing_aws/installing-aws-user-infra.html
[aws-local-zones]: https://aws.amazon.com/about-aws/global-infrastructure/localzones
[aws-local-zones-features]: https://aws.amazon.com/about-aws/global-infrastructure/localzones/features
[ocp-aws-localzone-limitations]: https://docs.openshift.com/container-platform/4.13/installing/installing_aws/installing-aws-localzone.html#cluster-limitations-local-zone_installing-aws-localzone
[ocp-aws-localzones-day2-user-workloads]: https://docs.openshift.com/container-platform/4.13/post_installation_configuration/cluster-tasks.html#installation-extend-edge-nodes-aws-local-zones_post-install-cluster-tasks