# Connecting worker nodes and pods to an IPv6 network

To connect your workers at the time of installation you can use [additionalNetworkIDs](https://github.com/openshift/installer/blob/master/docs/user/openstack/customization.md#additional-networks) parameter in the install config and set IPv6 network ID there:

Example OpenShift install config:

```yaml
...
compute:
- name: worker
  platform:
    openstack:
      additionalNetworkIDs:
      - <ipv6_network_id>
...
```

> **Note**
> To use Stateful IPv6 Networks, the arg `ip=dhcp,dhcp6` needs to be included in the kernel args of the Worker nodes,
> otherwise the Nodes won't get an IPv6 address due to a [bug][dhcpv6-bug].
> Use the [procedure][add-kernel-args] to add kernel argument to the Nodes.

## Enable connectivity to the pods

To enable connectivity between pods with additional Network on different Nodes, the Port security needs to be disabled for the IPv6 Port of the Server. This way it's possible to avoid adding an allowed-address-pairs with an IP and MAC address in the Server's IPv6 Port whenever a new pod gets created.

```sh
openstack port set --no-security-group --disable-port-security <worker-ipv6-port>
```

## Add IPv6 connectivity to pods

Create a file named `network.yaml` and specify [the desired CNI config][configuring-an-additional-network]. Here is an example of CNI config used for slaac address mode with macvlan:

```yaml
spec:
  additionalNetworks:
  - name: ipv6
    namespace: ipv6
    rawCNIConfig: '{ "cniVersion": "0.3.1", "name": "ipv6", "type": "macvlan", "master": "ens4"}'
    type: Raw
```

The node's interface specified in the Network attachment `master` field may differ from `ens4` when more additional networks are configured or when a different Kernel driver is used.

> **Note**
> When using Stateful address mode, specify the `ipam` section in the CNI config,
> otherwise no address is configured in the additional interface of the pod.
> Also, note that DHCPv6 is not yet supported by [Multus][dhcpv6-multus].

then run:

```sh
oc patch network.operator cluster --patch "$(cat network.yaml)" --type=merge
```

It takes a while for the network definition to be enforced.
You can check with the following command:

```sh
oc get network-attachment-definitions -A
```

## Create pods with IPv6 network

To create pods with IPv6 network make sure to create them on the same Namespace specified in the `additionalNetworks`
and specify the following annotation `k8s.v1.cni.cncf.io/networks: <additional-network-name>`.

[dhcpv6-bug]: https://issues.redhat.com/browse/OCPBUGS-2104
[add-kernel-args]: https://docs.openshift.com/container-platform/4.11/nodes/nodes/nodes-nodes-managing.html#nodes-nodes-kernel-arguments_nodes-nodes-managing
[configuring-an-additional-network]: https://docs.openshift.com/container-platform/4.11/networking/multiple_networks/configuring-additional-network.html#configuring-additional-network_approaches-managing-additional-network
[dhcpv6-multus]: https://issues.redhat.com/browse/SDN-2844
