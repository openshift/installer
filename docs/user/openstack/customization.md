# OpenStack Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, OpenStack-specific properties.

## Table of Contents

* [Cluster-scoped properties](#cluster-scoped-properties)
* [Machine pools](#machine-pools)
* [Examples](#examples)
  * [Minimal](#minimal)
  * [Custom-machine-pools](#custom-machine-pools)
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

## Further customization

For customizing the installation beyond what is possible with `openshift-install`, refer to the [UPI (User Provided Infrastructure) documentation](./install_upi.md).
