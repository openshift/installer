# Container Linux Module

This [Terraform][] [module][] supports `latest` versions for [Container Linux][container-linux] release channels.

## Example

From the module directory:

```console
$ terraform init
$ terraform apply
$ terraform output version
1745.7.0
$ terraform apply --var release_channel=alpha
$ terraform output version
1828.0.0
$ terraform apply --var release_version=1814.0.0
$ terraform output version
1814.0.0
```

When you're done, clean up by removing the `.terraform` directory created by `init` and the `terraform.tfstate*` files created by `apply`.

[container-linux]: https://coreos.com/os/docs/latest/
[module]: https://www.terraform.io/docs/modules/
[Terraform]: https://www.terraform.io/
