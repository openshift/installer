# AWS Platform Customization

The following options are available when using AWS:

- `machines.platform.aws.rootVolume.iops` - the reserved IOPS of the root volume
- `machines.platform.aws.rootVolume.size` - the size (in GiB) of the root volume
- `machines.platform.aws.rootVolume.type` - the storage type of the root volume
- `machines.platform.aws.type` - the EC2 instance type
- `machines.platform.aws.zones` - a list of the availability zones that the installer will use when creating machines of this pool
- `platform.aws.region` - the AWS region that the installer will use when creating resources
- `platform.aws.userTags` - a map of keys and values that the installer will add as tags to all resources it creates

## Examples

An example `install-config.yaml` is shown below. This configuration has been modified to show the customization that is possible via the install config.

```yaml
apiVersion: v1beta3
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
networking:
  clusterNetworks:
  - cidr: 10.128.0.0/14
    hostSubnetLength: 9
  machineCIDR: 10.0.0.0/16
  serviceCIDR: 172.30.0.0/16
  type: OpenshiftSDN
platform:
  aws:
    region: us-west-2
    userTags:
      adminContact: jdoe
      costCenter: 7536
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```
