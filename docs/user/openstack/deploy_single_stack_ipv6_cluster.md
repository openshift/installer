# Creating a single stack IPv6 cluster on OpenStack

## Table of Contents

- [Prerequisites](#prerequisites)
- [Creating Network for the cluster](#creating-network-for-the-cluster)
- [Creating IPv6 API and Ingress VIPs Ports for the cluster](#creating-ipv6-api-and-ingress-vips-ports-for-the-cluster)
- [Deploy OpenShift](#deploy-openshift)

## Prerequisites

* Installation with single stack IPv6 is only allowed when using one pre-created OpenStack IPv6 subnet.
* DNS must be configured in the Subnet.
* Add the IPv6 Subnet to a neutron router to provide router advertisements.
* The network MTU must accommodate the minimum MTU for IPv6, which is 1280, and OVN-Kubernetes encapsulation overhead, which is 100.
* API and Ingress VIPs ports needs to pre-created by the user and the addresses specified in the `install-config.yaml`.
* A local image registry needs to be pre-configured to mirror the images over IPv6.

Additional prerequisites are listed at the [OpenStack Platform Customization docs](./customization.md)

**Note**: Converting a dual-stack cluster to single stack IPv6 cluster is not supported with OpenStack.

## Creating Network for the cluster

You must create one network and add the IPv6 subnet. Here is an example:

```sh
$ openstack network create --project <project-name> --share --external --provider-physical-network <physical-network> --provider-network-type flat v6-network
$ openstack subnet create --project <project-name> v6-subnet --subnet-range fd2e:6f44:5dd8:c956::/64 --dhcp --dns-nameserver <dns-address> --network v6-network --ip-version 6 --ipv6-ra-mode stateful --ipv6-address-mode stateful
```

**Note**: using an IPv6 slaac subnet is not supported given a known [OpenStack issue](https://bugzilla.redhat.com/show_bug.cgi?id=2304331) that prevents DNS from working.

Given the above example uses a provider network, this network can be added to the router external gateway to enable external connectivity and router advertisements with the following command:
```sh
$ openstack router set --external-gateway v6-network <router-id>
```

**Note**: Any additional IPv6 Subnet that is used in the OpenShift cluster, should be added to a neutron router to provide router advertisements.

## Creating IPv6 API and Ingress VIPs Ports for the cluster

You must create the API and Ingress VIPs Ports with the following commands:

```sh
$ openstack port create api --network v6-network
$ openstack port create ingress --network v6-network
```

## Deploy OpenShift

Now that the Networking resources are pre-created you can deploy OpenShift. Here is an example of `install-config.yaml`:

```yaml
apiVersion: v1
baseDomain: mydomain.test
compute:
- name: worker
  platform:
    openstack:
      type: m1.xlarge
  replicas: 3
controlPlane:
  name: master
  platform:
    openstack:
      type: m1.xlarge
  replicas: 3
metadata:
  name: mycluster
networking:
  machineNetwork:
  - cidr: "fd2e:6f44:5dd8:c956::/64"
  clusterNetwork:
  - cidr: fd01::/48
    hostPrefix: 64
  serviceNetwork:
  - fd02::/112
platform:
  openstack:
    ingressVIPs: ['fd2e:6f44:5dd8:c956::383']
    apiVIPs: ['fd2e:6f44:5dd8:c956::9a']
    controlPlanePort:
      fixedIPs:
      - subnet:
          name: subnet-v6
      network:
        name: v6-network
imageContentSources:
- mirrors:
  - <mirror>
  source: quay.io/openshift-release-dev/ocp-v4.0-art-dev
- mirrors:
  - <mirror>
  source: registry.ci.openshift.org/ocp/release
additionalTrustBundle: |
<certificate-of-the-mirror>
```
There are important things to note:

The subnets under `platform.openstack.controlPlanePort.fixedIPs` can contain both id or name. The same applies to the network `platform.openstack.controlPlanePort.network`.

The image content sources contains the details of the mirror to be used. Please follow the docs to configure a [local image registry](https://docs.openshift.com/container-platform/4.16/installing/disconnected_install/installing-mirroring-creating-registry.html).