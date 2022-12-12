# Cluster Installation in existing VPC with Local Zones subnet

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
___

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
      "Sid": "Stmt1677614927608",
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
  --region us-west-2 \
  --stack-name ${CLUSTER_NAME}-vpc \
  | jq -r '.Stacks[0].Outputs[] | select(.OutputKey=="VpcId").OutputValue' )
```

- Extract the subnets IDs to the environment variable list `SUBNETS`:

```bash
mapfile -t SUBNETS < <(aws cloudformation describe-stacks \
  --region us-west-2 \
  --stack-name ${CLUSTER_NAME}-vpc \
  | jq -r '.Stacks[0].Outputs[0].OutputValue' | tr ',' '\n')
mapfile -t -O "${#SUBNETS[@]}" SUBNETS < <(aws cloudformation describe-stacks \
  --region us-west-2 \
  --stack-name ${CLUSTER_NAME}-vpc  \
  | jq -r '.Stacks[0].Outputs[1].OutputValue' | tr ',' '\n')
```

- Export the Public Route Table ID:

```bash
export PUBLIC_RTB_ID=$(aws cloudformation describe-stacks \
  --region us-west-2 \
  --stack-name ${CLUSTER_NAME}-vpc \
  | jq -r '.Stacks[0].Outputs[] | select(.OutputKey=="PublicRouteTableId").OutputValue' )
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
  | jq -r '.Stacks[0].Outputs[] | select(.OutputKey=="PublicSubnetIds").OutputValue' )

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
> Note: This list can be updated over the time
- `m6i.xlarge`
- `m5.xlarge`
- `c5d.2xlarge`

The `edge` compute pool will also create new labels to help developers to
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

## Use Cases <a name="use-cases"></a>

### Example of a sample application deployment <a name="uc-deployment"></a>

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

### User-workload ingress traffic <a name="uc-exposing-ingress"></a>

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


[openshift-install]: https://docs.openshift.com/container-platform/4.11/installing/index.html
[aws-cli]: https://aws.amazon.com/cli/
[aws-install-vpc]: https://docs.openshift.com/container-platform/4.11/installing/installing_aws/installing-aws-vpc.html
[aws-install-cloudformation]: https://docs.openshift.com/container-platform/4.11/installing/installing_aws/installing-aws-user-infra.html
[aws-local-zones]: https://aws.amazon.com/about-aws/global-infrastructure/localzones
[aws-local-zones-features]: https://aws.amazon.com/about-aws/global-infrastructure/localzones/features
