# Install Tectonic on VMware with Terraform

Following this guide will deploy a Tectonic cluster within a VMware vSphere infrastructure .

Generally, the VMware platform templates adhere to the standards defined by the project [conventions][conventions] and [generic platform requirements][generic]. This document aims to document the implementation details specific to the VMware platform.

## Prerequsities

1. Download the latest Container Linux Stable OVA from; https://coreos.com/os/docs/latest/booting-on-vmware.html.
1. Import `coreos_production_vmware_ova.ova` into vCenter. Generally, all settings can be kept as is. Consider "thin" provisioning and naming the template with CoreOS Container Linux Version.
1. Resize the Virtual Machine Disk size to 30 GB or larger
1. In the Virtual Machine Configuration View select "vApp Options" tab and un-check "Enable vApp Options".
1. Convert the Container Linux image into a Virtual Machine template.
1. Pre-Allocated IP addresses for the cluster and pre-create DNS records

### DNS and IP address allocation

Prior to the start of setup create required DNS records. Below is a sample table of 3 etcd nodes, 2 master nodes and 2 worker nodes. 

| Record | Type | Value |
|------|-------------|:-----:|
|mycluster.mycompany.com | A | 192.168.246.30 |
|mycluster.mycompany.com | A | 192.168.246.31 |
|mycluster-k8s.mycompany.com | A | 192.168.246.20 |
|mycluster-k8s.mycompany.com | A | 192.168.246.21 |
|mycluster-worker-0.mycompany.com | A | 192.168.246.30 |
|mycluster-worker-1.mycompany.com | A | 192.168.246.31 |
|mycluster-master-0.mycompany.com | A | 192.168.246.20 |
|mycluster-master-1.mycompany.com | A | 192.168.246.21 |
|mycluster-etcd-0.mycompany.com | A | 192.168.246.10 |
|mycluster-etcd-1.mycompany.com | A | 192.168.246.11 |
|mycluster-etcd-2.mycompany.com | A | 192.168.246.12 |

See [Tectonic on Baremetal DNS documentation][baremetaldns] for general DNS Requirements.

### Tectonic Account

Register for a [Tectonic Account][register], which is free for up to 10 nodes. You must provide the cluster license and pull secret during installation.

### ssh-agent

Ensure `ssh-agent` is running:
```
$ eval $(ssh-agent)
```

Add the SSH key that will be used for the Tectonic installation to `ssh-agent`:
```
$ ssh-add <path-to-ssh-key>
```

Verify that the SSH key identity is available to the ssh-agent:
```
$ ssh-add -L
```

Reference the absolute path of the **_public_** component of the SSH key in `tectonic_vmware_ssh_authorized_key`.

Without this, terraform is not able to SSH copy the assets and start bootkube.
Also make sure that the SSH known_hosts file doesn't have old records of the API DNS name (fingerprints will not match).

## Getting Started

The following steps must be executed on a machine that has network connectivity to VMware vCenter API and SSH access to Tectonic Master Server(s).

### Download and extract Tectonic Installer

Open a new terminal, and run the following commands to download and extract Tectonic Installer.

```bash
$ curl -O https://releases.tectonic.com/tectonic-1.7.1-tectonic.1.tar.gz # download
$ tar xzvf tectonic-1.7.1-tectonic.1.tar.gz # extract the tarball
$ cd tectonic
```

## Customize the deployment

Customizations to the base installation live in `examples/terraform.tfvars.<flavor>`. Export a variable that will be the cluster identifier:

```
$ export CLUSTER=my-cluster
```

Create a build directory to hold customizations and copy the example file into it:

```
$ mkdir -p build/${CLUSTER}
$ cp examples/terraform.tfvars.vmware build/${CLUSTER}/terraform.tfvars
$ cd build/${CLUSTER}/
```

Edit the parameters with details of the VMware infrastructure. View all of the [VMware][vmware] specific options and [the common Tectonic variables][vars].

## Deploy the cluster

Get the modules that Terraform will use to create the cluster resources:

```
$ terraform get ../../platforms/vmware
```

Test out the plan before deploying everything:

```
$ terraform plan ../../platforms/vmware
```

Terraform will prompt for vSphere credentials:

```
provider.vsphere.password
  The user password for vSphere API operations.

  Enter a value: 

provider.vsphere.user
  The user name for vSphere API operations.

  Enter a value: 
```

Next, deploy the cluster:

```
$ terraform apply ../../platforms/vmware
```

Wait for `terraform apply` until all tasks complete. Tectonic cluster should be ready upon completion of apply command. If any issues arrise please check the known issues and workarounds below.

## Access the cluster

The Tectonic Console should be up and running after the containers have downloaded. Console can be access by the DNS name configured as `tectonic_vmware_ingress_domain` in the `terraform.tfvars` variables file.

Credentials and secrets for Tectonic can be found in `/generated` folder including the CA if generated, and a kubeconfig. Use the kubeconfig file to control the cluster with `kubectl`:

```
$ export KUBECONFIG=generated/auth/kubeconfig
$ kubectl cluster-info
```

## Working with the cluster

### Scaling Tectonic VMware clusters

This document describes how to add cluster nodes to Tectonic clusters on VMware.

#### Scaling worker nodes

To scale worker nodes, adjust `tectonic_worker_count`, `tectonic_vmware_worker_hostnames` and `tectonic_vmware_worker_ip` variables in `terraform.tfvars` and run:

```
$ terraform plan \
  ../../platforms/vmware
$ terraform apply \
  ../../platforms/vmware
```
Shortly after running `terraform apply` new worker machines will appear on Tectonic console. This change may take several minutes.

#### Scaling master nodes

To scale master nodes, adjust `tectonic_master_count`, `tectonic_vmware_master_hostnames` and `tectonic_vmware_master_ip` variables in `terraform.tfvars` and run:

```
$ terraform plan \
  ../../platforms/vmware
$ terraform apply \
  ../../platforms/vmware
```
Shortly after running `terraform apply` master machines will appear on Tectonic console. This change may take several minutes.  

Make sure to add the new Controller nodes' IP addresses in DNS for `tectonic_vmware_controller_domain` variable or update the Load balancer to include new Controller nodes.

## Known issues and workarounds

See the [troubleshooting][troubleshooting] document for workarounds for bugs that are being tracked

## Delete the cluster

To delete Tectonic cluster, run:

```
$ terraform destroy ../../platforms/vmware
```

[register]: https://account.coreos.com
[baremetaldns]: https://coreos.com/tectonic/docs/latest/install/bare-metal/#dns 
[conventions]: ../../conventions.md
[generic]: ../../generic-platform.md
[downloadterraform]: https://www.terraform.io/downloads.html
[vmware]: ../../variables/vmware.md
[vars]: ../../variables/config.md
[troubleshooting]: ../../troubleshooting/faq.md
