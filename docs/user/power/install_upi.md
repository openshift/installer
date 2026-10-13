# Install: Power User Provided Infrastructure

The steps for performing a UPI-based install are outlined here. Example automation is provided [here](#example-power-upi-configuration) to help model your own.

## Table of contents

1. [Minimum compute requirements](#minimum-resource-requirements)

2. [Network topology requirements](#network-topology-requirements)

3. [DNS requirements](#dns-requirements)

4. [Getting Ignition configs for machines](#getting-ignition-configs-for-machines)

5. [Booting machines with RHCOS and Ignition configs](#booting-machines-with-rhcos-and-ignition-configs)

6. [Watching your installation (bootstrap complete, cluster available)](#watching-your-installation)

7. [Example Bare-Metal UPI deployment](#example-power-upi-configuration)

## Compute

The smallest OpenShift 4.x cluster requires the following hosts:

* 1 bootstrap machine.

* 3 control plane machines.

* 2 worker machines.

NOTE: The cluster requires the bootstrap machine to deploy the OpenShift cluster on to the 3 control plane machines. Once the installation completes you can remove the bootstrap machine.

All the machines must use Red Hat Enterprise Linux CoreOS (RHCOS) as the operating system.

### Minimum resource requirements

| Machine       | Operating System | vCPU | RAM   | Storage |
|---------------|------------------|------|-------|---------|
| Bootstrap     | RHCOS            | 2    | 16 GB | 120 GB  |
| Control Plane | RHCOS            | 2    | 16 GB | 120 GB  |
| Compute       | RHCOS            | 2    | 16 GB | 120 GB  |



## Network Topology Requirements

The easiest way to get started is to ensure all Power nodes have internet access to pull images for platform containers and provide telemetry data to Red Hat.
OpenShift 4.x also supports a restricted network installation.

### Load balancers

Before you install OpenShift, you must provision two load balancers.

* A load balancer for the control plane and bootstrap machines that targets port 6443 (Kubernetes APIServer) and 22623([Machine Config server][machine-config-server]). Port 6443 must be accessible to both clients external to the cluster and nodes within the cluster, and port 22623 must be accessible to nodes within the cluster.

  NOTE: Bootstrap machine can be deleted after cluster installation is finished.

* A load balancer for the machines that run the [ingress router][openshift-router] pods that balances ports 443 and 80. Both the ports must be accessible to both clients external to the cluster and nodes within the cluster.

    NOTE: A working configuration for the ingress router is required for an OpenShift 4.x cluster.

    NOTE: The default configuration for Cluster Ingress Operator  deploys the ingress router to `worker` nodes in the cluster. The administrator needs to configure the [ingress][openshift-router] after the control plane has been bootstrapped.

### Connectivity between machines

You must configure the network connectivity between machines to allow cluster components to communicate.

* etcd

    As the etcd members are located on the control plane machines. Each control plane machine requires connectivity to [etcd server][etcd-ports], [etcd peer][etcd-ports] and [etcd-metrics][etcd-ports] on every other control plane machine.

* OpenShift SDN

    All the machines require connectivity to certain reserved ports on every other machine to establish in-cluster networking. For further detail, please refer to the following [documentation][sdn-ports].

* Kubernetes NodePort

    All the machines require connectivity to Kubernetes NodePort range 30000-32767 on every other machine for OpenShift platform components.

* OpenShift reserved

    All the machines require connectivity to reserved port ranges 10250-12252 and 9000-9999 on every other machine for OpenShift platform components.

## DNS requirements

* Kubernetes API

    OpenShift 4.x requires the DNS records `api.$cluster_name.$base_domain` and `api-int.$cluster_name.$base_domain` to point to the load balancer targeting the control plane machines. Both records must be resolvable from all the nodes within the cluster. The `api.$cluster_name.$base_domain` must also be resolvable by clients external to the cluster.

* etcd

    For each control plane machine, OpenShift 4.x requires DNS records `etcd-$idx.$cluster_name.$base_domain` to point to `$idx`'th control plane machine. The DNS record must resolve to an unicast IPV4 address for the control plane machine and the records must be resolvable from all the nodes in the cluster.

    For each control plane machine, OpenShift 4.x also requires a SRV DNS record for etcd server on that machine with priority `0`, weight `10` and port `2380`. For 3 control plane cluster, the records look like:

    ```plain
    # _service._proto.name.                            TTL   class SRV priority weight port target.
    _etcd-server-ssl._tcp.$cluster_name.$base_domain   86400 IN    SRV 0        10     2380 etcd-0.$cluster_name.$base_domain.
    _etcd-server-ssl._tcp.$cluster_name.$base_domain   86400 IN    SRV 0        10     2380 etcd-1.$cluster_name.$base_domain.
    _etcd-server-ssl._tcp.$cluster_name.$base_domain   86400 IN    SRV 0        10     2380 etcd-2.$cluster_name.$base_domain.
    ```

* OpenShift Routes

    OpenShift 4.x requires the DNS record `*.apps.$cluster_name.$base_domain` to point to the load balancer targeting the machines running the ingress router pods. This record must be resolvable by both clients external to the cluster and from all the nodes within the cluster.

## Getting Ignition configs for machines

The OpenShift Installer provides administrators various assets that are required to create an OpenShift cluster, namely:

* Ignition configs: The OpenShift Installer provides Ignition configs that should be used to configure the RHCOS based bootstrap and control plane machines using `bootstrap.ign`  and `master.ign` respectively. The OpenShift Installer also provides `worker.ign` that can be used to configure the RHCOS based `worker` machines.

* Admin Kubeconfig: The OpenShift Installer provides a kubeconfig with admin level privileges to Kubernetes APIServer.

    NOTE: This kubeconfig is configured to use `api.$cluster_name.$base_domain` DNS name to communicate with the Kubernetes APIServer.

### Setting up install-config for installer

The OpenShift installer uses an [Install Config](../customization.md#platform-customization) to drive all install time configuration.

An example install config for bare-metal UPI is as follows:

```yaml
apiVersion: v1
## The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.
baseDomain: example.com
compute:
- name: worker
  replicas: 1
controlPlane:
  name: master
  replicas: 3
metadata:
  ## The name for the cluster
  name: test
platform:
  none: {}
## The pull secret that provides components in the cluster access to images for OpenShift components.
pullSecret: ''
## The default SSH key that will be programmed for `core` user.
sshKey: ''
```

Create a directory that will be used by the OpenShift installer to provide all the assets. For example `test-bare-metal`,

```console
$ mkdir test-bare-metal
$ tree test-bare-metal
test-bare-metal

0 directories, 0 files
```

Copy *your* `install-config` to the `INSTALL_DIR`. For example using the `test-bare-metal` as our `INSTALL_DIR`,

```console
$ cp <your-instal-config> test-bare-metal/install-config.yaml
$ tree test-bare-metal
test-bare-metal
└── install-config.yaml

0 directories, 1 file
```

NOTE: The filename for `install-config` in the `INSTALL_DIR` must be `install-config.yaml`

### Invoking the installer to get Ignition configs

Given that you have setup the `INSTALL_DIR` with the appropriate `install-config`, you can create the Ignition configs by using the `create ignition-configs` target. For example,

```console
$ openshift-install --dir test-bare-metal create ignition-configs
INFO Consuming "Install Config" from target directory
$ tree test-bare-metal
test-bare-metal
├── auth
│   └── kubeconfig
├── bootstrap.ign
├── master.ign
└── worker.ign

1 directory, 4 files
```

The `bootstrap.ign`, `master.ign`, and `worker.ign` files must be made available as http/https file downloads resolvable by the RHCOS nodes.

## Booting machines with RHCOS and Ignition configs

### Required kernel parameters for boot
A kernel parameter file must be created for each node with the following parameters:

* `rd.neednet=1`: [CoreOS Installer][coreos-installer] needs internet access to fetch the OS image that needs to be installed on the machine.

* IP configuration [arguments](https://docs.openshift.com/container-platform/4.3/installing/installing_bare_metal/installing-bare-metal-network-customizations.html#network-customization-config-yaml_installing-bare-metal-network-customizations) may be required to access the network.

* CoreOS Installer [arguments][coreos-installer-args] are required to be configured to install RHCOS and setup the Ignition config file for that machine.

* Refer to the following docs for details on booting a PowerVM machine
  - [iso boot](https://www.ibm.com/developerworks/community/wikis/home?lang=en#!/wiki/Power+Systems/page/Mounting+an+ISO+image+on+VIO+client+LPAR)
  - [network boot](https://www.ibm.com/developerworks/community/wikis/home?lang=en#!/wiki/Power+Systems/page/How+to+initiate+network+boot+of+an+LPAR)

## Watching your installation

### Monitor for bootstrap-complete

The administrators can use the `wait-for bootstrap-complete` target of the OpenShift Installer to monitor cluster bootstrapping. The command succeeds when it notices `bootstrap-complete` event from Kubernetes APIServer. This event is generated by the bootstrap machine after the Kubernetes APIServer has been bootstrapped on the control plane machines. For example,

```console
$ openshift-install --dir test-bare-metal wait-for bootstrap-complete
INFO Waiting up to 30m0s for the Kubernetes API at https://api.test.example.com:6443...
INFO API v1.16.2 up
INFO Waiting up to 30m0s for bootstrapping to complete...
```

### Configure Image Registry Storage Provisioner


The Cluster Image Registry [Operator][cluster-image-registry-operator] does not pick a storage backend for `None` platform. Therefore, the cluster operator will be stuck in progressing because it is waiting for the administrator to [configure][cluster-image-registry-operator-configuration] a storage backend for the image-registry.
[NFS][openshift-nfs] should be picked as a [storage-backend][nfs-storage-backend].


#### Configuring NFS

To make an existing NFS share accessible for OpenShift to use as persistent storage, users must first attach it as a Persistent Volume.  At least 100GB of NFS storage space must be available for the image registry claim.

```
apiVersion: v1
kind: PersistentVolume
spec:
  accessModes:
  - ReadWriteMany
  - ReadWriteOnce
  capacity:
    storage: 100Gi
  nfs:
    path: <NFS export path>
    server: <ip of NFS server>
  persistentVolumeReclaimPolicy: Recycle
  volumeMode: Filesystem
status: {}
```

Once the persistent volume is created, the image registry must be patched to use it.

```sh
oc patch configs.imageregistry.operator.openshift.io cluster --type merge --patch '{"spec":{"storage":{"pvc":{"claim":""}}, "managementState": "Managed"}}'
```

#### Configuring Local Storage (testing/development only)

Alternatively, for non-production clusters, `emptyDir` can be used for testing instead of NFS.

```sh
oc patch configs.imageregistry.operator.openshift.io cluster --type merge --patch '{"spec":{"storage":{"emptyDir":{}}, "managementState": "Managed"}}'
```


## Monitor for cluster completion

The administrators can use the `wait-for install-complete` target of the OpenShift Installer to monitor cluster completion. The command succeeds when it notices that Cluster Version Operator has completed rolling out the OpenShift cluster from Kubernetes APIServer.

```console
$ openshift-install wait-for install-complete
INFO Waiting up to 30m0s for the cluster to initialize...
```

## Example Power UPI configuration

An [example terraform configuration](https://github.com/ppc64le/ocp4_upi_powervm) for deploying a
self-contained, development/testing cluster on Power is available.  This example
configuration demonstrates a minimal set of infrastructure services to bring
up a running cluster.  It is not a production-ready configuration.

The repository includes examples of the following user-provided components,
which are intended to serve as a guide for designing a user's cluster
topology.

* DNS
* Load Balancing
* DHCP
* File Server (for Ignition configs)

[cluster-image-registry-operator-configuration]: https://github.com/openshift/cluster-image-registry-operator#registry-resource
[cluster-image-registry-operator]: https://github.com/openshift/cluster-image-registry-operator#image-registry-operator
[coreos-installer-args]: https://github.com/coreos/coreos-installer#kernel-command-line-options-for-coreos-installer-running-in-the-initramfs
[coreos-installer]: https://github.com/coreos/coreos-installer#coreos-installer
[csr-requests]: https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/#requesting-a-certificate
[etcd-ports]: https://github.com/openshift/origin/pull/21520
[machine-config-server]: https://github.com/openshift/machine-config-operator/blob/master/docs/MachineConfigServer.md
[openshift-router]: https://github.com/openshift/cluster-ingress-operator#openshift-ingress-operator
[rrdns]: https://tools.ietf.org/html/rfc1794
[sdn-ports]: https://github.com/openshift/origin/pull/21520
[upi-metal-example-pre-req]: ../../../upi/metal/README.md#pre-requisites
[upi-metal-example]: ../../../upi/metal/README.md
[openshift-nfs]: https://docs.openshift.com/container-platform/4.3/storage/persistent_storage/persistent-storage-nfs.html
[nfs-storage-backend]: https://docs.openshift.com/container-platform/4.3/registry/configuring_registry_storage/configuring-registry-storage-baremetal.html
