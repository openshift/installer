# Scaling Tectonic bare-metal clusters

This document describes how to add cluster nodes to Tectonic clusters on bare-metal.

## Scaling worker nodes

To add an additional worker node to a existing Tectonic cluster, add configuration for the new machine to the Matchbox provisioner machine, then boot the new node with PXE. The easiest way to add the new machine's configuration to the provisioner is to base it on the configuration of an existing worker node.

First, find two files in the `/var/lib/matchbox/groups/` directory on the provisioner machine: `tectonic-install-*.json` and `tectonic-node-*.json`. The files' names will contain the MAC address of an existing node. Make a copy of each, replacing the MAC in the file name with that of the machine being added. Then, edit the copies to match the configuration of the new worker node.

### Example: Adding a second node to a single node cluster

In this example, a new configuration for a machine with MAC
`00-11-22-33-44-56` is being created based on the configuration of an existing node with MAC `00-11-22-33-44-55`.

```
$ ls /var/lib/matchbox/groups/
tectonic-install-00-11-22-33-44-55.json
tectonic-node-00-11-22-33-44-55.json
```

```
$ cd /var/lib/matchbox/groups/
$ cp tectonic-install-00-11-22-33-44-55.json tectonic-install-00-11-22-33-44-56.json
$ cp tectonic-node-00-11-22-33-44-55.json tectonic-node-00-11-22-33-44-56.json
```

In the new file named `tectonic-install-*.json, change the `"mac"` and `"id"` values to match the MAC of the new node:

```
$ cat /var/lib/matchbox/groups/tectonic-install-00-11-22-33-44-56.json
{
    "id": "tectonic-install-00-11-22-33-44-56",
    "name": "CoreOS Install",
    "profile": "install-reboot",
    "selector": {
        "mac": "00:11:22:33:44:56"
    },
[...]
```

In the file named `tectonic-node-*.json`, edit the `"mac"` and `"domain_name"` values to match the new node's MAC and DNS hostname.

```
$ cat /var/lib/matchbox/groups/tectonic-node-00-11-22-33-44-56.json
{
    "id": "tectonic-node-00-11-22-33-44-56",
    "profile": "tectonic-worker",
    "selector": {
        "mac": "00:11:22:33:44:56",
        "os": "installed"
    },
    "metadata": {
        "domain_name": "node2.example.com",
        "etcd_initial_cluster": "node0=http://node2.example.com:2380",
[...]
```

After these two files are in place and edited with the new node's identifiers, PXE boot the new node to add it to the cluster. See the [Matchbox documentation][matchbox-docs] for details about the network boot provisioning environment.

## Scaling controller nodes

Adding controller nodes to an existing Tectonic cluster on bare metal is not officially supported. This feature is planned for a future release.

### etcd scaling on bare metal

Each controller in a default bare metal Tectonic cluster runs the etcd service. Alternatively, Tectonic Installer will accept the client endpoint of an existing external etcd v3 cluster to which it should connect instead.


[matchbox-docs]: https://coreos.com/matchbox/docs/latest
