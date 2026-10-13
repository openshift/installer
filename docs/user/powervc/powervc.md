# PowerVC Platform Customization

Beyond the [platform-agnostic `install-config.yaml` properties](../customization.md#platform-customization), the installer supports additional, PowerVC-specific properties.

## Table of Contents

- [Example](#example)

## Example

An example `install-config.yaml` is shown below.

```yaml
apiVersion: v1
baseDomain: domain.net
compute:
- architecture: ppc64le
  hyperthreading: Enabled
  name: worker
  platform:
    powervc:
      zones:
        - s1022
  replicas: 3
controlPlane:
  architecture: ppc64le
  hyperthreading: Enabled
  name: master
  platform:
    powervc:
      zones:
        - s1022
  replicas: 3
metadata:
  name: rdr-openstack
networking:
  clusterNetwork:
  - cidr: 10.116.0.0/14
    hostPrefix: 23
  machineNetwork:
  - cidr: 10.130.32.0/20
  networkType: OVNKubernetes
  serviceNetwork:
  - 172.30.0.0/16
platform:
  powervc:
    loadBalancer:
      type: UserManaged
    apiVIPs:
    - 10.130.41.250
    cloud: cloud-name
    clusterOSImage: rhcos-9.6.20250826-1-openstack-ppc64le
    defaultMachinePlatform:
      type: flavor
    ingressVIPs:
    - 10.130.41.250
    controlPlanePort:
      fixedIPs:
        - subnet:
            id: deadbeef-0000-1111-2222-cafefade1234
credentialsMode: Passthrough
pullSecret: '{"auths": ...}'
sshKey: ssh-ed25519 AAAA...
```

`loadBalancer` can only be of type `UserManaged`.  This means that the user must create a HAProxy server
which acts as a load balancer.  It must be updated during the cluster creation as new VMs are created.  Both
`apiVIPs` and `ingressVIPs` should be the same HAProxy server IP.

`clusterOSImage` is an image created in advance by the `powervc-image` tool.  The source can only be a `ova.gz`

`credentialsMode` can only be `Passthrough`.
