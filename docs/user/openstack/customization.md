# OpenStack Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, OpenStack-specific properties.

## Table of Contents

* [Cluster-scoped properties](#cluster-scoped-properties)
* [Machine pools](#machine-pools)
* [Examples](#examples)
  * [Minimal](#minimal)
  * [Custom-machine-pools](#custom-machine-pools)
* [Image Overrides](#image-overrides)
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

## Machine pools

* `type` (optional string): The OpenStack flavor name for machines in the pool.
* `rootVolume` (optional object): Defines the root volume for instances in the machine pool. The instances use ephemeral disks if not set.
  * `size` (required integer): Size of the root volume in GB.
  * `type` (required string): The volume pool to create the volume from.

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
    region: region1
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
    region: region1
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

If the user wants to reuse an existing Glance image without any uploading of binary data, then it is possible to set `clusterOSImage` install config parameter that specifies the Glance image name. In this case no new Glance images will be created, and the image will stay when the cluster is destroyed. In other words, if `clusterOSImage` is not an http(s) URL, then the installer will look into Glance for an image with that name.

Example:

```yaml
platform:
  openstack:
      clusterOSImage: my-rhcos
```

## Further customization

For customizing the installation beyond what is possible with `openshift-install`, refer to the [UPI (User Provided Infrastructure) documentation](./install_upi.md).
