# AWS Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, AWS-specific properties.

## Cluster-scoped properties

* `amiID` (optional string): The AMI that should be used to boot machines for the cluster.
    If set, the AMI should belong to the same region as the cluster.
* `region` (required string): The AWS region where the cluster will be created.
* `subnets` (optional array of strings): Existing subnets (by ID) where cluster resources will be created.
    Leave unset to have the installer create subnets in a new VPC on your behalf.
* `userTags` (optional object): Additional keys and values that the installer will add as tags to all resources that it creates.
    Resources created by the cluster itself may not include these tags.
* `defaultMachinePlatform` (optional object): Default [AWS-specific machine pool properties](#machine-pools) which applies to [machine pools](../customization.md#machine-pools) that do not define their own AWS-specific properties.

## Machine pools

* `rootVolume` (optional object): Defines the root volume for EC2 instances in the machine pool.
    * `iops` (optional integer): The amount of provisioned [IOPS][volume-iops].
        This is only valid for `type` `io1`.
    * `size` (optional integer): Size of the root volume in gibibytes (GiB).
    * `type` (optional string):  The [type of volume][volume-type].
* `type` (optional string): The [EC2 instance type][instance-type].
* `zones` (optional array of strings): The availability zones used for machines in the pool.

## Installing to Existing VPC & Subnetworks

The installer can use an existing VPC and subnets when provisioning an OpenShift cluster. A VPC will be inferred from the provided subnets. For a standard installation, a private and public subnet should be specified. ([see example below](#pre-existing-vpc--subnets)). Both of the subnets must be within the IP range specified in `networking.machineNetwork`. 

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
      rootVolume:
        iops: 4000
        size: 500
        type: io1
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
    subnets:
    - subnet-0e953079d31ec4c74
    - subnet-05e6864f66a954c27
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

[availablity-zones]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html
[instance-type]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html
[volume-iops]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-io-characteristics.html
[volume-type]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/EBSVolumeTypes.html
