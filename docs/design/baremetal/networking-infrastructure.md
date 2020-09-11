# Bare Metal IPI Networking Infrastructure

The `baremetal` platform (IPI for Bare Metal hosts) automates a number
of networking infrastructure requirements that are handled on other
platforms by cloud infrastructure services.

For an overview of the expected network environment that an administrator must
prepare for a `baremetal` platform cluster, see the [install
documentation](../../user/metal/install_ipi.md).

## Load-balanced control plane access

Access to the Kubernetes API (port 6443) from clients both external
and internal to the cluster should be load-balanced across control
plane machines.

Access to Ignition configs (port 22623) from clients within the
cluster should also be load-balanced across control plane machines.

In both cases, the installation process expects these ports to be
reachable on the bootstrap VM at first and then later on the
newly-deployed control plane machines.

On other platforms (for example, see [the AWS UPI
instructions](../../user/aws/install_upi.md)) an external
load-balancer is required to be configured in advance in order to
provide this access.

### API VIP (Virtual IP)

In the `baremetal` platform, a VIP (Virtual IP) is used to provide
failover of the API server across the control plane machines
(including the bootstrap VM). This "API VIP" is provided by the user
as an `install-config.yaml` parameter and the installation process
configures `keepalived` to manage this VIP.

The API VIP first resides on the bootstrap VM. The `keepalived`
instance here is managed by systemd and a script is used to generate
the `keepalived` configuration before launching the service using
`podman`. See [here](../../../data/data/bootstrap/baremetal/README.md)
for more information about the relevant bootstrap assets.

The VIP will move to one of the control plane nodes, but only after the
bootstrap process has completed and the bootstrap VM is stopped. This happens
because the `keepalived` instances on control plane machines are configured (in
`keepalived.conf`) with a lower
[VRRP](https://en.wikipedia.org/wiki/Virtual_Router_Redundancy_Protocol)
priority. This ensures that the API on the control plane nodes is fully
functional before the API VIP moves.

These `keepalived` instances are run as [static
pods](https://kubernetes.io/docs/tasks/administer-cluster/static-pod/) and the
relevant assets are [rendered by the Machine Config
Operator](https://github.com/openshift/machine-config-operator/pull/795). See
[here](FIXME: link to a README in MCO) for more information about these assets.

### API load balancing

Once the API VIP has moved to one of the control plane nodes, traffic sent from
external clients to this VIP first hits an `haproxy` load balancer running on
that control plane node.
This instance of `haproxy` will load balance the API traffic across all
of the control plane nodes.

The configuration of `haproxy` will be done by MCO once the following PR is
merged:

https://github.com/openshift/machine-config-operator/pull/795

See [here](FIXME: link to a README in MCO) for more detailed information about
the `haproxy` configuration.

## Internal DNS

Externally resolvable DNS records are required for:

* `api.$cluster_name.$base-domain` -
* `*.apps.$cluster_name.$base_domain` -

In addition, internally resolvable DNS records are required for:

* `api-int.$cluster_name.$base-domain` -

On other platforms (for example, see the CloudFormation templates
referenced by [the AWS UPI
instructions](../../user/aws/install_upi.md)), all of these records
are automatically created using a cloud platform's DNS service.

In the `baremetal` platform, the goal is is to automate as much of the
DNS requirements internal to the cluster as possible, leaving only a
small amount of public DNS configuration to be implemented by the user
before starting the installation process.

In a `baremetal` environment, we do not know the IP addresses of all hosts in
advance.  Those will come from an organization’s DHCP server.  Further, we can
not rely on being able to program an organization’s DNS infrastructure in all
cases.  We address these challenges in two ways:

1. Self host some DNS infrastructure to provide DNS resolution for records only
   needed internal to the cluster.
2. Make use of mDNS (Multicast DNS) to dynamically discover the addresses of
   hosts that we must resolve records for.

### api-int hostname resolution

The CoreDNS server performing our internal DNS resolution includes
configuration to resolve the `api-int` hostname. `api-int` will be resolved to
the API VIP.

### mdns-publisher

https://github.com/openshift/mdns-publisher

The `mdns-publisher` is the component that runs on each host to make itself
discoverable by other hosts in the cluster.  Control plane hosts currently
advertise `master-NN` names, and the worker nodes advertise
`worker-NN` names.  *Note: The `master-NN` and `worker-NN` names were based on records
previously published on AWS clusters, but may no longer be necessary.*

On masters and workers it is run by the `machine-config-operator`.

`mdns-publisher` is not run on the bootstrap node, as there is no need for any
other host to discover the IP address that the bootstrap VM gets from DHCP.

### coredns-mdns

https://github.com/openshift/coredns-mdns/

The `mdns` plugin for `coredns` was developed to perform DNS lookups
based on discoverable information from mDNS. The plugin will resolve the
`master-NN` and `worker-NN` records in the cluster domain.

The IP addresses that the `master-NN` host records resolve to comes from the
mDNS advertisement sent out by the `mdns-publisher` on that control plane node.

### DNS Resolution

Because the baremetal platform does not have a cloud DNS service available to
provide internal DNS records, it instead uses a coredns static pod configured
with the `coredns-mdns` plugin discussed above. There is one of these pods
running on every node in a deployment, and a NetworkManager dispatcher script
is used to configure resolv.conf to point at the node's public IP address.
`localhost` can't be used because `resolv.conf` is propagated into some
containers where that won't resolve to the actual host.

### Bootstrap Asset Details

See [here](../../../data/data/bootstrap/baremetal/README.md)
for more information about the relevant bootstrap assets.

## Ingress High Availability

There is a third VIP used by the `baremetal` platform, and that is for Ingress.
The Ingress VIP will always reside on a node running an Ingress controller.
This ensures that we provide high availability for ingress by default.

The mechanism used to determine which nodes are running an ingress controller
is that `keepalived` will try to reach the local `haproxy` stats port number
using `curl`.  This makes assumptions about the default ingress controller
behavior and may be improved in the future.

The Ingress VIP is managed by `keepalived`.  The `keepalived` configuration for
the Ingress VIP will be managed by MCO once the following PR is complete:

https://github.com/openshift/machine-config-operator/pull/795
