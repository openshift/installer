# Install: vSphere User Provided Infrastructure

The steps for performing a UPI-based install are outlined here. Several [Terraform][upi-vsphere-example] templates are provided as an example to help model your own.

## Table of contents

1. Compute

2. Network topology requirements

3. DNS requirements

4. Getting Ignition configs for machines

5. Getting OS related assets for machines

6. Configuring RHCOS VMs with Ignition configs

7. Watching your installation (bootstrap_complete, cluster available)

8. Example vSphere UPI deployment

## Compute

The smallest OpenShift 4.x clusters require the following VMs:

* 1 bootstrap machine.

* 3 control plane machines.

* at least 1 worker machine.

NOTE: The cluster requires the bootstrap machine to deploy the OpenShift cluster on to the 3 control plane machines, and you can remove the bootstrap machine.

The bootstrap and control plane machines must use Red Hat Enterprise Linux CoreOS (RHCOS) as the operating system.

All of the VMs created must reside within the same folder in all of the vCenters used. For example, if the VM folder used is named MyFolder, then all of the VMs running OpenShift must be in the MyFolder folder.

The disk UUID on the VMs must be enabled: the `disk.EnableUUID` value must be set to `True`. This step is necessary so that the VMDK always presents a consistent UUID to the VM, thus allowing the disk to be mounted properly.

### Minimum resource requirements

Processing
Memory
Storage
Networking

[todo-link-to-minimum-resource-requirements]

## Network Topology Requirements

OpenShift 4.x requires all nodes to have internet access to pull images for platform containers and provide telemetry data to Red Hat.

### Load balancers

Before you install OpenShift, you must provision two load balancers.

* A load balancer for the control plane machines that targets port 6443 (Kubernetes APIServer) and 22623([Machine Config server][machine-config-server]). Port 6443 must be accessible to both clients external to the cluster and nodes within the cluster, and port 22623 must be accessible to nodes within the cluster.

* A load balancer for the machines that run the [ingress router][openshift-router] pods that balances ports 443 and 80. Both the ports must be accessible to both clients external to the cluster and nodes within the cluster.

    NOTE: A working configuration for the ingress router is required for an OpenShift 4.x cluster.

    NOTE: The default configuration for Cluster Ingress Operator deploys the ingress router to `worker` nodes in the cluster. The administrator needs to configure the [ingress][openshift-router] after the control plane has been bootstrapped.

### Connectivity between machines

You must configure the network connectivity between machines to allow cluster components to communicate.

* etcd

    As the etcd members are located on the control plane machines, each control plane machine requires connectivity to [etcd server][etcd-ports], [etcd peer][etcd-ports] and [etcd-metrics][etcd-ports] on every other control plane machine.

    All the worker machines should have connectivity to [etcd server][etcd-ports] and [etcd-metrics][etcd-ports] ports on the control plane machines.

* OpenShift SDN

    All the machines require connectivity to certain reserved ports on every other machine to establish in-cluster networking. For more details refer [doc][snd-ports].

* Kubernetes NodePort

    All the machines require connectivity to Kubernetes NodePort range 30000-32767 on every other machine for OpenShift platform components.

* OpenShift reserved

    All the machines require connectivity to reserved port ranges 10250-12252 and 9000-9999 on every other machine for OpenShift platform components.

### Connectivity during machine boot

All the RHCOS machines require network in `initramfs` during boot to fetch Ignition config from the Machine Config Server [machine-config-server]. During the initial boot, the machines requires a DHCP server in order to establish a network connection to download their Ignition configs. After the initial boot, the machines can be configured to use a static IP address.

## DNS requirements

* Kubernetes API

    OpenShift 4.x requires the DNS record `api.$cluster_name.$base_domain` to point to the Load balancer targeting the control plane machines. This record must be resolvable by both clients external to the cluster and from all the nodes within the cluster.

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

    OpenShift 4.x requires the DNS record `*.apps.$cluster_name.$base_domain` to point to the Load balancer targeting the machines running the ingress router pods. This record must be resolvable by both clients external to the cluster and from all the nodes within the cluster.

## Getting ignition configs for machines

The OpenShift Installer provides administrators various assets that are required to create an OpenShift cluster, namely:

* Ignition configs: The OpenShift Installer provides Ignition configs that should be used to configure the RHCOS based bootstrap and control plane machines using `bootstrap.ign`  and `master.ign` respectively. The OpenShift Installer also provides `worker.ign` that can be used to configure the RHCOS based `worker` machines, but also can be used as source for configuring RHEL based machines [todo-link-to-BYO-RHEL].

* Admin Kubeconfig: The OpenShift Installer provides a kubeconfig with admin level privileges to Kubernetes APIServer.

    NOTE: This kubeconfig is configured to use `api.$cluster_name.$base_domain` DNS name to communicate with the Kubernetes APIServer.

### Setting up install-config for installer

The OpenShift installer uses an [Install Config][install-config] to drive all install time configuration.

An example install config for vSphere UPI is as follows:

```yaml
apiVersion: v1beta4
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
  vsphere:
    ## The hostname or IP address of the vCenter
    vcenter: your.vcenter.server
    ## The name of the user for accessing the vCenter
    username: your_vsphere_username
    ## The password associated with the user
    password: your_vsphere_password
    ## The datacenter in the vCenter
    datacenter: your_datacenter
    ## The default datastore to use.
    defaultDatastore: your_datastore
## The pull secret that provides components in the cluster access to images for OpenShift components.
pullSecret: ''
## The default SSH key that will be programmed for `core` user.
sshKey: ''
```

Create a directory that will be used by the OpenShift installer to provide all the assets. For example `test-vsphere`,

```console
$ mkdir test-vsphere
$ tree test-vsphere
test-vsphere

0 directories, 0 files
```

Copy *your* `install-config` to the `INSTALL_DIR`. For example using the `test-vsphere` as our `INSTALL_DIR`,

```console
$ cp <your-instal-config> test-vsphere/install-config.yaml
$ tree test-vsphere
test-vsphere
└── install-config.yaml

0 directories, 1 file
```

NOTE: The filename for `install-config` in the `INSTALL_DIR` must be `install-config.yaml`

### Invoking the installer to get Ignition configs

Given that you have setup the `INSTALL_DIR` with the appropriate `install-config`, you can create the Ignition configs by using the `create ignition-configs` target. For example,

```console
$ openshift-install --dir test-vsphere create ignition-configs
INFO Consuming "Install Config" from target directory
$ tree test-vsphere
test-vsphere
├── auth
│   └── kubeconfig
├── bootstrap.ign
├── master.ign
├── metadata.json
└── worker.ign

1 directory, 5 files

```

## Getting OS related assets for machines

TODO RHEL CoreOS does not have assets for vSphere.

## Configuring RHCOS VMs with Ignition configs

Set the vApp properties of the VM to set the Ignition config for the VM. The `guestinfo.ignition.config.data` property is the base64-encoded Ignition config. The `guestinfo.ignition.config.data.encoding` should be set to `base64`.

### Control Plane and Worker VMs

The Ignition configs supplied in the vApp properties of the control plane and worker VMs should be copies of the `master.ign` and `worker.ign` created by the OpenShift Installer.

### Bootstrap VM

The Ignition config supplied in the vApp properties of the bootstrap VM should be an Ignition config that has a URL from which the bootstrap VM can download the `bootstrap.ign` created by the OpenShift Installer. Note that the URL must be accessible by the bootstrap VM.

The Ignition config created by the OpenShift Installer cannot be used directly because there is a size limit on the length of vApp properties, and the Ignition config will exceed that size limit.

```json
{
  "ignition": {
    "config": {
      "append": [
        {
          "source": "bootstrap_ignition_config_url",
          "verification": {}
        }
      ]
    },
    "timeouts": {},
    "version": "2.1.0"
  },
  "networkd": {},
  "passwd": {},
  "storage": {},
  "systemd": {}
}
```
### Hostname

TODO: Link to document on modifying the Ignition configs to set the hostnames.

### Static IP Addresses

TODO: Link to document on modifying the Ignition configs to set static IP addresses.

## Watching your installation

### Monitor for bootstrap-complete

The administrators can use the `wait-for bootstrap-complete` command of the OpenShift Installer to monitor cluster bootstrapping. The command succeeds when it notices `bootstrap-complete` event from Kubernetes APIServer. This event is generated by the bootstrap machine after the Kubernetes APIServer has been bootstrapped on the control plane machines. For example,

```console
$ openshift-install --dir test-vsphere wait-for bootstrap-complete
INFO Waiting up to 30m0s for the Kubernetes API at https://api.test.example.com:6443...
INFO API v1.12.4+c53f462 up
INFO Waiting up to 30m0s for the bootstrap-complete event...
```

### Monitor for install completion

The administrators can use the `wait-for install-complete` command of the OpenShift Installer to monitor install completion. The command succeeds when it notices that Cluster Version Operator has completed rolling out the OpenShift cluster from Kubernetes APIServer.

```console
$ openshift-install --dir test-vsphere wait-for install-complete
INFO Waiting up to 30m0s for the cluster to initialize...
```

## Example vSphere UPI deployment

Terraform [templates][upi-vsphere] provides an example of using OpenShift Installer to create a vSphere UPI OpenShift cluster.

### Overview

* Compute:
    Uses `public` IPv4 addresses for each machine, so that all the machines are accessible on the Internet.

* DNS and Load Balancing
    Uses AWS [Route53](aws-route53) to configure the all the DNS records.
    Uses Round-Robin DNS [RRDNS][rrdns] in place of load balancing solutions.

Refer to the pre-requisites for using the example [here][upi-vsphere-example-pre-req]

### Creating the cluster

#### Installer assets

Use the OpenShift Installer to create [Ignition configs][#getting-ignition-configs-for-machines] that will be used to create bootstrap, control plane and worker machines.

#### Terraform variable file

Use the [example][upi-vsphere-example-tfvar] `create_tfvars.sh` script to create a Terraform variable file, and edit the `tfvars` file on your favorite editor.

```sh
cd test-vsphere
create_tfvars.sh
```

At a minimum, you will need to provide values for the following variables.
* bootstrap_ignition_url
* bootstrap_ip
* control_plane_ips
* compute_ips

Move the `tfvars` file to the directory where the example terraform is.

#### Creating resources

Initialize terraform to download all the required providers. For more info on terraform [init][terraform-init] and terraform [providers][terraform-providers]

```sh
terraform init
```

Create all the resources using terraform by invoking [apply][terraform-apply]

```sh
terraform apply -auto-approve
```

#### Monitoring bootstrap-complete and removing bootstrap resources

Use the bootstrap [monitoring][#monitor-for-bootstrap-complete] to track when cluster bootstrapping has finished. After the Kubernetes APIServer has been bootstrapped on the control plane machines, the bootstrap VM can be destroyed by following:

```sh
terraform apply -auto-approve -var 'boostrap_complete=true'
```

### Approving server certificates for nodes

To allow Kube APIServer to communicate with the kubelet running on nodes for logs, rsh etc. The administrator needs to approve the CSR [requests][csr-requests] generated by each kubelet.

You can approve all `Pending` CSR requests using,

```sh
oc get csr -ojson | jq -r '.items[] | select(.status == {} ) | .metadata.name' | xargs oc adm certificate approve
```

### Updating image-registry to emptyDir storage backend

The Cluster Image Registry [Operator][cluster-image-registry-operator] does not pick an storage backend for `vSphere` platform. Therefore, the cluster operator will be stuck in progressing because it is waiting for administrator to [configure][cluster-image-registry-operator-configuration] a storage backend for the image-registry. You can pick `emptyDir` for non-production clusters by following:

```sh
oc patch configs.imageregistry.operator.openshift.io cluster --type merge --patch '{"spec":{"storage":{"filesystem":{"volumeSource": {"emptyDir":{}}}}}}'
```

#### Monitoring cluster completion

Use the cluster finish [monitoring][#monitor-for-cluster-completion] to track when cluster has completely finished deploying.

#### Destroying the cluster

Use terraform [destroy][terraform-destroy] to destroy all the resources for the cluster. For example,

```console
terraform destroy -auto-approve
```

[aws-route53]: https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/Welcome.html
[csr-request]: https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/#requesting-a-certificate
[etcd-ports]: https://github.com/openshift/origin/pull/21520
[install-config]: https://godoc.org/github.com/openshift/installer/pkg/types#InstallConfig
[machine-config-server]: https://github.com/openshift/machine-config-operator/blob/master/docs/MachineConfigServer.md
[openshift-router]: https://github.com/openshift/cluster-ingress-operator#openshift-ingress-operator
[rrdns]: https://tools.ietf.org/html/rfc1794
[sdn-ports]: https://github.com/openshift/origin/pull/21520
[terraform-apply]: https://www.terraform.io/docs/commands/apply.html
[terraform-destroy]: https://www.terraform.io/docs/commands/destroy.html
[terraform-init]: https://www.terraform.io/docs/commands/init.html
[terraform-providers]: https://www.terraform.io/docs/providers/
[upi-vsphere-example-pre-req]: ../../../upi/vsphere/README.md#pre-requisites
[upi-vsphere-example-tfvar]: ../../../upi/vsphere/terraform.tfvar.example
[upi-vsphere-example]: ../../../upi/vsphere/README.md
