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

In the `baremetal` platform, a VIP (Virtual IP) is used to provide
failover of the API server across the control plane machines
(including the bootstrap VM). This "API VIP" is provided by the user
as an `install-config.yaml` parameter and the installation process
configures `keepalived` to manage this VIP.

The API VIP first resides on the bootstrap VM. The `keepalived`
instance here is managed by systemd and a script is used to generate
the `keepalived` configuration before launching the service using
`podman`. See [here](../../../data/data/bootstrap/baremetal/README.md)
for more informations about the relevant bootstrap assets.

Once the control plane machines come up, the VIP will move to the one
of these machines. This happens because the `keepalived` instances on
control plane machines are configured (in `keepalived.conf`) with a
higher
[VRRP](https://en.wikipedia.org/wiki/Virtual_Router_Redundancy_Protocol)
priority. These `keepalived` instances are run as [static
pods](https://kubernetes.io/docs/tasks/administer-cluster/static-pod/)
and the relevant assets are [rendered by the Machine Config
Operator](https://github.com/openshift/machine-config-operator/pull/795). See
[here](FIXME: link to a README in MCO) for more information about
these assets.
