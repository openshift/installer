# Installing a cluster on OpenStack that supports DPDK-connected compute machines

If the underlying OpenStack deployment has Open vSwitch with Data Plane Development Kit (OVS-DPDK) support enabled,
an OpenShift cluster can take advantage of the feature by providing access to [poll mode drivers][1].

## Pre-requisites

The following steps might be required to be checked before starting the deployment of OpenShift.

### Selecting hugepages size and CPU policy for OpenShift virtual machines

You require hugepages for every OCP cluster virtual machine running on the hosts with OVS-DPDK.
If hugepages are not present in the flavor, the machines will get created but the network interface won't work.
You must check with your OpenStack administrator the size and number of available hugepages which can be used per node.

This may be accomplished in two ways:

- The flavors used for the OCP are defined with appropriate values. For example:

    ``` bash
    openstack flavor create m1.xlarge.nfv --ram 16384 --disk 40 --vcpus 8
    openstack flavor set --property hw:cpu_policy=dedicated --property hw:mem_page_size=large m1.xlarge.nfv
    ```

    This method requires administrative access the OpenStack deployment because flavors need to be created and modified.

- The image properties maybe be defined in the install-config.yaml using [platform.openstack.clusterOSImageProperties](customization.md).
  For example:

    ``` yaml
    platform:
      openstack:
         clusterOSImageProperties:
            - hw_cpu_policy: dedicated
            - hw_mem_page_size: large
    ```

If using image properties, then the flavor used to deploy the virtual machines must allow the image property values selected.
For example, if you attempt to deploy a virtual machine using an image with the property `hw_mem_page_size: 1GB`,
but the flavor being used is set with `--property hw:mem_page_size=small`,
then the deployment of the virtual machine will fail because the flavor does not allow for using a 1GB page size (large).
The flavor property would need to either not be set, which is the same as `--property hw:mem_page_size=any`,
or the flavor property should be set to `--property hw:mem_page_size=large`.

### Create an aggregate group (optional)

If not all your OpenStack compute nodes are configured for OVS-DPDK, you must create an aggregate group and add the compatible nodes into it.

Create the group, add the nodes, set the property and update the flavor:

``` bash
openstack aggregate create dpdk_group
openstack aggregate add host dpdk_group <compute-hosts>
openstack aggregate set --property dpdk=true dpdk_group
openstack flavor set <flavor> --property aggregate_instance_extra_specs:dpdk=true
```

### Enable multiqueue (optional)

If you use multiqueue with OVS-DPDK, set the `hw_vif_multiqueue_enabled` property on the image used for OCP:

```
openstack image set --property hw_vif_multiqueue_enabled=true <coreos_image>
```

where `coreos_image` is the name of the RHCOS image to be used.
This is the value later entered for `platform.openstack.clusterOSImage`.
This can also be accomplished by using `platform.openstack.clusterOSImageProperties` in `install-config.yaml`:

``` yaml
platform:
  openstack:
    clusterOSImageProperties:
      - hw_cpu_policy: dedicated
      - hw_mem_page_size: 1GB
      - hw_numa_nodes: '1'
      - hw_vif_multiqueue_enabled: true
```

### Create additional network to attach poll mode driver interfaces

The `poll_mode_driver` network will be where the ports using poll mode drivers will be attached to.
In most cases, this will be a Neutron provider network.

``` bash
openstack network create dpdk_network --provider-physical-network dpdk-phy1 --provider-network-type vlan --provider-segment <VLAN-ID>
openstack subnet create dpdk_subnet --network dpdk_network --subnet-range 192.0.2.0/24 --dhcp
```

The ID of the network will be later used to inform the cluster which ports to attach the poll mode drivers to.

## OpenShift cluster deployment

### Deploy the cluster

1. Create and edit the `install-config.yaml`.

    ``` bash
    openshift-install create install-config --log-level=debug
    ```

    Edit the `install-config.yaml`, providing the following values:

    - `compute.platform.openstack.additionalNetworkIDs`

        The ID of the Neutron network where the ports will use the poll mode drivers will be attached to.

2. If required, [create the manifests](../customization.md#install-time-customization-for-machine-configuration) to bind the PMD drivers to ports on the desired nodes.

    Using a manifest at this stage only allows the selection of all control-plane nodes and/or all worker nodes.
    It is not possible to select a subgroup of workers or control-plane nodes.

3. Create the cluster.

    ``` bash
    openshift-install create cluster --log-level=debug
    ```

### Install Performance Addon Operator (PAO) and create a performance profile

To get the best performance out of the workers attached to the VFIO interface, the correct settings for hugepages and other relevant settings need to be used.
The performance-operator is used to configure these settings on the desired worker nodes.
For details on installing and using the performance-operator see the [official documentation][2].

The following is an example of the performance profile that should be applied.
The user should provide the values for kernel args, `HUGEPAGES`, `CPU_ISOLATED` and `CPU_RESERVED` depending on their system:

``` yaml
apiVersion: performance.openshift.io/v1
kind: PerformanceProfile
metadata:
  name: cnf-performanceprofile
spec:
  additionalKernelArgs:
    - nmi_watchdog=0
    - audit=0
    - mce=off
    - processor.max_cstate=1
    - idle=poll
    - intel_idle.max_cstate=0
    - default_hugepagesz=1GB
    - hugepagesz=1G
    - intel_iommu=on
  cpu:
    isolated: <CPU_ISOLATED>
    reserved: <CPU_RESERVED>
  hugepages:
    defaultHugepagesSize: 1G
    pages:
      - count: <HUGEPAGES>
        node: 0
        size: 1G
  nodeSelector:
    node-role.kubernetes.io/worker: ''
  realTimeKernel:
    enabled: false
```

### Enable `vfio-noiommu` driver

In the case of deploying workers in virtual machines, the underlying virtualization platform doesn't support
a virtualized IOMMU and we need to explicitly enable it via `options vfio enable_unsafe_noiommu_mode=1`.
This can be done via a MachineConfig resource:

``` yaml
kind: MachineConfig
apiVersion: machineconfiguration.openshift.io/v1
metadata:
  name: 99-vfio-noiommu 
  labels:
    machineconfiguration.openshift.io/role: worker
spec:
  config:
    ignition:
      version: 3.2.0
    storage:
      files:
      - path: /etc/modprobe.d/vfio-noiommu.conf
        mode: 0644
        contents:
          source: data:;base64,b3B0aW9ucyB2ZmlvIGVuYWJsZV91bnNhZmVfbm9pb21tdV9tb2RlPTEK
```

More documentation about that can be found [here][3].

### Bind the `vfio-pci` kernel driver to the appropriate NICs

The workers connected to the `VFIO_NETWORK` require the `vfi-pci` kernel driver to be bound to the ports attached to the `VFIO_NETWORK`.

We basically need to create a MachineSet for the workers which will be attached to the `VFIO_NETWORK` network.
The MachineConfig used in the MachineSet should contain the necessary SystemD unit definitions to install the poll mode driver (`vfio-pci`) on boot.
This [Ansible role][4] provides an easy way to generate the required MachineSet.
The playbooks also provide the ability to automatically install the MachineSet, but it is highly recommended to use the playbooks for generating a MachineSet file (.yaml).
Then the user should explicitly install this on their cluster using the appropriate `oc` commands.

To generate the MachineSet file, run the following commands:

``` bash
openstack network show <name of VFIO network> -f value -c id
ansible-playbook play.yaml -e network_ids="ID of VFIO network" --extra-vars "mc_config_file=/tmp/vhostuser.yaml"
oc create -f /tmp/vhostuser.yaml
```

In the case of multiple networks, a comma-separated list of IDs can be provided to `network_ids`.

The MachineSet will enable a script at boot which will inspects the config drive attached to the VM for the required information. If not found there, it queries the metadata server for it. Based on the network ID, the MAC address of the port is learned and this is translated into the correct PCI bus ID. This means that on boot the vfio-pci module is bound to any port that is attached to the network identified by the OpenStack network ID.

The result is that the interface will have the vfio-pci driver bound to it. You can verify this by executing the `lspci -k` command on the workers in question. The result should be similar to this:

``` bash
lspci -k
...
00:07.0 Ethernet controller: Red Hat, Inc. Virtio network device
	Subsystem: Red Hat, Inc. Device 0001
	Kernel driver in use: vfio-pci
```

As you can see, the Ethernet controller at bus ID `00:07.0` is using the vfio-pci kernel driver.

### Expose host-device interface to pod

The host-device CNI plugin can be used to expose an interface on the host to the pod.
The plugin moves the interface from the host network namespace to the pods namespace.
Effectively given direct control of the interface to the pod running in the namespace.

This can be done by [creating an additional network attachment with the host-device CNI plug-in][5]:
``` yaml
additionalNetworks:
- name: vhost1
  namespace: <your-cnf-namespace>
  type: Raw
  rawCNIConfig: '{
    "cniVersion": "0.3.1",
    "name": "vhost1",
    "type": "host-device",
    "pciBusId": "0000:00:04.0",
    "ipam": {}
    }
  }'
```

After this step, validation can be done by doing `oc -n <your-cnf-namespace> get net-attach-def` to ensure the networks are created on the namespace.

[1]: https://doc.dpdk.org/guides/prog_guide/poll_mode_drv.html
[2]: https://docs.openshift.com/container-platform/4.9/scalability_and_performance/cnf-performance-addon-operator-for-low-latency-nodes.html
[3]: https://github.com/k8snetworkplumbingwg/sriov-network-device-plugin/blob/master/docs/dpdk/README-virt.md
[4]: https://github.com/rh-nfv-int/shift-on-stack-vhostuser
[5]: https://docs.openshift.com/container-platform/4.9/networking/multiple_networks/configuring-additional-network.html#nw-multus-host-device-object_configuring-additional-network
