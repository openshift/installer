# Installing With SR-IOV Worker Nodes

## Table of Contents

- [Prerequisites](#prerequisites)
- [Preparing your OpenShift Cluster to use SR-IOV Networks](#preparing-your-openshift-cluster-to-use-sr-iov-networks)
- [Creating SR-IOV Networks for Worker Nodes](#creating-sr-iov-networks-for-worker-nodes)
- [Creating SR-IOV Worker Nodes in IPI](#creating-sr-iov-worker-nodes-in-ipi)
- [Creating SR-IOV Worker Nodes in UPI](#creating-sr-iov-worker-nodes-in-upi)

## Prerequisites

Single Root I/O Virtualization (SR-IOV) networking in OpenShift can benefit applications
that require high bandwidth and low latency. To plan an OpenStack deployment that uses SR-IOV network interface cards (NICs), refer to [the OSP 16.1 installation documentation](osp-sriov-install). you install an OpenShift cluster on OpenStack, make sure that the NICs that your OpenStack nodes use [are supported](supported-nics) for use with SR-IOV
in OpenShift, and that your tenant has access to them. Your OpenStack cluster must meet the following quota requirements for each OpenShift node that has an attached SR-IOV NIC:

- One instance from the RHOSP quota
- One port attached to the machines subnet
- One port for each SR-IOV Virtual Function
- A flavor with at least 16 GB memory, 4 vCPUs, and 25 GB storage space

For all clusters that use single-root input/output virtualization (SR-IOV), RHOSP compute nodes require a flavor that supports [huge pages](huge-pages).
Deploying worker nodes with SR-IOV networks is supported as a post-install operation for both IPI and UPI workflows. After you verify that your OpenStack cluster can support SR-IOV in OpenShift and you install an OpenShift cluster that meets the [minimum requirements](README.md#openstack-requirements), use the following steps and examples to create worker nodes with SR-IOV NICs.

## Preparing your OpenShift Cluster to use SR-IOV Networks

Before you use single root I/O virtualization (SR-IOV) on a cluster that runs on Red Hat OpenStack Platform (RHOSP), make the RHOSP metadata service mountable as a drive and enable the No-IOMMU Operator for the virtual function I/O (VFIO) driver.

### Enabling the RHOSP metadata service as a mountable drive

Create a machine config in your worker machine pool that makes the OpenStack metadata service available as a mountable drive. This machine config enables the SR-IOV operator to get the UUID of networks from OpenStack.

```yaml
kind: MachineConfig
apiVersion: machineconfiguration.openshift.io/v1
metadata:
  name: 20-mount-config
  labels:
    machineconfiguration.openshift.io/role: worker
spec:
  osImageURL: ''
  config:
    ignition:
      version: 2.2.0
    systemd:
      units:
        - name: create-mountpoint-var-config.service
          enabled: true
          contents: |
            [Unit]
            Description=Create mountpoint /var/config
            Before=kubelet.service

            [Service]
            ExecStart=/bin/mkdir -p /var/config

            [Install]
            WantedBy=var-config.mount

        - name: var-config.mount
          enabled: true
          contents: |
            [Unit]
            Before=local-fs.target
            [Mount]
            Where=/var/config
            What=/dev/disk/by-label/config-2
            [Install]
            WantedBy=local-fs.target
```

### Enabling the No-IOMMU feature for the RHOSP VFIO driver

Apply a machine config to your worker machine pool that enables the No-IOMMU feature for the Red Hat OpenStack Platform (RHOSP) virtual function I/O (VFIO) driver.

```yaml
kind: MachineConfig
apiVersion: machineconfiguration.openshift.io/v1
metadata:
  name: 99-vfio-noiommu
  labels:
    machineconfiguration.openshift.io/role: worker
spec:
  osImageURL: ''
  config:
    ignition:
      version: 2.2.0
    storage:
      files:
      - filesystem: root
        path: "/etc/modprobe.d/vfio-noiommu.conf"
        contents:
          source: data:text/plain;charset=utf-8;base64,b3B0aW9ucyB2ZmlvIGVuYWJsZV91bnNhZmVfbm9pb21tdV9tb2RlPTEK
          verification: {}
        mode: 0644
```

If you need to configure your deployment for real-time or low latency workloads, install the [performance addon operator](performance-addon-operator). If services on a node need to use the performance addon operator or DPDK, that node needs [additional configuration to support hugepages](huge-pages-perf-addon).

After your OpenShift control plane is running, you must install the SR-IOV Network Operator. To install the Operator, you will need access to an account on your OpenShift cluster that has `cluster-admin` privileges. After you log in to the account, [install the Operator](sriov-operator). Then, [configure your SR-IOV network device](configure-sriov-network-device).

## Creating SR-IOV Networks for Worker Nodes

You must create SR-IOV networks to attach to worker nodes before you create the nodes. Reference the following example of how to create radio and uplink provider networks in OpenStack:

```sh
# Create Networks
openstack network create radio --provider-physical-network radio --provider-network-type vlan --provider-segment 120
openstack network create uplink --provider-physical-network uplink --provider-network-type vlan --external

# Create Subnets
openstack subnet create --network radio --subnet-range <radio_network_subnet_range> radio
openstack subnet create --network uplink --subnet-range <uplink_network_subnet_range> uplink
```

## Creating SR-IOV Worker Nodes in IPI

You can create worker nodes as a post-IPI-install operation by using the machine API. To create a new set of worker nodes, [create a new machineSet in OpenShift](openstack-machine-sets).

```sh
oc get machineset -n openshift-machine-api <machineset_name> -o yaml > sriov_machineset.yaml
```

When editing an existing machineSet (or a copy of one) to create SR-IOV worker nodes, add each subnet that is configured for SR-IOV to the machineSet's `providerSpec`. The following example attaches ports from the `radio` and `uplink` subnets, which were created in the previous example, to all of the worker nodes in the machineSet. For all SR-IOV ports, you must set the following parameters:

 - `nicType: direct`
 - `portSecurity:false`

The SR-IOV Operator requires a config drive to be mounted to each instance using SR-IOV VFs, so always set `configDrive: true`. Note that security groups or allowedAddressPairs can not be set on a port if `portSecurity` is disabled. If you are using a network with port security disabled, then allowed address pairs and security groups cannot be used for any port in that network. Setting security groups on the instance will apply that security group to all ports attached to it, be aware of this when using networks with port security disabled. Right now, trunking is not enabled on ports defined in the `ports` list, only the ports created by entries in the `networks` or `subnets` lists. The name of the port will be `<machine-name>-<nameSuffix>`, and the `nameSuffix` is required field in the port definition. Optionally, you can add tags to ports by adding them to the `tags` list. The following example shows how a machineset can be created that creates SR-IOV capable ports on the `Radio` and `Uplink` networks and subnets that were defined in a previous example:


```yaml
apiVersion: machine.openshift.io/v1beta1
kind: MachineSet
metadata:
  labels:
    machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
    machine.openshift.io/cluster-api-machine-role: <node_role>
    machine.openshift.io/cluster-api-machine-type: <node_role>
  name: <infrastructure_ID>-<node_role>
  namespace: openshift-machine-api
spec:
  replicas: <number_of_replicas>
  selector:
    matchLabels:
      machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
      machine.openshift.io/cluster-api-machineset: <infrastructure_ID>-<node_role>
  template:
    metadata:
      labels:
        machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
        machine.openshift.io/cluster-api-machine-role: <node_role>
        machine.openshift.io/cluster-api-machine-type: <node_role>
        machine.openshift.io/cluster-api-machineset: <infrastructure_ID>-<node_role>
    spec:
      metadata:
      providerSpec:
        value:
          apiVersion: openstackproviderconfig.openshift.io/v1alpha1
          cloudName: openstack
          cloudsSecret:
            name: openstack-cloud-credentials
            namespace: openshift-machine-api
          flavor: <nova_flavor>
          image: <glance_image_name_or_location>
          serverGroupID: <optional_UUID_of_server_group>
          kind: OpenstackProviderSpec
          networks:
            - subnets:
              - uuid: <machines_subnet_uuid>
          ports:
            - networkID: <radio_network_uuid>
              nameSuffix: radio
              fixedIPs:
                - subnetID: <radio_subnet_uuid>
              tags:
                - sriov
                - radio
              vnicType: direct
              portSecurity: false
            - networkID: <uplink_network_uuid>
              nameSuffix: uplink
              fixedIPs:
                - subnetID: <uplink_subnet_uuid>
              tags:
                - sriov
                - uplink
              vnicType: direct
              portSecurity: false
          primarySubnet: <machines_subnet_uuid>
          securityGroups:
          - filter: {}
            name: <infrastructure_ID>-<node_role>
          serverMetadata:
            Name: <infrastructure_ID>-<node_role>
            openshiftClusterID: <infrastructure_ID>
          tags:
          - openshiftClusterID=<infrastructure_ID>
          trunk: true
          userDataSecret:
            name: <node_role>-user-data
          availabilityZone: <optional_openstack_availability_zone>
          configDrive: true
```

After you finish editing your machineSet, upload it to your OpenShift cluster:

```sh
oc create -f sriov_machineset.yaml
```

To create SR-IOV ports on a network with the port security disabled, you need to make additional changes to your machineSet due to security groups being set on the instance by default, and allowed address pairs automatically getting added to ports created through the `networks` and `subnets` interfaces. The solution is to define all of your ports with the `ports` interface in your machineSet. Remember that the port for the machines subnet needs:
- allowed address pairs for your API and ingress vip ports
- the worker security group
- to be attached to the machines network and subnet

```yaml
apiVersion: machine.openshift.io/v1beta1
kind: MachineSet
metadata:
  labels:
    machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
    machine.openshift.io/cluster-api-machine-role: <node_role>
    machine.openshift.io/cluster-api-machine-type: <node_role>
  name: <infrastructure_ID>-<node_role>
  namespace: openshift-machine-api
spec:
  replicas: <number_of_replicas>
  selector:
    matchLabels:
      machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
      machine.openshift.io/cluster-api-machineset: <infrastructure_ID>-<node_role>
  template:
    metadata:
      labels:
        machine.openshift.io/cluster-api-cluster: <infrastructure_ID>
        machine.openshift.io/cluster-api-machine-role: <node_role>
        machine.openshift.io/cluster-api-machine-type: <node_role>
        machine.openshift.io/cluster-api-machineset: <infrastructure_ID>-<node_role>
    spec:
      metadata: {}
      providerSpec:
        value:
          apiVersion: openstackproviderconfig.openshift.io/v1alpha1
          cloudName: openstack
          cloudsSecret:
            name: openstack-cloud-credentials
            namespace: openshift-machine-api
          flavor: <nova_flavor>
          image: <glance_image_name_or_location>
          kind: OpenstackProviderSpec
          configDrive: True
          ports:
            - allowedAddressPairs:
              - ipAddress: <api_vip_port_IP>
              - ipAddress: <ingress_vip_port_IP>
              fixedIPs:
                - subnetID: <machines_subnet_UUID>
              nameSuffix: nodes
              networkID: <machines_network_UUID>
              securityGroups:
                  - <worker_security_group_UUID>
            - networkID: <sriov_network_UUID>
              nameSuffix: sriov
              fixedIPs:
                - subnetID: <sriov_subnet_UUID>
              tags:
                - sriov
              vnicType: direct
              portSecurity: False
          primarySubnet: <machines_subnet_UUID>
          serverMetadata:
            Name: <infrastructure_ID>-<node_role>
            openshiftClusterID: <infrastructure_ID>
          tags:
          - openshiftClusterID=<infrastructure_ID>
          trunk: false
          userDataSecret:
            name: worker-user-data
```

## Creating SR-IOV Worker Nodes in UPI

Because UPI implementation depends largely on your deployment environment and requirements, there is no official script for deploying SR-IOV worker nodes. However, we can share a verified example that is based on the [compute-nodes.yaml](../../../upi/openstack/compute-nodes.yaml) script to help you understand the process. To use the script, open up a terminal to the location of the `inventory.yaml` and `common.yaml` UPI Ansible scripts. In the following example, we add provider networks named `radio` and `uplink` to the `inventory.yaml` file. Note that the count parameter specifies the number of virtual functions (VFs) to attach to each worker node. This code can also be found on [github](https://github.com/shiftstack/SRIOV-Compute-Nodes-Ansible-Automation).

```yaml
....
# If this value is non-empty, the corresponding floating IP will be
# attached to the bootstrap machine. This is needed for collecting logs
# in case of install failure.
    os_bootstrap_fip: '203.0.113.20'

    additionalNetworks:
    - id: radio
    count: 4
    type: direct
    port_security_enabled: no
    - id: uplink
    count: 4
    type: direct
    port_security_enabled: no
```
Next, create a file called `compute-nodes.yaml` with this Ansible script:

```yaml
- import_playbook: common.yaml

- hosts: all
  gather_facts: no

  vars:
    worker_list: []
    port_name_list: []
    nic_list: []

  tasks:
  # Create the SDN/primary port for each worker node
  - name: 'Create the Compute ports'
    os_port:
      name: "{{ item.1 }}-{{ item.0 }}"
      network: "{{ os_network }}"
      security_groups:
      - "{{ os_sg_worker }}"
      allowed_address_pairs:
      - ip_address: "{{ os_ingressVIP }}"
    with_indexed_items: "{{ [os_port_worker] * os_compute_nodes_number }}"
    register: ports

  # Tag each SDN/primary port with cluster name
  - name: 'Set Compute ports tag'
    command:
      cmd: "openstack port set --tag {{ cluster_id_tag }} {{ item.1 }}-{{ item.0 }}"
    with_indexed_items: "{{ [os_port_worker] * os_compute_nodes_number }}"

  - name: 'List the Compute Trunks'
    command:
      cmd: "openstack network trunk list"
    when: os_networking_type == "Kuryr"
    register: compute_trunks

  - name: 'Create the Compute trunks'
    command:
      cmd: "openstack network trunk create --parent-port {{ item.1.id }} {{ os_compute_trunk_name }}-{{ item.0 }}"
    with_indexed_items: "{{ ports.results }}"
    when:
    - os_networking_type == "Kuryr"
    - "os_compute_trunk_name|string not in compute_trunks.stdout"

  - name: ‘Call additional-port processing’
    include_tasks: additional-ports.yaml

  # Create additional ports in OpenStack
  - name: ‘Create additionalNetworks ports’
    os_port:
      name:  "{{ item.0 }}-{{ item.1.name }}"
      vnic_type: "{{ item.1.type }}"
      network: "{{ item.1.uuid }}"
      port_security_enabled: "{{ item.1.port_security_enabled|default(omit) }}"
      no_security_groups: "{{ 'true' if item.1.security_groups is not defined else omit }}"
      security_groups: "{{ item.1.security_groups | default(omit) }}"
    with_nested:
      - "{{ worker_list }}"
      - "{{ port_name_list }}"

  # Tag the ports with the cluster info
  - name: 'Set additionalNetworks ports tag'
    command:
      cmd: "openstack port set --tag {{ cluster_id_tag }} {{ item.0 }}-{{ item.1.name }}"
    with_nested:
      - "{{ worker_list }}"
      - "{{ port_name_list }}"

  # Build the nic list to use for server create
  - name: Build nic list
    set_fact:
      nic_list: "{{ nic_list | default([]) + [ item.name ] }}"
    with_items: "{{ port_name_list }}"

  # Create the servers
  - name: 'Create the Compute servers'
    vars:
      worker_nics: "{{ [ item.1 ] | product(nic_list) | map('join','-') | map('regex_replace', '(.*)', 'port-name=\\1') | list }}"
    os_server:
      name: "{{ item.1 }}"
      image: "{{ os_image_rhcos }}"
      flavor: "{{ os_flavor_worker }}"
      auto_ip: no
      userdata: "{{ lookup('file', 'worker.ign') | string }}"
      security_groups: []
      nics:  "{{ [ 'port-name=' + os_port_worker + '-' + item.0|string ] + worker_nics }}"
      config_drive: yes
    with_indexed_items: "{{ worker_list }}"
```

Create a new Ansible script named `additional-ports.yaml`:

```yaml
Build a list of worker nodes with indexes
- name: ‘Build worker list’
  set_fact:
    worker_list: "{{ worker_list | default([]) + [ item.1 + '-' + item.0 | string ] }}"
  with_indexed_items: "{{ [ os_compute_server_name ] * os_compute_nodes_number }}"

# Ensure that each network specified in additionalNetworks exists
- name: ‘Verify additionalNetworks’
  os_networks_info:
    name: "{{ item.id }}"
  with_items: "{{ additionalNetworks }}"
  register: network_info

# Expand additionalNetworks by the count parameter in each network definition
- name: ‘Build port and port index list for additionalNetworks’
  set_fact:
    port_list: "{{ port_list | default([]) + [ {
                    'net_name' : item.1.id,
                    'uuid' : network_info.results[item.0].openstack_networks[0].id,
                    'type' : item.1.type|default('normal'),
                    'security_groups' : item.1.security_groups|default(omit),
                    'port_security_enabled' : item.1.port_security_enabled|default(omit)
                    } ] * item.1.count|default(1) }}"
    index_list: "{{ index_list | default([]) + range(item.1.count|default(1)) | list }}"
  with_indexed_items: "{{ additionalNetworks }}"

# Calculate and save the name of the port
# The format of the name is cluster_name-worker-workerID-networkUUID(partial)-count
# i.e. fdp-nz995-worker-1-99bcd111-1
- name: ‘Calculate port name’
  set_fact:
    port_name_list: "{{ port_name_list | default([]) + [ item.1 | combine( {'name' : item.1.uuid | regex_search('([^-]+)') + '-' + index_list[item.0]|string } ) ] }}"
  with_indexed_items: "{{ port_list }}"
  when: port_list is defined
```

Finally, run the `compute-nodes.yaml` script as you normally would:

```sh
ansible-playbook -i inventory.yaml compute-nodes.yaml
```

Make sure to follow the documentation to [approve the CSRs](approve-csr-upi) for your worker nodes, and to [wait for the installation to complete](wait-for-install-complete) to finalize your deployment.

[wait-for-install-complete]: install_upi.md#wait-for-the-openshift-installation-to-complete
[approve-csr-upi]: install_upi.md#approve-the-worker-csrs
[machine-pool-customizations]: customization.md#machine-pools
[sriov-operator]: https://docs.openshift.com/container-platform/4.7/networking/hardware_networks/installing-sriov-operator.html
[configure-sriov-network-device]: https://docs.openshift.com/container-platform/4.7/virt/virtual_machines/vm_networking/virt-configuring-sriov-device-for-vms.html
[supported-nics]: https://docs.openshift.com/container-platform/4.7/networking/hardware_networks/about-sriov.html#supported-devices_about-sriov
[osp-sriov-install]: https://access.redhat.com/documentation/en-us/red_hat_openstack_platform/16.1/html-single/network_functions_virtualization_planning_and_configuration_guide/index#assembly_sriov_parameters
[openstack-machine-sets]: https://docs.openshift.com/container-platform/4.7/machine_management/creating_machinesets/creating-machineset-osp.html
[performance-addon-operator]: https://docs.openshift.com/container-platform/4.7/scalability_and_performance/cnf-performance-addon-operator-for-low-latency-nodes.html#about_hyperthreading_for_low_latency_and_real_time_applications_cnf-master
[huge-pages]: https://access.redhat.com/documentation/en-us/red_hat_openstack_platform/16.1/html/configuring_the_compute_service_for_instance_creation/configuring-compute-nodes-for-performance#configuring-huge-pages-on-compute-nodes-osp
[huge-pages-perf-addon]: https://docs.openshift.com/container-platform/4.7/scalability_and_performance/what-huge-pages-do-and-how-they-are-consumed-by-apps.html
