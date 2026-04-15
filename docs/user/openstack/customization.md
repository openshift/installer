# OpenStack Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, OpenStack-specific properties.

## Table of Contents

- [OpenStack Platform Customization](#openstack-platform-customization)
  - [Table of Contents](#table-of-contents)
  - [Cluster-scoped properties](#cluster-scoped-properties)
  - [Machine pools](#machine-pools)
  - [Bootstrap Flavor Optimization](#bootstrap-flavor-optimization)
  - [Examples](#examples)
    - [Minimal](#minimal)
    - [Custom machine pools](#custom-machine-pools)
    - [Bootstrap flavor](#bootstrap-flavor)
  - [Image Overrides](#image-overrides)
  - [Custom Subnets](#custom-subnets)
  - [Additional Networks](#additional-networks)
  - [Additional Security Groups](#additional-security-groups)
  - [Cloud Provider configuration](#cloud-provider-configuration)
  - [Further customization](#further-customization)

## Cluster-scoped properties

* `cloud` (required string): The name of the OpenStack cloud to use from `clouds.yaml`.
* `bootstrapFlavor` (optional string): The OpenStack flavor to use for the bootstrap machine. When not specified, inherits from `controlPlane.type` (see [Machine pools](#machine-pools)).
* `computeFlavor` (deprecated string): The OpenStack flavor to use for compute and control-plane machines.
* `externalDNS` (optional list of strings): The IP addresses of DNS servers to be used for the DNS resolution of all instances in the cluster. The total number of dns servers supported by an instance is three. That total includes any dns server provided by the underlying openstack infrastructure.
* `externalNetwork` (optional string): Name of external network the installer will use to provide access to the cluster. If defined, a floating IP from this network will be created and associated with the bootstrap node to facilitate debugging and connection to the bootstrap node during installation. The `apiFloatingIP` property is a floating IP address selected from this network.
* `apiFloatingIP` (optional string): Address of existing Floating IP from externalNetwork the installer will associate with the OpenShift API. This property is only valid if externalNetwork is defined. If externalNetwork is not defined, the installer will throw an error.
* `ingressFloatingIP` (optional string): Address of an existing Floating IP from externalNetwork the installer will associate with the ingress port. This property is only valid if externalNetwork is defined. If externalNetwork is not defined, the installer will throw an error.
* `octaviaSupport` (deprecated string): Whether OpenStack supports Octavia (`1` for true or `0` for false)
* `region` (deprecated string): The OpenStack region where the cluster will be created. Currently this value is not used by the installer.
* `trunkSupport` (deprecated string): Whether OpenStack ports can be trunked (`1` for true or `0` for false)
* `clusterOSImage` (optional string): Either a URL with `http(s)` or `file` scheme to override the default OS image for cluster nodes or an existing Glance image name.
* `clusterOSImageProperties` (key-value pairs): a list of properties to be added to the installer-uploaded ClusterOSImage in Glance. The default is to not set any properties. `clusterOSImageProperties` is ignored when `clusterOSImage` points to an existing image in Glance.
* `apiVIPs` (optional array of strings): IP address on the machineNetwork that will be assigned to the API VIP. If more than one are set, it must be one IPv4 and one IPv6.
* `ingressVIPs` (optional array of strings): IP address on the machineNetwork that will be assigned to the ingress VIP. If more than one are set, it must be one IPv4 and one IPv6.
* `controlPlanePort` (optional object): the UUID and/or Name of an OpenStack Network and its Subnets where to install the nodes of the cluster onto. For more information on how to install with a custom subnet, see the [custom subnets](#custom-subnets) section of the docs.
* `defaultMachinePlatform` (optional object): Default [OpenStack-specific machine pool properties](#machine-pools) which apply to [machine pools](../customization.md#machine-pools) that do not define their own OpenStack-specific properties.

## Machine pools

* `additionalNetworkIDs` (optional list of strings): IDs of additional networks for machines.
* `additionalSecurityGroupIDs` (optional list of strings): IDs of additional security groups for machines.
* `serverGroupPolicy` (optional string): Server group policy to apply to the group that will contain the machines in the pool. Defaults to "soft-anti-affinity". Allowed values are "affinity", "soft-affinity", "anti-affinity", "soft-anti-affinity".
  * It is not possible to change a server group policy or a server's affiliation to a group after creation
  * A strict "affinity" policy prevents migrations, and therefore affects OpenStack upgrades
  * An additional OpenStack host is needed when migrating instances with a strict "anti-affinity" policy
* `type` (optional string): The OpenStack flavor name for machines in the pool.
* `rootVolume` (optional object): Defines the root volume for instances in the machine pool. The instances use ephemeral disks if not set.
  * `size` (required integer): Size of the root volume in GB. Must be set to at least 25. For production clusters, this must be at least 100.
  * `type` (deprecated string): The volume pool to create the volume from. It was replaced by `types`.
  * `types` (required list of strings): The volume pool to create the volume from. If compute `zones` are defined with more than one type, the number of zones must match the number of types.
  * `zones` (optional list of strings): The names of the availability zones you want to install your root volumes on. If unset, the installer will use your default volume zone.
    If compute `zones` contains at least one value, `rootVolume.zones` must also contain at least one value.
    Indeed, when a machine is created with a compute availability zone and a storage root volume with no specified `rootVolume.availabilityZone`, [CAPO](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/9d183bd479fe9aed4f6e7ac3d5eee46681c518e7/pkg/cloud/services/compute/instance.go#L439-L442) will use the compute AZ for the volume AZ.
    This can be problematic if the AZ doesn't exist in Cinder, therefore we enforce that `rootVolume.zones` to be set if `zones` is set.
* `zones` (optional list of strings): The names of the availability zones you want to install your nodes on. If unset, the installer will use your default compute zone.

> **Note**
> The bootstrap node follows the `type`, `rootVolume`, `additionalNetworkIDs`, and `additionalSecurityGroupIDs` parameters from the `controlPlane` machine pool. The bootstrap flavor (`type`) can be overridden independently using the cluster-scoped [`bootstrapFlavor`](#cluster-scoped-properties) field.

> **Note**
> Note when deploying the control-plane machines with `rootVolume`, it is highly suggested to use an [additional ephemeral disk dedicated to etcd](./etcd-ephemeral-disk.md).

## Bootstrap Flavor Optimization

The bootstrap machine is a **temporary** instance that exists only during cluster installation. Understanding its lifecycle and resource requirements allows operators to select a smaller, cheaper flavor for it, reducing infrastructure costs during the installation phase.

### Bootstrap Machine Lifecycle

The bootstrap machine goes through the following lifecycle:

1. **Created** at the start of `openshift-install create cluster`
2. **Hosts** the temporary Kubernetes control plane and serves Ignition configs to control plane nodes
3. **Forms** the initial etcd cluster with the control plane nodes
4. **Transfers** control to the production control plane once it becomes self-hosting
5. **Destroyed** automatically once the control plane is fully operational (typically 15–30 minutes after installation begins)

Because the bootstrap machine exists only for the duration of the installation process, it does not need to match the production-grade flavor used for permanent control plane nodes.

### Cost Optimization Benefits

By specifying a smaller flavor for the bootstrap machine via the `bootstrapFlavor` field, operators can reduce the cost of the bootstrap phase by **30–50%** compared to using the same flavor as the control plane:

- Control plane nodes run indefinitely and require production-grade resources for etcd, the API server, and other components.
- The bootstrap machine runs a temporary, lightweight workload and is active for only 15–30 minutes.
- Using a right-sized bootstrap flavor avoids paying for idle compute capacity on an oversized instance during installation.

For example, if the control plane uses a flavor with 16 vCPUs and 64 GB RAM, specifying a 4 vCPU / 16 GB RAM bootstrap flavor can reduce the per-hour cost of the bootstrap instance by more than 50%, resulting in meaningful savings especially in environments where installations are performed frequently (e.g., CI/CD pipelines, automated cluster provisioning).

### Minimum Recommended Resources for Bootstrap

The bootstrap machine must have enough resources to run the temporary control plane services. The minimum recommended resources are:

| Resource | Minimum |
|----------|---------|
| vCPUs    | 4       |
| RAM      | 16 GB   |
| Disk     | 100 GB  |

> **Note**
> These are the **minimum** recommended values. Using a flavor that does not meet these requirements may cause the bootstrap process to fail or time out. The installer validates that the specified `bootstrapFlavor` meets at least the control plane flavor requirements.

### Choosing a Bootstrap Flavor

When selecting a bootstrap flavor, consider the following guidance:

- **Do** use a flavor that meets the minimum requirements above (4 vCPU, 16 GB RAM, 100 GB disk).
- **Do** choose a smaller flavor than your control plane nodes to reduce cost during installation.
- **Do not** use a baremetal flavor for the bootstrap machine unless your environment requires it — virtual flavors are sufficient and more cost-effective.
- **Do not** reduce the bootstrap flavor below the minimums, as this can cause installation failures due to resource exhaustion during the bootstrapping phase.

To configure a bootstrap flavor, set the `bootstrapFlavor` field in the cluster-scoped platform properties. See the [Bootstrap flavor](#bootstrap-flavor) example for a complete `install-config.yaml`.

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
        types:
        - performance
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

### Bootstrap flavor

The `bootstrapFlavor` field allows you to use a different (typically smaller) OpenStack flavor
for the temporary bootstrap machine, which is only needed during cluster installation and is
destroyed once the cluster is self-hosting. This lets you optimize resource usage and costs by
running the bootstrap node on a cheaper flavor while keeping the control plane nodes on a
larger, production-grade flavor.

#### With explicit bootstrapFlavor (cost-optimized)

The following example configures `m1.xlarge` for control plane machines and uses the smaller
`m1.medium` flavor for the bootstrap machine. Since the bootstrap node is temporary and hosts
lighter workloads than production control plane nodes, a smaller flavor is sufficient:

```yaml
apiVersion: v1
baseDomain: example.com
controlPlane:
  name: master
  platform:
    openstack:
      type: m1.xlarge  # Production-grade flavor for persistent control plane nodes
  replicas: 3
compute:
- name: worker
  platform:
    openstack:
      type: m1.large
  replicas: 3
metadata:
  name: test-cluster
platform:
  openstack:
    cloud: mycloud
    externalNetwork: external
    apiFloatingIP: 128.0.0.1
    bootstrapFlavor: m1.medium  # Smaller flavor for the temporary bootstrap node
                                # Reduces cost since bootstrap is destroyed after install
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

#### Without bootstrapFlavor (inheritance behavior)

When `bootstrapFlavor` is not specified, the bootstrap node inherits its flavor from
`controlPlane.platform.openstack.type` (or `defaultMachinePlatform.type` if the control
plane does not define its own flavor). This is the default behavior:

```yaml
apiVersion: v1
baseDomain: example.com
controlPlane:
  name: master
  platform:
    openstack:
      type: m1.xlarge  # Bootstrap node will also use m1.xlarge (inherited)
  replicas: 3
compute:
- name: worker
  platform:
    openstack:
      type: m1.large
  replicas: 3
metadata:
  name: test-cluster
platform:
  openstack:
    cloud: mycloud
    externalNetwork: external
    apiFloatingIP: 128.0.0.1
    # bootstrapFlavor is not set — bootstrap node inherits controlPlane.platform.openstack.type
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

> **Note**
> For this to work, the parameter value must be a valid http(s) URL.

> **Note**
> The optional `sha256` query parameter can be attached to the URL. This will force the installer to check the uncompressed image file checksum before uploading it into Glance.

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

In the `install-config.yaml` file, the value of the `controlPlanePort` property contains the Name and/or UUID of the network and subnet(s) where the Kubernetes endpoints of the nodes in your cluster are published. The Ingress and API ports are created on the subnets, too. By default, the installer creates a network and subnet for these endpoints and ports. Alternatively, you can use a OpenStack network containing one subnet or two, in case of dual-stack, of your own by specifying their and/or in the `controlPlanePort` property. To use this feature, you need to meet these requirements:

* Any subnet used by `controlPlanePort` have DHCP enabled.
* The CIDR of any subnet listed in `controlPlanePort.fixedIPs` matches the CIDRs listed on `networks.machineNetwork`.
* When using dual-stack or single stack IPv6 Network the api and ingress Ports needs to be pre-created by the user. Also, the installer user must have permission to add tags and security groups to those pre-created Ports. The value of the fixed IPs of the Ports needs to be specified at the `apiVIPs` and `ingressVIPs` options in the `install-config.yaml`.
* If not using dual-stack or single stack IPv6, the installer user must have permission to create ports on this network, including ports with fixed IP addresses.

You should also be aware of the following limitations:

* If you plan to install a cluster that uses floating IPs, the `controlPlanePort` must be attached to a router that is connected to the `externalNetwork`.
* The installer will not create a private network or subnet for your OpenShift machines if the `controlPlanePort` is set in the `install-config.yaml`.
* By default when not using dual-stack or single stack IPv6, the API and Ingress VIPs use the .5 and .7 of your network CIDR. To prevent other services from taking the ports that are assigned to the API and Ingress VIPs, set the `apiVIP` and `ingressVIP` options in the `install-config.yaml` to addresses that are outside of the DHCP allocation pool.
* You cannot use the `externalDNS` property at the same time as a custom `controlPlanePort`. If you want to add a DNS to your cluster while using a custom subnet, [add it to the subnet in OpenStack](https://docs.openstack.org/neutron/rocky/admin/config-dns-res.html).

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

> **Note**
> Allowed address pairs won't be created for the additional networks.

> **Note**
> The additional networks attached to the Control Plane machine will also be attached to the bootstrap node.

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

> **Note**
> The additional security groups attached to the Control Plane machine will also be attached to the bootstrap node.

## Cloud Provider configuration

You may want to modify cloud provider configuration in order to make it work with your OpenStack cloud. This is possible if you'll let the installer generate the manifests before running the installation:

```sh
openshift-install --dir <directory> create manifests
```

Then modify the manifest containing the [cloud provider configuration](https://v1-18.docs.kubernetes.io/docs/concepts/cluster-administration/cloud-providers/#cloud-conf):

```sh
vi <directory>/manifests/cloud-provider-config.yaml
```

As an example in order to tweak support for LoadBalancer Services you can modify options regarding Octavia configuration in the `[LoadBalancer]` section of `config` key. In particular:

* `enabled = false` disables the Octavia integration.
* `lb-provider = <"amphora" or "ovn">` lets you choose the Octavia provider to use when creating load balancers. Please note that setting "ovn" requires setting `lb-method = SOURCE_IP_PORT` as this is the only method supported by OVN.
* `floating-network-id = <uuid>` is required to be set if your OpenStack cluster has multiple external networks. The network set here will be used by cloud provider to create floating IPs on.

After saving the file you can continue the installation normally:

```sh
openshift-install --dir <directory> create cluster
```

## Further customization

For customizing the installation beyond what is possible with `openshift-install`, refer to the [UPI (User Provided Infrastructure) documentation](./install_upi.md).
