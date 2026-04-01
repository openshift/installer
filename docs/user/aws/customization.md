# AWS Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, AWS-specific properties.

## Cluster-scoped properties

* `amiID` (optional string): The AMI that should be used to boot machines for the cluster.
    If set, the AMI should belong to the same region as the cluster. This field is now deprecated and `defaultMachinePlatform` should be used instead.
* `region` (required string): The AWS region where the cluster will be created.
* `vpc` (optional object): The configurations of an existing VPC for the cluster.
  * `subnets` (optional array of objects): Defines the subnets in an existing VPC and can optionally specify their intended roles.
    Leave unset to have the installer create subnets in a new VPC on your behalf.
    When using dual-stack networking (`ipFamily` set to `DualStackIPv4Primary` or `DualStackIPv6Primary`), all subnets must have both IPv4 and IPv6 CIDR blocks associated with them.
    * `id` (required string): The subnet ID (e.g., `subnet-0e953079d31ec4c74`).
    * `roles` (optional array of objects): The roles that the subnet will provide in the cluster.
      If no roles are specified on any subnet, subnet roles are decided automatically.
      * `type` (required string): The role type. Valid values: `ClusterNode`, `EdgeNode`, `BootstrapNode`, `IngressControllerLB`, `ControlPlaneExternalLB`, `ControlPlaneInternalLB`.
* `subnets` (optional array of strings): Existing subnets (by ID) where cluster resources will be created.
    **Deprecated:** use `vpc.subnets` instead.
* `userTags` (optional object): Additional keys and values that the installer will add as tags to all resources that it creates.
    Resources created by the cluster itself may not include these tags.
* `defaultMachinePlatform` (optional object): Default [AWS-specific machine pool properties](#machine-pools) which applies to [machine pools](../customization.md#machine-pools) that do not define their own AWS-specific properties.
* `ipFamily` (optional string): Specifies the IP address family for the cluster network.
    * `IPv4` (default): IPv4-only networking
    * `DualStackIPv4Primary`: Dual-stack networking with both IPv4 and IPv6 addresses, with IPv4 as the primary address family
    * `DualStackIPv6Primary`: Dual-stack networking with both IPv4 and IPv6 addresses, with IPv6 as the primary address family

## Machine pools

* `rootVolume` (optional object): Defines the root volume for EC2 instances in the machine pool.
    * `iops` (optional integer): The amount of provisioned [IOPS][volume-iops].
        This is only valid for `type` `io1`.
    * `throughput` (optional integer): The amount of throughput in MiB/s [Throughput Performance][volume-throughput].
        This is only valid for `type` `gp3`.
    * `size` (optional integer): Size of the root volume in gibibytes (GiB).
    * `type` (optional string):  The [type of volume][volume-type].
    * `kmsKeyARN` (optional string): The [ARN of KMS key][kms-key] that should be used to encrypt the EBS volume.
        When no key is specified by user, the account's [default KMS Key][kms-key-default] for the region will be used.
        Example ARN values are: `arn:aws:kms:us-east-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab` or `arn:aws:kms:us-east-1:111122223333:alias/my-key`
* `type` (optional string): The [EC2 instance type][instance-type].
* `zones` (optional array of strings): The availability zones used for machines in the pool.
* `amiID` (optional string): The AMI that should be used to boot machines.
    If set, the AMI should belong to the same region as the cluster.

## Installing to Existing VPC & Subnetworks

The installer can use an existing VPC and subnets when provisioning an OpenShift cluster. A VPC will be inferred from the provided subnets. For a standard installation, a private and public subnet should be specified. ([see example below](#pre-existing-vpc--subnets)). Both of the subnets must be within the IP range specified in `networking.machineNetwork`.

For dual-stack installations, all subnets must have both IPv4 and IPv6 CIDR blocks associated with them, and the `networking.machineNetwork` must include both IPv4 and IPv6 CIDRs that match the subnet ranges. 

## Examples

Some example `install-config.yaml` are shown below.
For examples of platform-agnostic configuration fragments, see [here](../customization.md#examples).

### Minimal

An example minimal AWS install config is:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
platform:
  aws:
    region: us-west-2
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Custom machine pools

An example AWS install config with custom machine pools:

```yaml
apiVersion: v1
baseDomain: example.com
controlPlane:
  name: master
  platform:
    aws:
      zones:
      - us-west-2a
      - us-west-2b
      type: m5.xlarge
  replicas: 3
compute:
- name: worker
  platform:
    aws:
      amiID: ami-123456
      rootVolume:
        iops: 4000
        size: 500
        type: io1
        kmsKeyARN: arn:aws:kms:us-east-1:111122223333:key/1234abcd-12ab-34cd-56ef-1234567890ab
      type: c5.9xlarge
      zones:
      - us-west-2c
  replicas: 5
metadata:
  name: test-cluster
platform:
  aws:
    region: us-west-2
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Pre-existing VPC & Subnets

An example install config for installing to an existing VPC and subnets is:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
networking:
  machineNetwork:
  - cidr: 10.190.0.0/16
platform:
  aws:
    region: us-west-2
    vpc:
      subnets:
      - id: subnet-0e953079d31ec4c74
      - id: subnet-05e6864f66a954c27
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Dual-stack networking (installer-provisioned VPC)

An example install config for dual-stack networking with IPv6 as the primary address family, where the installer creates the VPC and subnets with both IPv6 and IPv4 CIDR blocks:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
networking:
  machineNetwork:
  - cidr: 10.0.0.0/16
  clusterNetwork:
  - cidr: fd01::/48
    hostPrefix: 64
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  serviceNetwork:
  - fd02::/112
  - 172.30.0.0/16
platform:
  aws:
    region: us-west-2
    ipFamily: DualStackIPv6Primary
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

The installer will automatically:
* Create a VPC with IPv4 CIDR (10.0.0.0/16 by default) and an Amazon-provided IPv6 CIDR block
* Create subnets with both IPv4 and IPv6 CIDR blocks
* Configure dual-stack networking for the cluster

### Dual-stack networking (existing VPC)

An example install config for dual-stack networking using existing subnets that have both IPv4 and IPv6 CIDR blocks:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
networking:
  machineNetwork:
  - cidr: 10.0.0.0/16
  - cidr: 2600:1f13:fe4::/56
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  - cidr: fd01::/48
    hostPrefix: 64
  serviceNetwork:
  - 172.30.0.0/16
  - fd02::/112
platform:
  aws:
    region: us-west-2
    ipFamily: DualStackIPv4Primary
    vpc:
      subnets:
      - id: subnet-0e953079d31ec4c74  # must have both IPv4 and IPv6 CIDR blocks
      - id: subnet-05e6864f66a954c27  # must have both IPv4 and IPv6 CIDR blocks
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

[availablity-zones]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html
[instance-type]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html
[kms-key-default]: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_GetEbsDefaultKmsKeyId.html
[kms-key]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSEncryption.html
[volume-iops]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-io-characteristics.html
[volume-throughput]: https://docs.aws.amazon.com/ebs/latest/userguide/general-purpose.html#gp3-ebs-volume-type
[volume-type]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html
