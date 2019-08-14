# OpenStack Platform Customization

In the OpenShift Installer `install-config.yaml` you can set the following options regarding the OpenStack platform:

- `machines.platform.openstack.region` - The OpenStack region where the cluster will get created
- `machines.platform.openstack.cloud` - Name of the OpenStack cloud to use from clouds.yaml
- `machines.platform.openstack.externalNetwork` - The OpenStack external network name to be used for installation
- `machines.platform.openstack.computeFlavor` - The OpenStack compute flavor to use for master servers
- `machines.platform.openstack.lbFloatingIP` - Existing Floating IP to associate with API loadbalancer
- `machines.platform.openstack.trunkSupport` - Whether OpenStack ports can be trunked. True or False
- `machines.platform.openstack.octaviaSupport` - Whether OpenStack supports Octavia. True of False
- `machines.platform.openstack.defaultMachinePlatform` - (optional) The default configuration used when installing on OpenStack for machine pools

For more technical definitions, see the [go docs](https://godoc.org/github.com/openshift/installer/pkg/types/openstack#Platform).

## Examples

The example `install-config.yaml` below showcases all the possible OpenStack customizations.

```yaml
apiVersion: v1
baseDomain: example.com
clusterID: os-test
controlPlane:
  name: master
  platform: {}
  replicas: 3
compute:
- name: worker
  platform:
    openstack:
      type: ml.large
  replicas: 3
metadata:
  name: example
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  machineCIDR: 10.0.0.0/16
  serviceNetwork:
  - 172.30.0.0/16
  networkType: OpenShiftSDN
platform:
  openstack:
    region: region1
    cloud: mycloud
    externalNetwork: external
    computeFlavor: m1.xlarge
    lbFloatingIP: 128.0.0.1
    trunkSupport: false
    octaviaSupport: false
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```
