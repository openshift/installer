# OpenStack Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, OpenStack-specific properties.

## Cluster-scoped properties

* `cloud` (required string): The name of the OpenStack cloud to use from `clouds.yaml`.
* `computeFlavor` (required string): The OpenStack flavor to use for compute and control-plane machines.
    This is currently required, but has lower precedence than [the `type` property](#machine-pools) on [the `compute` and `controlPlane` machine-pools](../customization.md#platform-customization).
* `externalNetwork` (required string): The OpenStack external network name to be used for installation.
* `lbFloatingIP` (required string): Existing Floating IP to associate with the API load balancer.
* `octaviaSupport` (optional string): Whether OpenStack supports Octavia (`1` for true or `0` for false)
* `region` (required string): The OpenStack region where the cluster will be created.
* `trunkSupport` (optional string): Whether OpenStack ports can be trunked (`1` for true or `0` for false)

## Machine pools

* `type` (optional string): The OpenStack flavor name for machines in the pool.

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
