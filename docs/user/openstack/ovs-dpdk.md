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
          hw_cpu_policy: dedicated
          hw_mem_page_size: large
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
      hw_cpu_policy: dedicated
      hw_mem_page_size: 1GB
      hw_numa_nodes: '1'
      hw_vif_multiqueue_enabled: true
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

### Create a performance profile

To get the best performance out of the workers attached to the VFIO interface, the correct settings for hugepages and other relevant settings need to be used.
The node-tuning-operator component is used to configure these settings on the desired worker nodes.

The following is an example of the performance profile that should be applied.
The user should provide the values for kernel args, `HUGEPAGES`, `CPU_ISOLATED` and `CPU_RESERVED` depending on their system:

``` yaml
apiVersion: performance.openshift.io/v2
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

### Label your CNF workers

Once the CNF workers are deployed, you must label them SR-IOV capable:

```sh
oc label node <node-name> feature.node.kubernetes.io/network-sriov.capable="true"
```

These workers will be used to run the CNF workloads.

### Install the SRIOV Network Operator and configure a network device

You must install the SR-IOV Network Operator. To install the Operator, you will need access to an account on your OpenShift cluster that has `cluster-admin` privileges. After you log in to the account, [install the Operator][2].

Then, [configure your SR-IOV network device][3]. Note that only `netFilter` needs to be used from the `nicSelector`, as we'll give the Neutron network ID used for DPDK traffic.

Example of `SriovNetworkNodePolicy` named `dpdk1`:

```
apiVersion: sriovnetwork.openshift.io/v1
kind: SriovNetworkNodePolicy
metadata:
  name: dpdk1
  namespace: openshift-sriov-network-operator
spec:
  deviceType: vfio-pci
  isRdma: false
  nicSelector:
    netFilter: openstack/NetworkID:55a54d05-9ec1-4051-8adb-1b5a7be4f1b6
  nodeSelector:
    feature.node.kubernetes.io/network-sriov.capable: 'true'
  numVfs: 1
  priority: 99
  resourceName: dpdk1
```

Note: If the network device plugged to the network is not from Intel and is from Mellanox, then `deviceType` must be set to `netdevice` and `isRdma` set to `true`. 

The SR-IOV network operator will automatically discover the devices connected on that network for each worker, and make them available for use by the CNF pods later.


## Deploy a testpmd pod

This pod is an example how we can create a container that uses the hugepages, the reserved CPUs and the DPDK port:

```
apiVersion: v1
kind: Pod
metadata:
  name: testpmd-dpdk
  namespace: mynamespace
spec:
  containers:
  - name: testpmd
    command: ["sleep", "99999"]
    image: registry.redhat.io/openshift4/dpdk-base-rhel8:v4.9
    securityContext:
      capabilities:
        add: ["IPC_LOCK","SYS_ADMIN"]
      privileged: true
      runAsUser: 0
    resources:
      requests:
        memory: 1000Mi
        hugepages-1Gi: 1Gi
        cpu: '2'
        openshift.io/dpdk1: 1
      limits:
        hugepages-1Gi: 1Gi
        cpu: '2'
        memory: 1000Mi
        openshift.io/dpdk1: 1
    volumeMounts:
      - mountPath: /dev/hugepages
        name: hugepage
        readOnly: False
  volumes:
  - name: hugepage
    emptyDir:
      medium: HugePages
```

More examples are documented [here][4].

[1]: https://doc.dpdk.org/guides/prog_guide/poll_mode_drv.html
[2]: https://docs.openshift.com/container-platform/4.10/networking/hardware_networks/installing-sriov-operator.html
[3]: https://docs.openshift.com/container-platform/4.10/networking/hardware_networks/configuring-sriov-device.html
[4]: https://docs.openshift.com/container-platform/4.10/networking/hardware_networks/using-dpdk-and-rdma.html#example-vf-use-in-dpdk-mode-intel_using-dpdk-and-rdma
