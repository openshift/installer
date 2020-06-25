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
that process in the [Virtual IPs section](#virtual-ips).

## CoreDNS-mDNS

https://github.com/openshift/CoreDNS-mdns/

The `mDNS` plugin for `CoreDNS` was developed to perform DNS lookups
based on discoverable information from mDNS. This plugin will resolve both the
`etcd-NNN` records, as well as the `_etcd-server-ssl._tcp.` SRV record. It is also
able to resolve the name of the nodes.

The list of `etcd` hosts included in the SRV record is based on the list of
control plane nodes currently running.

The IP addresses that the `etcd-NNN` host records resolve to comes from the
mDNS advertisement sent out by the `mDNS-publisher` on that control plane node.

## Virtual IPs

We use virtual IP addresses, VIPs, managed by Keepalived to provide high
availability access to essential APIs and services. For more info on how this
works, please read about what [Keepalived is](https://www.keepalived.org/) and
about the underlying [VRRP
algorithm](https://en.wikipedia.org/wiki/Virtual_Router_Redundancy_Protocol)
that it runs. In our current implementation, we have 2 highly available VIPs
that we manage.  Ingress VIP handles requests to services managed by OpenShift
and the API VIP handles requests to the openshift API. Our VIP addresses are
chosen and validated from the nodes subnet in the openshift installer, however,
the services we run to manage the internal networking infrastructure, such as
Keepalived, dns, and loadbalancers, are manged by the
[machine-config-operator(MCO)](https://github.com/openshift/machine-config-operator/tree/master/docs).
The MCO has been configured to run static pods on the Bootstrap, Master, and
Worker Nodes that run our internal networking infrastructure. Files run on the
bootstrap node can be found
[here](https://github.com/openshift/machine-config-operator/tree/master/manifests/openstack).
Files run on both master and worker nodes can be found
[here](https://github.com/openshift/machine-config-operator/tree/master/templates/common/openstack/files).
Files run on only master nodes can be found
[here](https://github.com/openshift/machine-config-operator/tree/master/templates/master/00-master/openstack/files).
Lastly, files run only on worker nodes can be found
[here](https://github.com/openshift/machine-config-operator/tree/master/templates/worker/00-worker/openstack/files).

## Infrastructure Walkthrough

The bootstrap node is responsible for running temporary networking infrastructure while the Master
nodes are still coming up. The bootstrap node will run a CoreDNS instance, as well as
Keepalived. While the bootstrap node is up, it will have priority running the API VIP.

The Master nodes run dhcp, HAProxy, CoreDNS, mDNS-publisher, and Keepalived. Haproxy loadbalances incoming requests
to the API across all running masters. It also runs a stats and healthcheck server. Keepalived manages both VIPs on the master, where each
master has an equal chance of being assigned one of the VIPs. Initially, the bootstrap node has the highest priority for hosting the API VIP, so they will point to addresses there at startup. Meanwhile, the master nodes will try to get the control plane, and the OpenShift API up. Keepalived implements periodic health checks for each VIP that are used to determine the weight assigned to each server. The server with the highest weight is assigned the VIP. Keepalived has two seperate healthchecks that attempt to reach the OpenShift API and CoreDNS on the localhost of each master node. When the API on a master node is reachable, Keepalived substantially increases it's weight for that VIP, making its priority higher than that of the bootstrap node and any node that does not yet have the that service running. This ensures that nodes that are incapable of serving DNS records or the OpenShift API do not get assigned the respective VIP. The Ingress VIP is also managed by a healthcheck that queries for an OCP Router HAProxy healthcheck, not the HAProxy we stand up in  static pods for the API. This makes sure that the Ingress VIP is pointing to a server that is running the necessary OpenShift Ingress Operator resources to enable external access to the node.

The Worker Nodes run dhcp, CoreDNS, mDNS-publisher, and Keepalived. On workers, Keepalived is only responsible for managing
the Ingress VIP. It's algorithm is the same as the one run on the masters.
