# OpenStack Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, OpenStack-specific properties.

## Table of Contents

* [Cluster-scoped properties](#cluster-scoped-properties)
* [Machine pools](#machine-pools)
* [Examples](#examples)
  * [Minimal](#minimal)
  * [Custom-machine-pools](#custom-machine-pools)
* [Image Overrides](#image-overrides)
* [Additional Networks](#additional-networks)
* [Additional Security Groups](#additional-security-groups)
* [Further customization](#further-customization)

## Cluster-scoped properties

* `cloud` (required string): The name of the OpenStack cloud to use from `clouds.yaml`.
* `computeFlavor` (required string): The OpenStack flavor to use for compute and control-plane machines.
    This is currently required, but has lower precedence than [the `type` property](#machine-pools) on [the `compute` and `controlPlane` machine-pools](../customization.md#platform-customization).
* `externalDNS` (optional list of strings): The IP addresses of DNS servers to be used for the DNS resolution of all instances in the cluster
* `externalNetwork` (required string): The OpenStack external network name to be used for installation.
* `lbFloatingIP` (required string): Existing Floating IP to associate with the API load balancer.
* `octaviaSupport` (optional string): Whether OpenStack supports Octavia (`1` for true or `0` for false)
* `region` (deprecated string): The OpenStack region where the cluster will be created. Currently this value is not used by the installer.
* `trunkSupport` (optional string): Whether OpenStack ports can be trunked (`1` for true or `0` for false)
* `clusterOSImage` (optional string): Either a URL with `http(s)` or `file` scheme to override the default OS image for cluster nodes or an existing Glance image name.
* `apiVIP` (optional string): An IP addresss on the machineNetwork that will be assigned to the API VIP. Be aware that the `10` and `11` of the machineNetwork will be taken by neutron dhcp by default, and wont be available.
* `ingressVIP` (optional string): An IP address on the machineNetwork that will be assigned to the ingress VIP. Be aware that the `10` and `11` of the machineNetwork will be taken by neutron dhcp by default, and wont be available.
* `machinesSubnet` (optional string): the UUID of an openstack subnet to install the nodes of the cluster onto. The first CIDR in `networks.machineNetwork` must match the cidr of the `machinesSubnet`. In order to support more complex networking configurations, we expect the subnet passed to already be connected to an external network in some way. When this option is set, we will no longer attempt to create a router. Also note that setting `externalDNS` while setting `machinesSubnet` is invalid usage. If you want to add a DNS to your cluster while using a custom subnet, add it to the subnet in openstack [like this](https://docs.openstack.org/neutron/rocky/admin/config-dns-res.html). 

## Machine pools

* `additionalNetworkIDs` (optional list of strings): IDs of additional networks for machines.
* `additionalSecurityGroupIDs` (optional list of strings): IDs of additional security groups for machines.
* `type` (optional string): The OpenStack flavor name for machines in the pool.
* `rootVolume` (optional object): Defines the root volume for instances in the machine pool. The instances use ephemeral disks if not set.
  * `size` (required integer): Size of the root volume in GB.
  * `type` (required string): The volume pool to create the volume from.

**NOTE:** The bootstrap node follows the `type` and `rootVolume` parameters from the `controlPlane` machine pool.

## Examples

Some example `install-config.yaml` are shown below.
For examples of platform-agnostic configuration fragments, see [here](../customization.md#examples).

### Minimal

An example minimal OpenStack install config is:

```yaml
apiVersion: v1
baseDomain: example.com
metadata:
  name: test-cluster
platform:
  openstack:
    cloud: mycloud
    computeFlavor: m1.s2.xlarge
    externalNetwork: external
    externalDNS:
      - "8.8.8.8"
      - "192.168.1.12"
    lbFloatingIP: 128.0.0.1
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

### Custom machine pools

An example OpenStack install config with custom machine pools:

```yaml
apiVersion: v1
baseDomain: example.com
controlPlane:
  name: master
  replicas: 3
compute:
- name: worker
  platform:
    openstack:
      type: ml.large
      rootVolume:
        size: 30
        type: performance
  replicas: 3
metadata:
  name: test-cluster
platform:
  openstack:
    cloud: mycloud
    computeFlavor: m1.s2.xlarge
    externalNetwork: external
    lbFloatingIP: 128.0.0.1
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

## Image Overrides

Normally the installer downloads the RHCOS image from a predetermined location described in [data/data/rhcos.json](/data/data/rhcos.json)). But the download URL can be overridden, notably for disconnected installations.

To do so and upload binary data from a custom location the user may set `clusterOSImage` parameter in the install config that points to that location, and then start the installation. In all other respects the process will be consistent with the default.

**NOTE:** For this to work, the parameter value must be a valid http(s) URL.

**NOTE:** The optional `sha256` query parameter can be attached to the URL, which will force the installer to check the image file checksum before uploading it into Glance.

Example:

```yaml
platform:
  openstack:
      clusterOSImage: http://mirror.example.com/images/rhcos-43.81.201912131630.0-openstack.x86_64.qcow2.gz?sha256=ffebbd68e8a1f2a245ca19522c16c86f67f9ac8e4e0c1f0a812b068b16f7265d
```

If the user wants to upload the image from the local file system, he can set `clusterOSImage` as `file:///path/to/file`. In this case the installer will take this file and automatically create an image in Glance.

Example:

```yaml
platform:
  openstack:
      clusterOSImage: file:///home/user/rhcos.qcow2
```

If the user wants to reuse an existing Glance image without any uploading of binary data, then it is possible to set `clusterOSImage` install config parameter that specifies the Glance image name. In this case no new Glance images will be created, and the image will stay when the cluster is destroyed. In other words, if `clusterOSImage` is not an "http(s)" or "file" URL, then the installer will look into Glance for an image with that name.

Example:

```yaml
platform:
  openstack:
      clusterOSImage: my-rhcos
```

## Additional Networks

You can set additional networks for your machines by defining `additionalNetworkIDs` parameter in the machine configuration. The parameter is a list of strings with additional network IDs:

```yaml
additionalNetworkIDs:
- <network1_uuid>
- <network2_uuid>
```

You can attach this parameter for both `controlPlane` and `compute` machines:

Example:

```yaml
compute:
- name: worker
  platform:
    openstack:
      additionalNetworkIDs:
      - fa806b2f-ac49-4bce-b9db-124bc64209bf
controlPlane:
  name: master
  platform:
    openstack:
      additionalNetworkIDs:
      - fa806b2f-ac49-4bce-b9db-124bc64209bf
```

**NOTE:** Allowed address pairs won't be created for the additional networks.

## Additional Security Groups

You can set additional security groups for your machines by defining `additionalSecurityGroupIDs` parameter in the machine configuration. The parameter is a list of strings with additional security group IDs:

```yaml
additionalSecurityGroupIDs:
- <security_group1_id>
- <security_group2_id>
```

You can attach this parameter for both `controlPlane` and `compute` machines:

Example:

```yaml
compute:
- name: worker
  platform:
    openstack:
      additionalSecurityGroupIDs:
      - 7ee219f3-d2e9-48a1-96c2-e7429f1b0da7
controlPlane:
  name: master
  platform:
    openstack:
      additionalSecurityGroupIDs:
      - 7ee219f3-d2e9-48a1-96c2-e7429f1b0da7
```

## Further customization

For customizing the installation beyond what is possible with `openshift-install`, refer to the [UPI (User Provided Infrastructure) documentation](./install_upi.md).
