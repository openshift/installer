# Azure platform architecture notes

Generally, the Azure platform templates adhere to the standards defined in *conventions.md* and *generic-platform.md*.

This document aims to document the implementation details specific to the Azure platform.

## Example cluster configuration (tfvars file)

Use this example to customize your cluster configuration. Save it as `terraform.tfvars` in the folder where your build state will reside.
When you invoke `terraform apply <path/to/platforms/azure>` in this folder, the tfvars file will be automatically picked up.

```
tectonic_worker_count = "4"

tectonic_master_count = "2"

tectonic_etcd_count = "1"

tectonic_base_domain = "azure.ifup.org"

tectonic_cluster_name = "mycluster"

tectonic_dns_name = "mycluster"

tectonic_pull_secret_path = "/Users/alex/Downloads/tectonic-quay.json"

tectonic_license_path = "/Users/alex/Downloads/tectonic-license.txt"

tectonic_cl_channel = "stable"

tectonic_admin_email = "alex.somesan@coreos.com"

tectonic_admin_password_hash = "<redacted - generate with bcrypt tool>"

tectonic_ca_cert = ""

tectonic_ca_key = ""

tectonic_azure_ssh_key = "/Users/alex/.ssh/id_rsa.pub"

tectonic_azure_vnet_cidr_block = "10.0.0.0/16"

tectonic_azure_etcd_vm_size = "Standard_DS3"

tectonic_azure_master_vm_size = "Standard_DS3"

tectonic_azure_worker_vm_size = "Standard_DS3"

tectonic_azure_location = "northeurope"
```

## Top-level templates

* The top-level templates that invoke the underlying component modules reside `./platforms/azure`
* Point terraform to this location to start applying: `terraform apply ./platforms/azure`

## Etcd nodes

* Discovery is currently not implemented so DO NOT scale the etcd cluster to more than 1 node, for now.
* Etcd cluster nodes are managed by the terraform module `modules/azure/etcd`
* Node VMs are created as stand-alone instances (as opposed to VM scale sets). This is mostly historical and could change.
* A load-balancer fronts the etcd nodes to provide a simple discovery mechanism, via a VIP + DNS record.
* Currently, the LB is configured with a public IP address. This is not optimal and it should be converted to an internal LB.

## Master nodes

* Master node VMs are managed by the templates in `modules/azure/master`
* An Azure VM Scaling Set resource is used to spin-up multiple identical VM configured as master nodes.
* Master VMs all share the same identical Ignition config
* Master nodes are fronted by one load-balancer for the API one for the Ingress controller.
* The API LB is configured with SourceIP session stickiness, to ensure that TCP (including SSH) sessions from the same client land reliably on the same master node. This allows for provisioning the assets and starting bootkube reliably via SSH.
* a `null_resource` terraform provisioner in the tectonic.tf top-level template will copy the assets and run bootkube automatically on one of the masters.
* make sure the SSH key specifyied in the tfvars file is also added to the SSH agent on the machine running terraform. Without this, terraform is not able to SSH copy the assets and start bootkube. Also make sure that the SSH known_hosts file doesn't have old records of the API DNS name (fingerprints will not match).

## Worker nodes

* Worker node VMs are managed by the templates in `modules/azure/worker`
* An Azure VM Scaling Set resource is used to spin-up multiple identical VM configured as worker nodes.
* Worker VMs all share the same identical Ignition config
* Worker nodes are not fronted by any LB and don't have public IP addresses. They can be accessed through SSH from any of the master nodes.