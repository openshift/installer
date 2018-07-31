# Container Linux AMI Module

This [Terraform][] [module][] supports `latest` versions for [Container Linux][container-linux] release channels and returns an appropriate [AMI][].

## Example

From the module directory:

```console
$ terraform init
$ terraform apply --var region=us-east-1
$ terraform output id
ami-ab6963d4
$ terraform apply --var region=us-east-1 --var release_channel=alpha
$ terraform output id
ami-985953e7
$ terraform apply --var region=us-east-2 --var release_channel=alpha --var release_version=1814.0.0
$ terraform output id
ami-c25f66a7
```

When you're done, clean up by removing the `.terraform` directory created by `init` and the `terraform.tfstate*` files created by `apply`.

[AMI]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html
[container-linux]: https://coreos.com/os/docs/latest/
[module]: https://www.terraform.io/docs/modules/
[Terraform]: https://www.terraform.io/
