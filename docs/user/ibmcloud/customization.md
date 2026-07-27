# IBM Cloud Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, IBM Cloud-specific properties.

* `region` (required string): The IBM Cloud region where the cluster should be created.
* `resourceGroupName` (optional string): The name of an existing resource group where the cluster resources should be created. If not specified, a resource group named after the cluster's infrastructure ID is created.
* `networkResourceGroupName` (optional string): The name of the existing resource group containing the pre-existing VPC and subnets. Required when using an existing VPC.
* `vpcName` (optional string): The name of an existing VPC where the cluster infrastructure should be provisioned.
* `controlPlaneSubnets` (optional array of strings): The names of existing subnets which should be used by the cluster control plane.
* `computeSubnets` (optional array of strings): The names of existing subnets which should be used by the cluster compute nodes.
* `defaultMachinePlatform` (optional object): Default [IBM Cloud-specific machine pool properties](#machine-pools) applied to every [machine pool](../customization.md#machine-pools). Properties set on a machine pool, including individual `bootVolume` fields, take precedence over the default; unset properties are inherited from it. Inherited values are validated together with the pool's own settings, so for example a default `iops` combined with a pool `profile` of `general-purpose` is rejected.
* `serviceEndpoints` (optional array of objects): A list of custom service endpoint overrides. Each entry has a `name` and a `url`.

## Machine pools

* `type` (optional string): The [IBM Cloud VPC instance profile][instance-profiles] (for example `bx2-4x16`).
* `zones` (optional array of strings): The availability zones used for machines in the pool.
* `bootVolume` (optional object): Configuration for the machine's boot volume. The sizes, IOPS and bandwidth that can be configured depend on the selected [block storage profile][storage-profiles] and on the account. The installer checks each value on its own against the range the profile advertises in the region. Profiles can have additional restrictions, such as IOPS limits that depend on the volume size, so IBM Cloud may still reject a combination that passes validation; see the documentation for the full limits.
    * `encryptionKey` (optional string): The CRN of a Key Protect or Hyper Protect Crypto Services root key used to encrypt the boot volume. If not specified, a provider managed encryption key is used.
    * `profile` (optional string): The [block storage profile][storage-profiles] for the boot volume. Allowed values are `general-purpose`, `5iops-tier`, `10iops-tier`, `custom` and `sdp`. Defaults to `general-purpose`. Not every profile is available to every account and region; see the documentation linked above for current availability.
    * `sizeGiB` (optional integer): The size of the boot volume in GiB. Must be at least as large as the image's minimum provisioned size, and within the capacity range the [block storage profile][storage-profiles] supports for the account. If not specified, the IBM Cloud default boot volume size is used.
    * `iops` (optional integer): The maximum I/O operations per second for the boot volume. Only configurable for the `custom` and `sdp` profiles; for the other profiles the IOPS is fixed by the profile itself.
    * `bandwidth` (optional integer): The maximum bandwidth in megabits per second for the boot volume. Only configurable for the `sdp` profile. If not specified, the bandwidth is set by IBM Cloud.
* `dedicatedHosts` (optional array of objects): The dedicated hosts on which to provision the machines. The number of entries must match the number of `zones`, and `type` is required when dedicated hosts are used.
    * `name` (optional string): The name of an existing dedicated host to provision the machine on.
    * `profile` (optional string): The profile of a new dedicated host to create for the machine. Must be in the same profile family as `type`.

## Examples

### Boot volumes using the `sdp` profile

An example `install-config.yaml` applying an `sdp` boot volume to both the control plane and compute machine pools:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
controlPlane:
  architecture: amd64
  name: master
  platform:
    ibmcloud:
      type: bx2-8x32
      bootVolume:
        profile: sdp
        sizeGiB: 200
        iops: 10000
        bandwidth: 2000
  replicas: 3
compute:
- architecture: amd64
  name: worker
  platform:
    ibmcloud:
      type: bx2-4x16
      bootVolume:
        profile: sdp
        sizeGiB: 500
        iops: 16000
  replicas: 3
platform:
  ibmcloud:
    region: us-east
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### A default boot volume for all machine pools

`defaultMachinePlatform` applies to every machine pool. Properties set on a pool override the default field by field, so a pool that sets only `bootVolume.sizeGiB` still inherits `profile: sdp` from it:

```yaml
platform:
  ibmcloud:
    region: us-east
    defaultMachinePlatform:
      bootVolume:
        profile: sdp
        sizeGiB: 250
```

### An encrypted boot volume

```yaml
platform:
  ibmcloud:
    region: us-east
    defaultMachinePlatform:
      bootVolume:
        encryptionKey: crn:v1:bluemix:public:kms:us-east:a/accountid:instanceid:key:keyid
```

[instance-profiles]: https://cloud.ibm.com/docs/vpc?topic=vpc-profiles
[storage-profiles]: https://cloud.ibm.com/docs/vpc?topic=vpc-block-storage-profiles
