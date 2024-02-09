# Install OpenShift on AWS using custom/owned Public IPv4 Pool

Steps to create a cluster on AWS using Public IPv4 address pool
that you bring to your AWS account with BYOIP.

## Prerequisites

- Public IPv4 Pool Provisioned in the Account
- Total of ( (Zones*3 ) + 1) of Public IPv4 available in the pool, where: Zones is the total numbber of AWS zones used to deploy the OpenShift cluster.
    - Example to query the IPv4 pools available in the account, which returns the  `TotalAvailableAddressCount`:
```
aws ec2 describe-public-ipv4-pools --region us-east-1
```

## Steps

- Create the install config setting the field `platform.aws.publicIpv4PoolId`, and create the cluster:

```yaml
apiVersion: v1
baseDomain: ${CLUSTER_BASE_DOMAIN}
metadata:
  name: ocp-byoip
platform:
  aws:
    region: ${REGION}
    publicIpv4Pool: ipv4pool-ec2-123456789abcde
publish: External
pullSecret: '...'
sshKey: |
  '...'
```

- Create the cluster

```sh
openshift-install create cluster
```
