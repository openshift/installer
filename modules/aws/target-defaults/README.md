# AWS Default Module

This [Terraform][] [module][] gives us a central location for [AWS][]-specific defaults.
For example, on AWS our default [etcd][] count depends on the availability zone count.
On libvirt, on the other hand, our default etcd count is one.
Because Terraform's [HashiCorp Configuration Language (HCL)[hcl] does not allow us to declare such conditional defaults directly, we use this module to convert user-supplied default values to values appropriate to the target system.

## Example

From the module directory:

```console
$ terraform init
$ terraform apply --var region=us-east-1  # currently 6 zones
$ terraform output etcd_count
5
$ terraform apply --var region=us-east-2
$ terraform output etcd_count
3
$ terraform apply --var region=us-east-2 --var etcd_count=1  # currently 3 zones
$ terraform output etcd_count
1
```

When you're done, clean up by removing the `.terraform` directory created by `init` and the `terraform.tfstate*` files created by `apply`.

[AWS]: https://aws.amazon.com/
[etcd]: https://github.com/coreos/etcd
[module]: https://www.terraform.io/docs/modules/
[Terraform]: https://www.terraform.io/
