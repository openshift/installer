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

Once the API VIP has moved to one of the control plane nodes, traffic sent to
this VIP first hits an `haproxy` load balancer running on that control plane
node.  This instance of `haproxy` will load balance the API traffic across all
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
* `etcd-NNN.$cluster_name.$base-domain` -
* SRV record: `_etcd-server-ssl._tcp.$cluster_name.$base-domain` - Lists the
  `etcd-NNN` records.

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
advertise both `etcd-NN` and `master-NN` names, and the worker nodes advertise
`worker-NN` names.  *Note: The `master-NN` and `worker-NN` names were based on records
previously published on AWS clusters, but may no longer be necessary.*

It will be run by the `machine-config-operator`, once the following PR is
complete:

https://github.com/openshift/machine-config-operator/pull/795

`mdns-publisher` is not run on the bootstrap node, as there is no need for any
other host to discover the IP address that the bootstrap VM gets from DHCP.

### coredns-mdns

https://github.com/openshift/coredns-mdns/

The `mdns` plugin for `coredns` was developed to perform DNS lookups
based on discoverable information from mDNS.  This plugin will resolve both the
`etcd-NNN` records, as well as the `_etcd-server-ssl._tcp.` SRV record.

The list of `etcd` hosts included in the SRV record is based on the list of
control plane nodes currently running.

The IP addresses that the `etcd-NNN` host records resolve to comes from the
mDNS advertisement sent out by the `mdns-publisher` on that control plane node.

### etcd clustering with DNS SRV

This section is not `baremetal` platform specific, but it’s important to
understand as the DNS automation put in place is to support this process.

`etcd` includes support for using DNS as a discovery mechanism.  The `etcd`
documentation for this is here: https://github.com/openshift/etcd/blob/master/Documentation/op-guide/clustering.md#dns-discovery

OpenShift makes use of this DNS based discovery mechanism to form the `etcd`
cluster.  You can see the following if you look at one of the `etcd` member
pods in a running cluster:

```
$ oc get pod etcd-member-master-0 -n openshift-etcd -o yaml | grep -A 1 "exec etcd"
      exec etcd \
        --discovery-srv ostest.test.metalkube.org \
```

This static pod definition is generated by the `machine-config-operator` from
this [template](https://github.com/openshift/machine-config-operator/blob/master/templates/master/00-master/_base/files/etc-kubernetes-manifests-etcd-member.yaml).

### DNS Resolution During Bootstrapping

One of the of virtual IP addresses required by the `baremetal` platform is used
for our self-hosted internal DNS - the “DNS VIP”.  The location of the DNS VIP
is managed by `keepalived`, similar to the management of the API VIP.

The control plane nodes are configured to use the DNS VIP as their primary DNS
server.  The DNS VIP resides on the bootstrap host until the control plane
nodes come up, and then it will move to one of the control plane nodes.

`coredns-mdns` on the bootstrap node is configured a little bit differently
than `coredns-mdns` on the control plane nodes.  On the bootstrap node, it is
configured with a minimum number of hosts for the `etcd` SRV record based on
the number of control plane nodes requested in `install-config.yaml`.

This configuration is important.  It ensures that whenever `etcd` resolves this
DNS SRV record to form the cluster, all of the expected control plane hosts are
present in the results.

### DNS Resolution Post-Install

After the bootstrap process completes, the DNS VIP moves to one of the control
plane nodes.  All control plane nodes will continue to use this DNS VIP as
their primary DNS server.

If a control plane node is replaced, the contents of the DNS SRV record will
automatically reflect the new node once it comes up and advertises itself via
DNS.

The current [etcd cluster recovery process](TODO-link?) does not make use of
the DNS SRV record and instead requires manual intervention, but it could
perhaps be useful in the future.

### Bootstrap Asset Details

See [here](../../../data/data/bootstrap/baremetal/README.md)
for more information about the relevant bootstrap assets.
