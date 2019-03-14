# Bootstrap Module

This [Terraform][] [module][] manages [vSphere][] DNS resources. 
It uses [implicit provider inheritance][implicit-provider-inheritance] to access the [vSphere provider][vSphere-provider].

## Example

Set up a `main.tf` with:

```hcl
module "dns" {
  source = "./route53"

  base_domain       = "example.com"
  cluster_domain    = "openshift.example.com"
  cluster_id        = "my-cluster"
  etcd_count        = "3"
  etcd_ip_addresses = ["10.0.0.2", "10.0.0.3", "10.0.0.4"]
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
