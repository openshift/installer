# Bare-Metal: Installation requirements

The Tectonic Installer creates bare-metal Tectonic clusters within networks with PXE infrastructure and the `matchbox` service.

For more information about `matchbox`, reference the [`matchbox` documentation][matchbox].

## Network

Bare-metal Tectonic clusters are provisioned in a PXE network environment. Cluster nodes will PXE boot from the `matchbox` service running on a provisioner node. Familiarity with your network topology is required.

Tectonic bare metal clusters store credentials in `user-data`, and etcd peer to peer communication is not currently encrypted with TLS. To restrict access to sensitive information, provision bare metal machines within a trusted network and ensure that a firewall exists between cluster controllers and the public internet.

### Services

Ensure DHCP, TFTP and DNS services are available on your network. CoreOS provides a [dnsmasq][matchbox-dnsmasq] container, if you wish to use rkt or Docker for this.

### PXE

Familiarize yourself with PXE booting. Cluster nodes should PXE boot from the network and delegate to the `matchbox` service which serves configs to provision clusters. At a high level, you will need to:

* Chainload PXE firmwares to iPXE
* Point iPXE client machines to the `matchbox` iPXE HTTP endpoint (e.g. `http://matchbox.example.com:8080/boot.ipxe`)

### DNS

The installer will prompt for "Controller" and "Tectonic" DNS names. For the controller DNS name, add a record which resolves to the node you plan to use as a controller.

By default, Tectonic Ingress runs as a [Kubernetes Daemon Set][daemonset] across workers. For the Tectonic DNS name, add a record which resolves to any node(s) you plan to use as workers.

* Add a DNS name which resolves to the provisioner (e.g. `matchbox.example.com`)
* Add a DNS name which resolves to any controller node (e.g. `k8s.example.com`)
* Add a DNS name which resolves to any worker nodes (e.g. `tectonic.example.com`)

### Machines

* Know the MAC address and stable DNS name for each cluster node
* Configure cluster nodes to favor booting from disk. Be able to use IPMI to request a PXE boot.
* Add a DNS name (and static IP) so that cluster nodes have stable names which can be used during cluster configuration (e.g. `node3.example.com`)

### Egress whitelist

Cluster nodes will need to be able to pull docker images from [quay.io][quay.io] and gcr.io. Be sure to whitelist these domains.

## Machines

A minimum of 3 machines are required to run Tectonic.

### Cluster nodes

Tectonic clusters consist of two types of nodes:

* Controller Nodes - Controller nodes run `etcd` and the control plane of the cluster.
* Worker Nodes - Worker nodes run your applications. New worker nodes will join the cluster by talking to controller nodes for admission.

Each node should meet the following tech-specs.

| Requirement | Value                        |
|-------------|------------------------------|
| RAM         | 8GB / node                   |
| CPU         | 2 cores / node               |
| Storage     | 30GB / node                  |

#### Boot from disk

Configure cluster nodes to favor booting from disk, and use IPMI to request a [PXE boot during installation and re-provisioning][reprovision]. Booting from disk allows Container Linux automatic updates to function normally and is the recommended configuration after provisioning.

Sites where cluster nodes always boot from PXE must plan to regularly update the Container Linux image served to clients.

Client machines:

* PXE boot from the matchbox service
* Favor boot from disk before PXE
* Pull images from Quay.io
* Expose ports 8080, 443, 2379, and 2380 for etcd and Kubernetes services

### Provisioner Node

A provisioner node (or Kubernetes cluster) runs the `matchbox` network boot and provisioning service, along with PXE services if you don't already run them elsewhere. You may use CoreOS or any Linux distribution for this node. It serves provisioning configs to nodes, but does not join Tectonic clusters.

The provisioner must:

* Run a new systemd service or pod for `matchbox`
* Expose a port of the read-only API and the read-write API
* Add /var/lib/matchbox and /etc/matchbox directories
* Add a matchbox user/group or use an existing non-root account
* Be resolvable at a DNS name (e.g. `matchbox.example.com`)
* Generate TLS server credentials (along with client credentials)
* Serve CoreOS PXE and install images

## Tectonic Installer

The Tectonic Installer app runs on a user's laptop as a GUI for creating new clusters and pushing the right configs to `matchbox`.

User machines:

* Run a Linux or Darwin binary Installer app
* Can resolve the matchbox read-write API (e.g. `matchbox.example.com`)
* Use matchbox TLS client credentials (generated via docs)
* Can resolve cluster nodes to provisioning progress (e.g. `node3.example.com`)
* Can SSH to one of the controller nodes (currently required to finish bootstrapping)
* Have a Tectonic software license and Docker pull secret from `tectonic.com`


[daemonset]: https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/
[reprovision]: uninstall.md
[matchbox-dnsmasq]: https://github.com/coreos/matchbox/tree/master/contrib/dnsmasq
[matchbox]: https://coreos.com/matchbox
[quay.io]: https://quay.io
