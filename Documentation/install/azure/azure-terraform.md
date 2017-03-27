# Install Tectonic on Azure with Terraform

Following this guide will deploy a Tectonic cluster within your Azure account.

Generally, the Azure platform templates adhere to the standards defined by the project [conventions][conventions] and [generic platform requirements][generic]. This document aims to document the implementation details specific to the Azure platform.

## Prerequsities

 - **DNS** - Setup your DNS zone in a resource group called `tectonic-dns-group` or specify a different resource group using the `tectonic_azure_dns_resource_group` variable below. We use a separate resource group assuming that you have a zone that you already want to use. Follow the [docs to set one up][azure-dns].
 - **Make** - This guide uses `make` to download a customized version of Terraform, which is pinned to a specific version and includes required plugins.
 - **Tectonic Account** - Register for a [Tectonic Account][register], which is free for up to 10 nodes. You will need to provide the cluster license and pull secret below.

## Getting Started

First, clone the Tectonic Installer repository in a convenient location:

```
$ git clone https://github.com/coreos/tectonic-installer.git
$ git checkout tags/v0.0.1
```

Download Terraform with via `make`. This will download the pinned Terraform binary and modules:

```
$ cd tectonic-installer
$ make terraform-download
```

After downloading, you will need to source this new binary in your `$PATH`. This is important, especially if you have another verison of Terraform installed. Run this command to add it to your path:

```
$ export PATH=/path/to/tectonic-installer/bin/terraform:$PATH
```

You can double check that you're using the binary that was just downloaded:

```
$ which terraform
/Users/coreos/tectonic-installer/bin/terraform/terraform
```

Next, get the modules that Terraform will use to create the cluster resources:

```
$ terraform get platforms/azure
Get: file:///Users/tectonic-installer/modules/azure/vnet
Get: file:///Users/tectonic-installer/modules/azure/etcd
Get: file:///Users/tectonic-installer/modules/azure/master
Get: file:///Users/tectonic-installer/modules/azure/worker
Get: file:///Users/tectonic-installer/modules/azure/dns
Get: file:///Users/tectonic-installer/modules/bootkube
Get: file:///Users/tectonic-installer/modules/tectonic
```

Generate credentials using the Azure CLI. If you're not logged in, execute `az login` first. See the [docs][login] for more info.

```
$ az ad sp create-for-rbac -n "http://tectonic" --role contributor
Retrying role assignment creation: 1/24
Retrying role assignment creation: 2/24
{
 "appId": "generated-app-id",
 "displayName": "azure-cli-2017-01-01",
 "name": "http://tectonic-coreos",
 "password": "generated-pass",
 "tenant": "generated-tenant"
}
```

Export variables that correspond to the data that was just generated. The subscription is your Azure Subscription ID.

```
$ export ARM_SUBSCRIPTION_ID=abc-123-456
$ export ARM_CLIENT_ID=generated-app-id
$ export ARM_CLIENT_SECRET=generated-pass
$ export ARM_TENANT_ID=generated-tenant
```

Last, let's create a local build directory `build/<cluster-name>` which holds all module references, Terraform state files, and custom variable files:

```
$ PLATFORM=azure CLUSTER=my-cluster make localconfig
```

Now we're ready to specify our cluster configuration.

## Customize the deployment

Use this example to customize your cluster configuration. A few fields require special consideration:

 - **tectonic_base_domain** - domain name that is set up with in a resource group, as described in the prerequisites.
 - **tectonic_pull_secret_path** - path on disk to your downloaded pull secret. You can find this on your [Account dashboard][account].
 - **tectonic_license_path** - path on disk to your downloaded Tectonic license. You can find this on your [Account dashboard][account].
 - **tectonic_admin_password_hash** - generate a hash with the [bcrypt-hash tool][bcrypt] that will be used for your admin user.

Here's an example of the full file:

**build/<cluster>/terraform.tfvars**

```
tectonic_worker_count = "4"

tectonic_master_count = "2"

tectonic_etcd_count = "1"

tectonic_base_domain = "azure.example.com"

tectonic_cluster_name = "mycluster"

tectonic_pull_secret_path = "/Users/coreos/Downloads/config.json"

tectonic_license_path = "/Users/coreos/Downloads/tectonic-license.txt"

tectonic_cl_channel = "stable"

tectonic_admin_email = "first.last@example.com"

tectonic_admin_password_hash = "<redacted - generate with bcrypt tool>"

tectonic_ca_cert = "" # path on disk, keep empty to generate one

tectonic_ca_key = "" # path on disk, keep empty to generate one

tectonic_azure_ssh_key = "/Users/coreos/.ssh/id_rsa.pub"

tectonic_azure_vnet_cidr_block = "10.0.0.0/16"

tectonic_azure_etcd_vm_size = "Standard_DS2"

tectonic_azure_master_vm_size = "Standard_DS2"

tectonic_azure_worker_vm_size = "Standard_DS2"

tectonic_azure_location = "eastus"
```

## Deploy the cluster

Test out the plan before deploying everything:

```
$ PLATFORM=azure CLUSTER=my-cluster make plan
```

Next, deploy the cluster:

```
$ PLATFORM=azure CLUSTER=my-cluster make apply
```

This should run for a little bit, and when complete, your Tectonic cluster should be ready.

If you encounter any issues, check the known issues and workarounds below.

To delete your cluster, run:

```
$ PLATFORM=azure CLUSTER=my-cluster make destroy
```

### Known issues and workarounds

See the [troubleshooting][troubleshooting] document for work arounds for bugs that are being tracked.

## Scaling the cluster

To scale worker nodes, adjust `tectonic_worker_count` in `terraform.vars` and invoke `terraform apply -target module.workers platforms/azure`.

## Under the hood

### Top-level templates

* The top-level templates that invoke the underlying component modules reside `./platforms/azure`
* Point terraform to this location to start applying: `terraform apply ./platforms/azure`

### Etcd nodes

* Discovery is currently not implemented so DO NOT scale the etcd cluster to more than 1 node, for now.
* Etcd cluster nodes are managed by the terraform module `modules/azure/etcd`
* Node VMs are created as stand-alone instances (as opposed to VM scale sets). This is mostly historical and could change.
* A load-balancer fronts the etcd nodes to provide a simple discovery mechanism, via a VIP + DNS record.
* Currently, the LB is configured with a public IP address. This is not optimal and it should be converted to an internal LB.

### Master nodes

* Master node VMs are managed by the templates in `modules/azure/master`
* An Azure VM Scaling Set resource is used to spin-up multiple identical VM configured as master nodes.
* Master VMs all share the same identical Ignition config
* Master nodes are fronted by one load-balancer for the API one for the Ingress controller.
* The API LB is configured with SourceIP session stickiness, to ensure that TCP (including SSH) sessions from the same client land reliably on the same master node. This allows for provisioning the assets and starting bootkube reliably via SSH.
* a `null_resource` terraform provisioner in the tectonic.tf top-level template will copy the assets and run bootkube automatically on one of the masters.
* make sure the SSH key specifyied in the tfvars file is also added to the SSH agent on the machine running terraform. Without this, terraform is not able to SSH copy the assets and start bootkube. Also make sure that the SSH known_hosts file doesn't have old records of the API DNS name (fingerprints will not match).

### Worker nodes

* Worker node VMs are managed by the templates in `modules/azure/worker`
* An Azure VM Scaling Set resource is used to spin-up multiple identical VM configured as worker nodes.
* Worker VMs all share the same identical Ignition config
* Worker nodes are not fronted by any LB and don't have public IP addresses. They can be accessed through SSH from any of the master nodes.

[conventions]: ../conventions.md
[generic]: ../generic-platform.md
[register]: https://account.tectonic.com/signup/summary/tectonic-2016-12
[account]: https://account.tectonic.com
[bcrypt]: https://github.com/coreos/bcrypt-tool/releases/tag/v1.0.0
[plan-docs]: https://www.terraform.io/docs/commands/plan.html
[copy-docs]: https://www.terraform.io/docs/commands/apply.html
[troubleshooting]: ../troubleshooting.md
[login]: https://docs.microsoft.com/en-us/cli/azure/get-started-with-azure-cli
[azure-dns]: https://docs.microsoft.com/en-us/azure/dns/dns-getstarted-portal
