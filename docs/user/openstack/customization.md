# OpenStack Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, OpenStack-specific properties.

## Table of Contents

- [OpenStack Platform Customization](#openstack-platform-customization)
  - [Table of Contents](#table-of-contents)
  - [Cluster-scoped properties](#cluster-scoped-properties)
  - [Machine pools](#machine-pools)
  - [Examples](#examples)
    - [Minimal](#minimal)
    - [Custom machine pools](#custom-machine-pools)
  - [Image Overrides](#image-overrides)
  - [Custom Subnets](#custom-subnets)
  - [Additional Networks](#additional-networks)
  - [Additional Security Groups](#additional-security-groups)
  - [Further customization](#further-customization)

## Cluster-scoped properties

* `cloud` (required string): The name of the OpenStack cloud to use from `clouds.yaml`.
* `computeFlavor` (deprecated string): The OpenStack flavor to use for compute and control-plane machines.
* `externalDNS` (optional list of strings): The IP addresses of DNS servers to be used for the DNS resolution of all instances in the cluster. The total number of dns servers supported by an instance is three. That total includes any dns server provided by the underlying openstack infrastructure.
* `externalNetwork` (optional string): Name of external network the installer will use to provide access to the cluster. If defined, a floating IP from this network will be created and associated with the bootstrap node to facilitate debugging and connection to the bootstrap node during installation. The `apiFloatingIP` property is a floating IP address selected from this network.
* `apiFloatingIP` (optional string): Address of existing Floating IP from externalNetwork the installer will associate with the OpenShift API. This property is only valid if externalNetwork is defined. If externalNetwork is not defined, the installer will throw an error.
* `ingressFloatingIP` (optional string): Address of an existing Floating IP from externalNetwork the installer will associate with the ingress port. This property is only valid if externalNetwork is defined. If externalNetwork is not defined, the installer will throw an error.
* `octaviaSupport` (deprecated string): Whether OpenStack supports Octavia (`1` for true or `0` for false)
* `region` (deprecated string): The OpenStack region where the cluster will be created. Currently this value is not used by the installer.
* `trunkSupport` (deprecated string): Whether OpenStack ports can be trunked (`1` for true or `0` for false)
* `clusterOSImage` (optional string): Either a URL with `http(s)` or `file` scheme to override the default OS image for cluster nodes or an existing Glance image name.
* `clusterOSImageProperties` (optional list of strings): a list of properties to be added to the installer-uploaded ClusterOSImage in Glance. The default is to not set any properties. `clusterOSImageProperties` is ignored when `clusterOSImage` points to an existing image in Glance.
* `apiVIP` (optional string): An IP address on the machineNetwork that will be assigned to the API VIP. Be aware that the `10` and `11` of the machineNetwork will be taken by neutron dhcp by default, and wont be available.
* `ingressVIP` (optional string): An IP address on the machineNetwork that will be assigned to the ingress VIP. Be aware that the `10` and `11` of the machineNetwork will be taken by neutron dhcp by default, and wont be available.
* `machinesSubnet` (optional string): the UUID of an OpenStack subnet to install the nodes of the cluster onto. For more information on how to install with a custom subnet, see the [custom subnets](#custom-subnets) section of the docs.
* `defaultMachinePlatform` (optional object): Default [OpenStack-specific machine pool properties](#machine-pools) which apply to [machine pools](../customization.md#machine-pools) that do not define their own OpenStack-specific properties.

## Machine pools

* `additionalNetworkIDs` (optional list of strings): IDs of additional networks for machines.
* `additionalSecurityGroupIDs` (optional list of strings): IDs of additional security groups for machines.
* `type` (optional string): The OpenStack flavor name for machines in the pool.
* `rootVolume` (optional object): Defines the root volume for instances in the machine pool. The instances use ephemeral disks if not set.
  * `size` (required integer): Size of the root volume in GB. Must be set to at least 25.
  * `type` (required string): The volume pool to create the volume from.
  * `zones` (optional list of strings): The names of the availability zones you want to install your root volumes on. If unset, the installer will use your default volume zone.
* `zones` (optional list of strings): The names of the availability zones you want to install your nodes on. If unset, the installer will use your default compute zone.

**NOTE:** The bootstrap node follows the `type`, `rootVolume`, `additionalNetworkIDs`, and `additionalSecurityGroupIDs` parameters from the `controlPlane` machine pool.

**NOTE:** Note when deploying with `Kuryr` there is an Octavia API loadbalancer VM that will not fulfill the Availability Zones restrictions due to Octavia lack of support for it. In addition, if Octavia only has the amphora provider instead of also the OVN-Octavia provider, all the OpenShift services will be backed up by Octavia Load Balancer VMs which will not fulfill the Availability Zone restrictions either.


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
    apiFloatingIP: 128.0.0.1
    cloud: mycloud
    defaultMachinePlatform:
      type: m1.s2.xlarge
    externalNetwork: external
    externalDNS:
      - "8.8.8.8"
      - "192.168.1.12"
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
    apiFloatingIP: 128.0.0.1
    cloud: mycloud
    defaultMachinePlatform:
      type: m1.s2.xlarge
    externalNetwork: external
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

## Image Overrides

The OpenShift installer pins the version of RHEL CoreOS and normally handles uploading the image to the target OpenStack instance.

If you want to download the image manually, see [CoreOS bootimages](../overview.md#coreos-bootimages) for more information
about bootimages.  This is useful, for example, to perform a disconnected installation.  To do this,
download the `qcow2` and host it at a custom location.  Then set the `openstack.clusterOSImage`
parameter field in the install config to point to that location.   The install process will
then use that mirrored image.
In all other respects the process will be consistent with the default.

**NOTE:** For this to work, the parameter value must be a valid http(s) URL.

**NOTE:** The optional `sha256` query parameter can be attached to the URL. This will force the installer to check the uncompressed image file checksum before uploading it into Glance.

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

## Custom Subnets

In the `install-config.yaml` file, the value of the `machinesSubnet` property is the subnet where the Kubernetes endpoints of the nodes in your cluster are published. The Ingress and API ports are created on this subnet, too. By default, the installer creates a network and subnet for these endpoints and ports. Alternatively, you can use a subnet of your own by setting the value of the `machinesSubnet` property to the UUID of an existing OpenStack subnet. To use this feature, you need to meet these requirements:

* The subnet that is used by `machinesSubnet` has DHCP enabled.
* The CIDR of `machinesSubnet` matches the CIDR of `networks.machineNetwork`.
* The installer user must have permission to create ports on this network, including ports with fixed IP addresses.

You should also be aware of the following limitations:

* If you plan to install a cluster that uses floating IPs, the `machinesSubnet` must be attached to a router that is connected to the `externalNetwork`.
* The installer will not create a private network or subnet for your OpenShift machines if the `machinesSubnet` is set in the `install-config.yaml`.
* By default, the API and Ingress VIPs use the .5 and .7 of your network CIDR. To prevent other services from taking the ports that are assigned to the API and Ingress VIPs, set the `apiVIP` and `ingressVIP` options in the `install-config.yaml` to addresses that are outside of the DHCP allocation pool.
* You cannot use the `externalDNS` property at the same time as a custom `machinesSubnet`. If you want to add a DNS to your cluster while using a custom subnet, [add it to the subnet in OpenStack](https://docs.openstack.org/neutron/rocky/admin/config-dns-res.html).

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

**NOTES:**
* Allowed address pairs won't be created for the additional networks.
* The additional networks attached to the Control Plane machine will also be attached to the bootstrap node.

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

**NOTE:** The additional security groups attached to the Control Plane machine will also be attached to the bootstrap node.

## Further customization

For customizing the installation beyond what is possible with `openshift-install`, refer to the [UPI (User Provided Infrastructure) documentation](./install_upi.md).
