# Install: BareMetal User Provided Infrastructure

The steps for performing a UPI-based install are outlined here. Several [Terraform][upi-metal-example] templates are provided as an example to help model your own.

## Table of contents

1. Minimum compute requirements

2. Network topology requirements

3. DNS requirements

4. Getting Ignition configs for machines

5. Getting OS related assets for machines

6. Booting machines with RHCOS and Ignition configs

7. Watching your installation (bootstrap_complete, cluster available)

8. Example Bare-Metal UPI deployment

## Compute

The smallest OpenShift 4.x clusters require the following host:

* 1 bootstrap machine.

* 3 control plane machines.

* at least 1 worker machine.

NOTE: The cluster requires the bootstrap machine to deploy the OpenShift cluster on to the 3 control plane machines, and you can remove the bootstrap machine.

The bootstrap and control plane machines must use Red Hat Enterprise Linux CoreOS (RHCOS) as the operating system.

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

* A load balancer for the control plane and bootstrap machines that targets port 6443 (Kubernetes APIServer) and 22623([Machine Config server][machine-config-server]). Port 6443 must be accessible to both clients external to the cluster and nodes within the cluster, and port 22623 must be accessible to nodes within the cluster.

  NOTE: Bootstrap machine can be deleted as target after cluster installation is finished.

* A load balancer for the machines that run the [ingress router][openshift-router] pods that balances ports 443 and 80. Both the ports must be accessible to both clients external to the cluster and nodes within the cluster.

    NOTE: A working configuration for the ingress router is required for an OpenShift 4.x cluster.

    NOTE: The default configuration for Cluster Ingress Operator  deploys the ingress router to `worker` nodes in the cluster. The administrator needs to configure the [ingress][openshift-router] after the control plane has been bootstrapped.

### Connectivity between machines

You must configure the network connectivity between machines to allow cluster components to communicate.

* Etcd

    As the etcd members are located on the control plane machines. Each control plane machine requires connectivity to [etcd server][etcd-ports], [etcd peer][etcd-ports] and [etcd-metrics][etcd-ports] on every other control plane machine.

* OpenShift SDN

    All the machines require connectivity to certain reserved ports on every other machine to establish in-cluster networking. For more details refer [doc][sdn-ports].

* Kubernetes NodePort

    All the machines require connectivity to Kubernetes NodePort range 30000-32767 on every other machine for OpenShift platform components.

* OpenShift reserved

    All the machines require connectivity to reserved port ranges 10250-12252 and 9000-9999 on every other machine for OpenShift platform components.

### Connectivity during machine boot

All the RHCOS machines require network in `initramfs` during boot to fetch Ignition config from the Machine Config Server [machine-config-server].

## DNS requirements

* Kubernetes API

    OpenShift 4.x requires the DNS records `api.$cluster_name.$base_domain` and `api-int.$cluster_name.$base_domain` to point to the Load balancer targeting the control plane machines. Both records must be resolvable from all the nodes within the cluster. The `api.$cluster_name.$base_domain` must also be resolvable by clients external to the cluster.

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

The OpenShift installer uses an [Install Config](../customization.md#platform-customization) to drive all install time configuration. Note that, the Openshift Installer currently can not create the Install Config for baremetal and you will have to manually create this file.

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
  replicas: 1
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

## Getting OS related assets for machines

TODO RHEL CoreOS does not have assets for bare-metal.

## Booting machines with RHCOS and Ignition configs

### Required kernel parameters during PXE

* `rd.neednet=1`: [CoreOS Installer][coreos-installer] needs internet access to fetch the OS image that needs to be installed on the machine.

* CoreOS Installer [arguments][coreos-installer-args] are required to be configured to install RHCOS and setup the Ignition config file for that machine.

## Watching your installation

### Monitor for bootstrap-complete

The administrators can use the `wait-for bootstrap-complete` target of the OpenShift Installer to monitor cluster bootstrapping. The command succeeds when it notices `bootstrap-complete` event from Kubernetes APIServer. This event is generated by the bootstrap machine after the Kubernetes APIServer has been bootstrapped on the control plane machines. For example,

```console
$ openshift-install --dir test-bare-metal wait-for bootstrap-complete
INFO Waiting up to 30m0s for the Kubernetes API at https://api.test.example.com:6443...
INFO API v1.12.4+c53f462 up
INFO Waiting up to 30m0s for the bootstrap-complete event...
```

## Monitor for cluster completion

The administrators can use the `wait-for install-complete` target of the OpenShift Installer to monitor cluster completion. The command succeeds when it notices that Cluster Version Operator has completed rolling out the OpenShift cluster from Kubernetes APIServer.

```console
$ openshift-install wait-for install-complete
INFO Waiting up to 30m0s for the cluster to initialize...
```

## Example Bare-Metal UPI deployment

Terraform [templates][upi-metal-example] provides an example of using OpenShift Installer to create an bare-metal UPI OpenShift cluster on Packet.net

### Overview

* Compute: Uses Packet.net to deploy bare-metal machines.
    Uses [matchbox] to serve PXE scripts and Ignition configs for bootstrap, control plane and worker machines.
    Uses `public` IPv4 addresses for each machine, so that all the machines are accessible on the internet.

* DNS and Load Balancing
    Uses AWS [Route53](aws-route53) to configure the all the DNS records.
    Uses Round-Robin DNS [RRDNS][rrdns] in place of load balancing solutions.

Refer to the pre-requisites for using the example [here][upi-metal-example-pre-req]

### Creating the cluster

#### Installer assets

Use the OpenShift Installer to create [Ignition configs](#getting-ignition-configs-for-machines) that will be used to create bootstrap, control plane and worker machines.

#### Terraform variable file

Use the [example][upi-metal-example-tfvar] Terraform variable file to create terraform variable file, and edit the `tfvars` file on your favorite editor.

```sh
cp terraform.tfvars{.example,}
```

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

Use the bootstrap [monitoring](#monitor-for-bootstrap-complete) to track when cluster bootstrapping has finished. After the Kubernetes APIServer has been bootstrapped on the control plane machines, the bootstrap machine can be removed from the API pool by following:

```sh
terraform apply -auto-approve -var=bootstrap_dns="false"
```

NOTE: The bootstrap resources like the bootstrap machines currently cannot be removed using terraform. You can use the Packet.net console to remove the bootstrap machine. All the resources will be cleaned up by `terraform destroy`

### Approving server certificates for nodes

To allow Kube APIServer to communicate with the kubelet running on nodes for logs, rsh etc. The administrator needs to approve the CSR [requests][csr-requests] generated by each kubelet.

You can approve all `Pending` CSR requests using,

```sh
oc get csr -ojson | jq -r '.items[] | select(.status == {} ) | .metadata.name' | xargs oc adm certificate approve
```

### Updating image-registry to emptyDir storage backend

The Cluster Image Registry [Operator][cluster-image-registry-operator] does not pick an storage backend for `None` platform. Therefore, the cluster operator is will be stuck in progressing because it is waiting for administrator to [configure][cluster-image-registry-operator-configuration] a storage backend for the image-registry. You can pick `emptyDir` for non-production clusters by following:

```sh
oc patch configs.imageregistry.operator.openshift.io cluster --type merge --patch '{"spec":{"storage":{"emptyDir":{}}}}'
```

#### Monitoring cluster completion

Use the cluster finish [monitoring](#monitor-for-cluster-completion) to track when cluster has completely finished deploying.

#### Destroying the cluster

Use terraform [destroy][terraform-destroy] to destroy all the resources for the cluster. For example,

```console
terraform destroy -auto-approve
```

[aws-route53]: https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/Welcome.html
[cluster-image-registry-operator-configuration]: https://github.com/openshift/cluster-image-registry-operator#registry-resource
[cluster-image-registry-operator]: https://github.com/openshift/cluster-image-registry-operator#image-registry-operator
[coreos-installer-args]: https://github.com/coreos/coreos-installer#kernel-command-line-options-for-coreos-installer-running-in-the-initramfs
[coreos-installer]: https://github.com/coreos/coreos-installer#coreos-installer
[csr-requests]: https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/#requesting-a-certificate
[etcd-ports]: https://github.com/openshift/origin/pull/21520
[machine-config-server]: https://github.com/openshift/machine-config-operator/blob/master/docs/MachineConfigServer.md
[matchbox]: https://github.com/coreos/matchbox
[openshift-router]: https://github.com/openshift/cluster-ingress-operator#openshift-ingress-operator
[rrdns]: https://tools.ietf.org/html/rfc1794
[sdn-ports]: https://github.com/openshift/origin/pull/21520
[terraform-apply]: https://www.terraform.io/docs/commands/apply.html
[terraform-destroy]: https://www.terraform.io/docs/commands/destroy.html
[terraform-init]: https://www.terraform.io/docs/commands/init.html
[terraform-providers]: https://www.terraform.io/docs/providers/
[upi-metal-example-pre-req]: ../../../upi/metal/README.md#pre-requisites
[upi-metal-example-tfvar]: ../../../upi/metal/terraform.tfvars.example
[upi-metal-example]: ../../../upi/metal/README.md
