# OpenStack IPI Networking Infrastructure

The `OpenStack` platform installer uses an internal networking solution that
is based heavily on the [baremetal networking infrastructure](../baremetal/networking-infrastructure.md).
For an overview of the quotas required, and the entrypoints created when
you build an OpenStack IPI cluster, see the [user docs](../../user/openstack/README.md).


## Load-balanced control plane access

Access to the Kubernetes API (port 6443) from clients both external
and internal to the cluster, and access to ignition configs (port 22623) from clients within the
cluster is load-balanced across control plane machines.
These services are initially hosted by the bootstrap node until the control
plane is up. Then, control is pivoted to the control plane machines. We will go into further detail on
that process in the [Virtual IP's section](#virtual-ips).

## coredns-mdns

https://github.com/openshift/coredns-mdns/

The `mdns` plugin for `coredns` was developed to perform DNS lookups
based on discoverable information from mDNS.  This plugin will resolve both the
`etcd-NNN` records, as well as the `_etcd-server-ssl._tcp.` SRV record.

The list of `etcd` hosts included in the SRV record is based on the list of
control plane nodes currently running.

The IP addresses that the `etcd-NNN` host records resolve to comes from the
mDNS advertisement sent out by the `mdns-publisher` on that control plane node.

## Virtual IP's

We use virtual IP addresses, VIPs, managed by keepalived
to provide high availability access to essential APIs and services. For more info
on how this works, please read about what [keepalived is](https://www.keepalived.org/) and
about the underlying [VRRP algorithm](https://en.wikipedia.org/wiki/Virtual_Router_Redundancy_Protocol)
that it runs. In our current implementation, we have 3 highly available VIPs that we manage: 
Ingress, which handles requests to services managed by OpenShift, DNS, which handles internal dns requests, and API, which handles requests to the openshift API. Our VIP addresses are chosen and validated from the nodes subnet in the openshift
installer, however our VIPs along with our internal dns and loadbalancing is managed
by the [machine-config-operator(MCO)](https://github.com/openshift/machine-config-operator/tree/master/docs).
The MCO has been configured to run static pods on the Bootstrap, Master, and Worker Nodes that
run our internal networking infrastructure. Files run on the bootstrap node can be found 
[here](https://github.com/openshift/machine-config-operator/tree/master/manifests/openstack).
Files run on both master and worker nodes can be found 
[here](https://github.com/openshift/machine-config-operator/tree/master/templates/common/openstack/files).
Files run on only master nodes can be found 
[here](https://github.com/openshift/machine-config-operator/tree/master/templates/master/00-master/openstack/files).
Lastly, files run only on worker nodes can be found 
[here](https://github.com/openshift/machine-config-operator/tree/master/templates/worker/00-worker/openstack/files).

## Infrastructure Walkthrough

The bootstrap node is responsible for running temporary networking infrastructure while the Master
nodes are still coming up. The bootstrap node will run a coredns instance, as well as 
keepalived. While the bootstrap node is up, it will have priority running the API and DNS
VIPs.

The Master nodes run dhcp, haproxy, coredns, mdns publisher, and keepalived. Haproxy loadbalances incoming requests 
to the api across all running masters. It also runs a stats and healthcheck server. Keepalived manages all 3 VIPs on the master, where each
master has an equal chance of being assigned one of the VIPs. Be aware, the bootstrap node has the highest priority for hosting the DNS
and API VIPs, so they will point to addresses there for as long as they can. Once the master nodes have stood up a fully functioning control plane, 
the bootstrap node will be torn down by the installer. Keepalived implements periodic health checks to make sure
that essential services can be reached by their VIPs, when it determines that they cannot, it alters the VIP to point to a server it deems fit.  
To ensure the api is reachable through the API VIP, keepalived periodically attempts to reach the api through the API VIP. It will do the same 
for the DNS VIP. This means that when the bootstrap node gets torn down, keepalived will detect it, and move the API and DNS VIPs to point to a working master node. The Ingress VIP is also managed by a set of health checks that qurey the haproxy healthcheck servers to make sure that the Ingress VIP is pointing to
a server that is able to service it.

The Worker Nodes run dhcp, coredns, mdns publisher, and keepalived. On workers, keepalived is only responsible for managing
the Ingress VIP. Just like on masters, it will periodically check the haproxy healthcheck server to make sure the node the Ingress
VIP points to is up. 