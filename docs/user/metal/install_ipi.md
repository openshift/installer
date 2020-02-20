# Bare Metal IPI (Installer Provisioned Infrastructure) Overview

Current Status: **Experimental**

This document discusses the installer support for an IPI (Installer Provisioned
Infrastructure) install for bare metal hosts.  This includes platform support
for the management of bare metal hosts, as well as some automation of DNS and
load balancing to bring up the cluster.

The upstream project that provides Kubernetes-native management of bare metal
hosts is [metal3.io](http://metal3.io).

For UPI (User Provisioned Infrastructure) based instructions for bare metal
deployments, see [install_upi.md](install_upi.md).

## Prerequisites

### Network Requirements

It is assumed that all hosts have at least 2 NICs, used for the following
purposes:

* **NIC #1 - External Network**
  * This network is the main network used by the cluster, including API traffic
    and application traffic.
  * ***DHCP***
    * External DHCP is assumed on this network.  It is **strongly** recommended
      to set up DHCP reservations for each of the hosts in the cluster to
      ensure that they retain stable IP addresses.
    * A pool of dynamic addresses should also be available on this network, as
      the provisioning host and temporary bootstrap VM will also need addresses
      on this network.
  * ***NTP***
    * A time source must be accessible from this network.
  * ***Reserved VIPs (Virtual IPs)*** - 3 IP addresses must be reserved on this
    network for use by the cluster.  Specifically, these IPs will serve the
    following purposes:
    * API - This IP will be used to reach the cluster API.
    * Ingress - This IP will be used by cluster ingress traffic
    * DNS - This IP will be used internally by the cluster for automating
      internal DNS requirements.
  * ***External DNS*** - While the cluster automates the internal DNS
    requirements, two external DNS records must be created in whatever DNS
    server is appropriate for this environment.
    * `api.<cluster-name>.<base-domain>` - pointing to the API VIP
    * `*.apps.<cluster-name>.<base-domain>` - pointing to the Ingress VIP

* **NIC #2 - Provisioning Network**
  * A private, non-routed network, used for PXE based provisioning.
  * DHCP is automated for this network.
  * Addressing for this network is currently hard coded as `172.22.0.0/24`, but
    will be made configurable in the future.

* **Out-of-band Management Network**
  * Servers will typically have an additional NIC used by the onboard
    management controllers (BMCs).  These BMCs must be accessible and routed to
    the host.

### Provisioning Host

The installer must be run from a host that is attached to the same networks as
the cluster, as described in the previous section.  We refer to this host as
the *provisioning host*.  The easiest way to provide a provisioning host is to
use one of the hosts that is intended to later become a worker node in the same
cluster.  That way it is already connected to the proper networks.

It is recommended that the provisioning host be a bare metal host, as it must be
able to use libvirt to launch the OpenShift bootstrap VM locally.

### Supported Hardware

The architecture is intended to support a wide variety of hardware.  This was
one of the reasons Ironic is used as an underlying technology.  However, so far
development and testing has focused on PXE based provisioning using IPMI for
out-of-band management of hosts.  Other provisioning approaches will be added,
tested, and documented over time.

## Installation Process

Once an environment has been prepared according to the documented
pre-requisites, the install process is the same as other IPI based platforms.

`openshift-install create cluster`

However, it is recommended to prepare an `install-config.yaml` file in advance,
containing all of the details of the bare metal hosts to be provisioned.

### Install Config

The `install-config.yaml` file requires some additional details.  Most of the
information is teaching the installer and the resulting cluster enough about
the available hardware so that it is able to fully manage it.

Here is an example `install-config.yaml` with the required `baremetal` platform
details.

**IMPORTANT NOTE:** The current install configuration for the `baremetal`
platform should be considered experimental and still subject to change without
backwards compatibility.  In particular, some items likely to change soon
include:

* The `hardwareProfile` is currently exposed as a way to allow specifying
  different hardware parameters for deployment.  By default, we will deploy
  RHCOS to the first disk, but that may not be appropriate for all hardware.
  The `hardwareProfile` is the field we have available to change that.  This
  interface is subject to change.  In the meantime, hardware profiles can be
  found here:
  https://github.com/metal3-io/baremetal-operator/blob/master/pkg/hardware/profile.go#L48

```yaml
apiVersion: v1
baseDomain: test.metalkube.org
metadata:
  name: ostest
networking:
  machineNetwork:
  - cidr: 192.168.111.0/24
compute:
- name: worker
  replicas: 1
controlPlane:
  name: master
  replicas: 3
  platform:
    baremetal: {}
platform:
  baremetal:
    apiVIP: 192.168.111.5
    ingressVIP: 192.168.111.4
    dnsVIP: 192.168.111.3
    hosts:
      - name: openshift-master-0
        role: master
        bmc:
          address: ipmi://192.168.111.1:6230
          username: admin
          password: password
        bootMACAddress: 00:11:07:4e:f6:68
        hardwareProfile: default
      - name: openshift-master-1
        role: master
        bmc:
          address: ipmi://192.168.111.1:6231
          username: admin
          password: password
        bootMACAddress: 00:11:07:4e:f6:6c
        hardwareProfile: default
      - name: openshift-master-2
        role: master
        bmc:
          address: ipmi://192.168.111.1:6232
          username: admin
          password: password
        bootMACAddress: 00:11:07:4e:f6:70
        hardwareProfile: default
      - name: openshift-worker-0
        role: worker
        bmc:
          address: ipmi://192.168.111.1:6233
          username: admin
          password: password
        bootMACAddress: 00:11:07:4e:f6:71
        hardwareProfile: default
pullSecret: ...
sshKey: ...
```

## Work in Progress

Integration of the `baremetal` platform is still a work-in-progress across
various parts of OpenShift.  This section discusses key items that are not yet
fully integrated, and their workarounds.

Note that once this work moves into the `openshift/installer` repository, new
issues will get created or existing issues will be moved to track these gaps
instead of the leaving the existing issues against the KNI fork of the installer.

### `destroy cluster` support

`openshift-install destroy cluster` is not supported for the `baremetal`
platform.

https://github.com/openshift-metal3/kni-installer/issues/74

### install gather not supported

When an installation fails, `openshift-install` will attempt to gather debug
information from hosts.  This is not yet supported by the `baremetal` platform.

https://github.com/openshift-metal3/kni-installer/issues/79

### Provisioning subnet not fully configurable

There are some install-config parameters to control templating of the provisioning
network configuration, but fully supporting alternative subnets for the
provisioning network is incomplete.

https://github.com/openshift/installer/issues/2091
