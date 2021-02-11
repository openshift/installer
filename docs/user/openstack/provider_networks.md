# Provider Networks


## Table of Contents
- [Provider Networks](#provider-networks)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Common prerequisites](#common-prerequisites)
  - [Considerations when creating provider networks](#considerations-when-creating-provider-networks)
  - [Deploying cluster with primary interface on a provider network with IPI](#deploying-cluster-with-primary-interface-on-a-provider-network-with-ipi)
  - [Known issues](#known-issues)


## Introduction

Provider networks map directly to an existing physical network in a data center.
Example of network types include flat (untagged), VLAN (802.1Q tagged) and VXLAN. 
OpenShift clusters that are installed on provider networks do not require tenant networks or floating IP addresses (FIPs).
Therefore, the OpenShift installer does not create these resources during installation.
More information can be found about provider networks terminology [here][1].

Here is a basic architecture of one OCP cluster running on a provider network and another one
on a tenant network:

![OCP on a provider network](provider-network.png)


## Prerequisites

* The [Neutron service][2] is enabled and accessible through the [OpenStack Networking API][3].
* The Neutron service is configured with the [port-security and allowed-address-pairs][4] extensions so the installer can
  add the `allowed_address_pairs` attribute to ports.


## Considerations when creating provider networks

* The provider network has to be shared with other tenants, otherwise Nova won't be able to request ports on that external
  network. For more details, see [BZ#1933047][5].

        openstack network create --share (...)

  To secure that network, it is advised to create [RBAC][6] rules so the network can be only usable by a specific project.


* The project that you use to install OpenShift must own the provider network.

    The provider network and the subnet must be owned by the project that is used to install OpenShift instead of `admin`.
    If they are not, you will have to run the installer from the admin user to create ports on the network.

    It is important that the provider network and the subnet are owned by the same project that will be used
    to install OpenShift (from the clouds.yaml) and we don't want them to be owned by `admin` otherwise
    it'll cause Terraform to fail creating the ports.

    Example commands to create a network and subnet for a project that is named `openshift`:

        openstack network create --project openshift (...)
        openstack subnet create --project openshift (...)

    More information can be found about how to create provider networks [here][7].

* You'll have to make sure that the provider network can reach
  the Metadata IP (169.254.169.254) which, depending on the OpenStack SDN and how Neutron
  is configured (e.g. DHCP servers provide metadata network routes) might involve
  to provide the route when creating the subnet:

    openstack subnet create --dhcp --host-route destination=169.254.169.254/32,gateway=$ROUTER_IP" (...)

**Note:** We're working on removing the nova-metadata requirement but for now it is strongly required to be
          enabled in the cloud and reachable from the provider network.


## Deploying cluster with primary interface on a provider network with IPI


- Considerations: make sure all prerequisites documented previously have been met.

- Create install-config.yaml:

    - Set `platform.openstack.apiVIP` to the IP address for the API VIP.
    - Set `platform.openstack.ingressVIP` to the IP address for the Ingress VIP.
    - Set `platform.openstack.machinesSubnet` to the subnet ID of the provider network subnet.
    - Set `networking.machineNetwork.cidr` to the CIDR of the provider network subnet.

**Note:**

`platform.openstack.apiVIP` and `platform.openstack.ingressVIP` both need to be an unassigned IP
address on the `networking.machineNetwork.cidr`.

    Example:

        (...)
        platform:
          openstack:
            apiVIP: <IP address in the provider network reserved for the API VIP>
            ingressVIP: <IP address in the provider network reserved for the Ingress VIP>
            machinesSubnet: <provider network subnet ID>
            (...)
        networking:
          machineNetwork:
          - cidr: <provider network subnet CIDR>

- Run the OpenShift installer:

      ./openshift-install create cluster --log-level debug

- Wait for the installer to complete.


[1]: <https://docs.openstack.org/neutron/latest/admin/archives/adv-features.html#provider-networks>
[2]: <https://docs.openstack.org/neutron>
[3]: <https://docs.openstack.org/api-ref/network>
[4]: <https://docs.openstack.org/api-ref/network/v2/#allowed-address-pairs>
[5]: <https://bugzilla.redhat.com/show_bug.cgi?id=1933047>
[6]: <https://docs.openstack.org/neutron/latest/admin/config-rbac.html>
[7]: <https://access.redhat.com/documentation/en-us/red_hat_openstack_platform/16.1/html/networking_guide/sec-networking-concepts#provider-networks>
