# Bootstrap Module

This [Terraform][] [module][] manages [vSphere][] resources only needed during cluster bootstrapping.
It uses [implicit provider inheritance][implicit-provider-inheritance] to access the [vSphere provider][vSphere-provider].

## Example

Set up a `main.tf` with:

```hcl
provider "vsphere" {
  user           = "administrator@vsphere.local"
  password       = "password"
  vsphere_server 	 = "vcsa.vmware.example.com"

  # If you have a self-signed cert
  allow_unverified_ssl = true
}

module "bootstrap" {
  source = "./bootstrap"

  bootstrap_ip       = "10.0.0.1"
  cluster_id         = "my-cluster"
  machine_cidr       = "10.0.0.0/24"
  vsphere_cluster    = "vsphere-cluster"
  vsphere_datacenter = "vsphere-datacenter"
  vsphere_datastore  = "vsphere-datastore"
  resource_pool_id   = "vsphere-pool-id"
  vm_base_domain     = "vmware.example.com"
  vm_network         = "VM Network"
  vm_template        = "rhocs-template"
}

```

Then run:

```console
$ terraform init
$ terraform plan
```

[vSphere]: https://www.vmware.com/products/vsphere.html
[vSphere-provider]: https://www.terraform.io/docs/providers/vsphere/
[implicit-provider-inheritance]: https://www.terraform.io/docs/modules/usage.html#implicit-provider-inheritance
[module]: https://www.terraform.io/docs/modules/
[Terraform]: https://www.terraform.io/
