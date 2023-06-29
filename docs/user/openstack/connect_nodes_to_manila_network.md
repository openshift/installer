# Connecting worker nodes to a dedicated Manila network

Depending on your OpenStack platform configuration, it is very likely that Manila uses a dedicated network for its shares. When that is the case, you need to attach your OpenShift compute nodes to this network otherwise pods might fail to mount the shares, as shown in the following example:

```txt
Mounting command: mount
Mounting arguments: -t nfs 172.16.32.1:/volumes/_nogroup/891cb5d9-a417-43a5-9d1c-5b160a16e7be /var/lib/kubelet/pods/c36e7573-853a-44f9-99bd-0de630edb3b9/volumes/kubernetes.io~csi/pvc-6b632043-2580-4ada-a634-ae696db4b96c/mount
Output: mount.nfs: Connection timed out
```

You will also see this error message on pods, which will be stuck in
`ContainerCreating` state with the following warning:

```
Warning  FailedMount  3m46s (x271 over 17h)  kubelet  MountVolume.SetUp failed for volume "<uuid>" : rpc error: code = DeadlineExceeded desc = context deadline exceeded
```

> **Note**
> Consult with your OpenStack administrator to know what network ID Manila exposes its shares on.

To connect your workers at the time of installation you can use [additionalNetworkIDs](https://github.com/openshift/installer/blob/master/docs/user/openstack/customization.md#additional-networks) parameter in the install config and set Manila network ID there:

Example OpenShift install config:

```yaml
...
compute:
- name: worker
  platform:
    openstack:
      additionalNetworkIDs:
      - <manila_network_id>
...
```

As day 2 operation you need to add new network at `networks` section of your machineset's [provider spec](https://github.com/openshift/installer/blob/master/docs/user/openstack/README.md#defining-a-machineset-that-uses-multiple-networks).
After that Cluster API Provider OpenStack will automatically connect your workers to the network.

Example of OpenStack Machine Spec:

```yaml
networks:
  ...
  - noAllowedAddressPairs: true
    uuid: <manila_network_id>
  ...
```

> **Note**
> The `noAllowedAddressPairs` option ensures we do not configure allowed address
> pairs when creating the port connected to this network. This is often
> prevented by policy and is correct behavior as this network will not be used
> for API access. This is behavior also seen for all networks specified via
> `compute.*.platform.openstack.additionalNetworkIDs` in
> `install-config.yaml`.
