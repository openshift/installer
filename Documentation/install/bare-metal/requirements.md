# Bare Metal Installation requirements

The Tectonic Installer creates bare metal Tectonic clusters within networks with PXE infrastructure and the `matchbox` service.

Installation requires the following items, which are discussed in more detail below:

* [Tectonic Account][account]. Register for a Tectonic Account, which is free for up to 10 nodes. The cluster license and pull secret are required during installation of Tectonic.
* Terraform. Tectonic Installer includes and requires a specific version of Terraform. This is included in the Tectonic Installer tarball. See the [Tectonic Installer release notes][release-notes] for information about which Terraform versions are compatible.
* [Matchbox v0.6+][matchbox-latest] installation with TLS client credentials and the gRPC API enabled. For more information, see the Matchbox [installation guide][matchbox-install].
* PXE network boot environment with DHCP, TFTP, and DNS services. For more information, see [Network setup][network-setup].
* DNS records for the Kubernetes controller(s) and Tectonic Ingress worker(s). See [DNS][dns].
* Machines with BIOS options set to boot from disk normally, but PXE prior to installation. `ipmitool` or `virt-install` will be used to actually boot the machines.
* Machines with known MAC addresses and stable domain names.
* A SSH keypair whose private key is present in your system's [ssh-agent][ssh-agent].

## Tectonic Installer

The Tectonic Installer app runs on a user's laptop as a GUI for creating new clusters and pushing the right configs to `matchbox`.

User machines must:

* Run a Linux or Darwin binary Installer app
* Resolve the matchbox read-write API (e.g. `matchbox.example.com`)
* Use matchbox TLS client credentials (generated via docs)
* Resolve cluster nodes to provisioning progress (e.g. `node3.example.com`)
* Have a Tectonic software license and Docker pull secret from `tectonic.com`
* Be able to SSH to one of the controller nodes (currently required to finish bootstrapping)
* Have an SSH keypair whose private key is present in the system's ssh-agent

## Network

Bare metal Tectonic clusters are provisioned in a PXE network environment. Cluster nodes will PXE boot from the `matchbox` service running on a provisioner node. Familiarity with your network topology is required.

Tectonic bare metal clusters store credentials in `user-data`. To restrict access to sensitive information, provision bare metal machines within a trusted network and ensure that a firewall exists between cluster controllers and the public internet.

### Services

Ensure DHCP, TFTP and DNS services are available on your network. CoreOS provides a [dnsmasq][matchbox-dnsmasq] container, if you wish to use rkt or Docker for this.

### PXE

Familiarize yourself with PXE booting. Cluster nodes should PXE boot from the network and delegate to the `matchbox` service which serves configs to provision clusters. At a high level, you must:

* Chainload PXE firmwares to iPXE.
* Point iPXE client machines to the `matchbox` iPXE HTTP endpoint (e.g. `http://matchbox.example.com:8080/boot.ipxe`).

### DNS

For best results, assign DNS names to each node. The following three records are required for Tectonic Installer:

* A DNS name which resolves to the provisioner (e.g. `matchbox.example.com`).
* A DNS name which resolves to any controller node (e.g. `k8s.example.com`).
* A DNS name which resolves to any worker nodes (e.g. `tectonic.example.com`).

### Egress whitelist

Cluster nodes must be able to pull docker images from [quay.io][quay.io] and gcr.io. Be sure to whitelist these domains. If you must whitelist by IP, run `dig quay.io` to list associated IP addresses.

### ssh-agent

Tectonic installer will add the installer machine's public SSH key to all machines in the cluster. The key must be on the installing machine's [ssh-agent][ssh-agent], and it is used to configure nodes.

Check if a key already exists in the ssh-agent using `ssh-add -l`. If a key must be added to the agent, use `ssh-add Path/ToYour/KeyFile`. Note that on OSX it may be necessary to re-add keys from your keyring to the agent on each login.

## Machines

A minimum of 3 machines are required to run Tectonic. To configure machines:

* Know the MAC address and stable DNS name for each cluster node.
* Configure cluster nodes to favor booting from disk. Be able to use IPMI to request a PXE boot.
* Add a DNS name (and static IP) so that cluster nodes have stable names which can be used during cluster configuration (e.g. `node3.example.com`)

### Cluster nodes

Tectonic clusters consist of two types of nodes:

* Controller Nodes: Controller nodes run `etcd` and the control plane of the cluster.
* Worker Nodes: Worker nodes run your applications. New worker nodes will join the cluster by talking to controller nodes for admission.

Each node should meet the following technical specs:

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

### Provisioner node

A provisioner node (or Kubernetes cluster) runs the `matchbox` network boot and provisioning service, along with PXE services if you don't already run them elsewhere. You may use CoreOS Container Linux or any Linux distribution for this node. It serves provisioning configs to nodes, but does not join Tectonic clusters.

The provisioner must:

* Run a new systemd service or pod for `matchbox`
* Expose a port of the read-only API and the read-write API
* Add /var/lib/matchbox and /etc/matchbox directories
* Add a matchbox user/group or use an existing non-root account
* Be resolvable at a DNS name (e.g. `matchbox.example.com`)
* Generate TLS server credentials (along with client credentials)
* Serve CoreOS PXE and install images


[account]: https://account.coreos.com
[daemonset]: https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/
[reprovision]: uninstall.md
[matchbox]: https://coreos.com/matchbox
[matchbox-dnsmasq]: https://github.com/coreos/matchbox/tree/master/contrib/dnsmasq
[matchbox-install]: https://coreos.com/matchbox/docs/latest/deployment.html
[matchbox-latest]: https://github.com/coreos/matchbox/releases
[quay.io]: https://quay.io
[ssh-agent]: https://www.freebsd.org/cgi/man.cgi?query=ssh-agent&sektion=1
[network-setup]: https://coreos.com/matchbox/docs/latest/network-setup.html
[dns]: index.md#dns
