# Bootstrap Module

This [Terraform][] [module][] manages [vSphere][] resource pool for VM use. 
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

module "resource_pool" {
  source = "./resource_pool"

  vsphere_cluster    	 = "vsphere-cluster"
  vsphere_datacenter 	 = "vsphere-datacenter"
  vsphere_resource_pool  = "vsphere-resource-pool"
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
