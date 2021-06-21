# Bare Metal IPI (Installer Provisioned Infrastructure) Overview

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

You have the choice of a single or dual NIC setup, depending on whether
you would like to use PXE/DHCP-based provisioning or not. Please note
that disabling the provisioning network means that host BMC's must be
accessible over the external network which may not be desirable.

* **NIC #1 - External Network**
  * This network is the main network used by the cluster, including API traffic
    and application traffic.
  * ***DHCP***
    * External DHCP is assumed on this network. Hosts *must* have stable
      IP addresses, therefore you should set up DHCP reservations for
      each of the hosts in the cluster. The addresses assigned by DHCP
      need to be in the same subnet as the Virtual IPs discussed below.
    * A pool of dynamic addresses should also be available on this network, as
      the provisioning host and temporary bootstrap VM will also need addresses
      on this network.
  * ***NTP***
    * A time source must be accessible from this network.
  * ***Reserved VIPs (Virtual IPs)*** - 3 IP addresses must be reserved on this
	network for use by the cluster. These Virtual IPs are managed using VRRP
	(v2 for IPv4 and v3 for IPv6). Specifically, these IPs will serve the
    following purposes:
    * API - This IP will be used to reach the cluster API.
    * Ingress - This IP will be used by cluster ingress traffic
      internal DNS requirements.
  * ***External DNS*** - While the cluster automates the internal DNS
    requirements, two external DNS records must be created in whatever DNS
    server is appropriate for this environment.
    * `api.<cluster-name>.<base-domain>` - pointing to the API VIP
    * `*.apps.<cluster-name>.<base-domain>` - pointing to the Ingress VIP

* **NIC #2 - Provisioning Network (optional) **
  * A private network used for PXE based provisioning.
  * You can specify `provisioningNetworkInterface` to indicate which
    interface is connected to this network on the control plane nodes.
    If not specified the interface is derived from the bootMacAddress.
    If set, all the control plane nodes must have this interface.
  * The provisioning network may be "Managed" (default), "Unmanaged," or
    "Disabled."
  * In managed mode, DHCP and TFTP are configured to run in the cluster. In
    unmanaged mode, TFTP is still available but you must configure DHCP
    externally.
  * Addressing for this network defaults to `172.22.0.0/24`, but is
    configurable by setting the `provisioningNetworkCIDR` option.
  * Two IP's are required to be available for use, one for the bootstrap
    host, and one as a provisioning IP in the running cluster. By
    default, these are the 2nd and 3rd addresses in the
    `provisioningNetworkCIDR` (e.g. 172.22.0.2, and 172.22.0.3).
  * To specify the name of the provisioning network interface,
    set the `provisioningNetworkInterface` option. This is the network interface
    on a master that is connected to the provisioning network.

* **Out-of-band Management Network**
  * Servers will typically have an additional NIC used by the onboard
    management controllers (BMCs).  These BMCs must be accessible and routed to
    the host.

When the Virtual IPs are managed using multicast (VRRPv2 or VRRPv3), there is a
limitation for 255 unique virtual routers per multicast domain. In case you
have pre-existing virtual routers using the standard IPv4 or IPv6 multicast
groups, you can learn the VIPs the installation will choose by running the
following command:

    $ podman run quay.io/openshift/origin-baremetal-runtimecfg:TAG vr-ids cnf10
    APIVirtualRouterID: 147
    DNSVirtualRouterID: 158
    IngressVirtualRouterID: 2

Where `TAG` is the release you are going to install, e.g., 4.5. Let's see another example:

    $ podman run quay.io/openshift/origin-baremetal-runtimecfg:TAG vr-ids cnf11
    APIVirtualRouterID: 228
    DNSVirtualRouterID: 239
    IngressVirtualRouterID: 147

In the example output above you can see that installing two clusters in the
same multicast domain with names `cnf10` and `cnf11` would lead to a conflict.
You should also take care that none of those are taken by other independent
VRRP virtual routers running in the same broadcast domain.

### Provisioning Host

The installer must be run from a host that is attached to the same networks as
the cluster, as described in the previous section.  We refer to this host as
the *provisioning host*.  The easiest way to provide a provisioning host is to
use one of the hosts that is intended to later become a worker node in the same
cluster.  That way it is already connected to the proper networks.

It is recommended that the provisioning host be a bare metal host, as it must be
able to use libvirt to launch the OpenShift bootstrap VM locally. Additionally,
the installer creates a directory backed libvirt storage pool in the
`/var/lib/libvirt/openshift-images` directory. Sufficient disk space must be
available in the directory to host the bootstrap VM volume.

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

Note for baremetal the installer must be built with both `libvirt` and
`baremetal` tags - in releases such a binary is included,
named `openshift-baremetal-install`

The installer supports interactive mode, but it is recommended to prepare an
`install-config.yaml` file in advance, containing all of the details of the
bare metal hosts to be provisioned.

### Install Config

The `install-config.yaml` file requires some additional details.  Most of the
information is teaching the installer and the resulting cluster enough about
the available hardware so that it is able to fully manage it. There are
[additional customizations](customization_ipi.md) possible.

Here is an example `install-config.yaml` with the required `baremetal` platform
details.

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
    hosts:
      - name: openshift-master-0
        role: master
        bmc:
          address: ipmi://192.168.111.1:6230
          username: admin
          password: password
        bootMACAddress: 00:11:07:4e:f6:68
        rootDeviceHints:
          minSizeGigabytes: 20
        bootMode: legacy
      - name: openshift-master-1
        role: master
        bmc:
          address: ipmi://192.168.111.1:6231
          username: admin
          password: password
        bootMACAddress: 00:11:07:4e:f6:6c
        rootDeviceHints:
          minSizeGigabytes: 20
        bootMode: UEFI
      - name: openshift-master-2
        role: master
        bmc:
          address: ipmi://192.168.111.1:6232
          username: admin
          password: password
        bootMACAddress: 00:11:07:4e:f6:70
        rootDeviceHints:
          minSizeGigabytes: 20
      - name: openshift-worker-0
        role: worker
        bmc:
          address: ipmi://192.168.111.1:6233
          username: admin
          password: password
        bootMACAddress: 00:11:07:4e:f6:71
        rootDeviceHints:
          minSizeGigabytes: 20
pullSecret: ...
sshKey: ...
```

#### Required Inputs

| Parameter | Default | Description |
| --- | --- | --- |
`hosts` | | Details about bare metal hosts to use to build the cluster. See below for more details. |
`defaultMachinePlatform` | | The default configuration used for machine pools without a platform configuration. |
`apiVIP` | `api.<clusterdomain>` | The VIP to use for internal API communication. |
`ingressVIP` | `test.apps.<clusterdomain>` | The VIP to use for ingress traffic. |

##### VIP Settings

The `apiVIP` and `ingressVIP` settings must either be provided or
pre-configured in DNS so that the default names resolve correctly (see
the defaults in the table above).

##### Describing Hosts

The `hosts` parameter is a list of separate bare metal assets that
should be used to build the cluster. The number of assets must be at least greater or equal to the sum of the configured `ControlPlane` and `compute` `Replicas`.

| Name | Default | Description |
| --- | --- | --- |
| `name` | | The name of the `BareMetalHost` resource to associate with the details. It must be unique. |
| `role` | | Either `master` or `worker`. |
| `bmc` | | Connection details for the baseboard management controller. See below for details. |
| `bootMACAddress` | | The MAC address of the NIC the host will use to boot on the provisioning network. It must be unique. |
| `rootDeviceHints` | | How to choose the target disk for the OS during provisioning - for more details see [upstream docs](https://github.com/metal3-io/baremetal-operator/blob/master/docs/api.md). |
| `bootMode` | `UEFI` | Choose `legacy` (BIOS) or `UEFI` mode for booting. Use `UEFISecureBoot` to enable UEFI and secure boot on the server. |

The `bmc` parameter for each host is a set of values for accessing the
baseboard management controller in the host.

| Name | Default | Description |
| --- | --- | --- |
| `username` | | The username for authenticating to the BMC |
| `password` | | The password associated with `username`. |
| `address` | | The URL for communicating with the BMC controller, based on the provider being used. See [BMC Addressing](#bmc-addressing) for details. It must be unique. |

##### BMC Addressing

The `address` field for each `bmc` entry is a URL with details for
connecting to the controller, including the type of controller in the
URL scheme and its location on the network.

IPMI hosts use `ipmi://<host>:<port>`. An unadorned `<host>:<port>` is
also accepted. If the port is omitted, the default of 623 is used.

Dell iDRAC hosts use `idrac://` (or `idrac+http://` to disable TLS).

Fujitsu iRMC hosts use `irmc://<host>:<port>`, where `<port>` is
optional if using the default.

For Redfish, use `redfish://` (or `redfish+http://` to disable
TLS). The hostname (or IP address) and the path to the system ID are
both required.  For example
`redfish://myhost.example/redfish/v1/Systems/System.Embedded.1` or
`redfish://myhost.example/redfish/v1/Systems/1`

To use virtual media instead of PXE for attaching the provisioning
image to the host, use `redfish-virtualmedia://` or `idrac-virtualmedia://`

Please note that when the provisioning network is disabled, the only
supported BMC's are virtual media.

## Known Issues

### `destroy cluster` support

`openshift-install destroy cluster` is not supported for the `baremetal`
platform.

https://github.com/openshift/installer/issues/2005

## Troubleshooting

General troubleshooting for OpenShift installations can be found
[here](/docs/user/troubleshooting.md).

### Bootstrap

The bootstrap VM by default runs on the same host as the installer. This
bootstrap VM runs the Ironic services needed to provision the control
plane. Ironic being available is dependent on having successfully
downloaded the machine OS and Ironic agent images. In some cases, this
may fail, and the installer will report a timeout waiting for the Ironic
API.

To login to the bootstrap VM, you will need to ssh to the VM using the
`core` user, and the SSH key defined in your install config.

The VM obtains an IP address from your DHCP server on the external
network. When using a development environment with
[dev-scripts](https://github.com/openshift-metal3/dev-scripts), it uses
the `baremetal` libvirt network unless an override is specified. The IP
can be retrieved with `virsh net-dhcp-leases baremetal`. If the install
is far enough along to have brought up the provisioning network, you may
use the provisioning bootstrap IP which defaults to 172.22.0.2.

Viewing the virtual machine's console with virt-manager may also be
helpful.

You can view the Ironic logs by sshing to the bootstrap VM, and
examining the logs of the `ironic` service, `journalctl -u ironic`. You
may also view the logs of the individual containers:

  - `podman logs ipa-downloader`
  - `podman logs coreos-downloader`
  - `podman logs ironic-api`
  - `podman logs ironic-conductor`
  - `podman logs ironic-inspector`
  - `podman logs ironic-dnsmasq`
  - `podman logs ironic-deploy-ramdisk-logs`
  - `podman logs ironic-inspector-ramdisk-logs`


### Control Plane

Once Ironic is available, the installer will provision the three control
plane hosts. For early failures, it may be useful to look at the console
(using virt-manager if emulating baremetal with vbmc, or through the BMC
like iDRAC) and see if there are any errors reported.

Additionally, if the cluster comes up enough that the bootstrap is destroyed,
but commands like `oc get clusteroperators` shows degraded operators, it
may be useful to examine the logs of the pods within the
`openshift-kni-infra` namespace.

### Ironic

You may want to examine Ironic itself and look at the state of the
hosts. On the bootstrap VM there is a `/opt/metal3/auth/clouds.yaml`
file which may be used with the Ironic `baremetal` client.

To interact with Ironic running on the cluster, it will be necessary
to create a similar `clouds.yaml` using the content from the metal3-ironic
secrets in the openshift-machine-api namespace, and the hostIP of the
controlplane host running the metal3 pod.
