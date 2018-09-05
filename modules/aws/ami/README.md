# Container Linux AMI Module

This [Terraform][] [module][] supports `latest` versions for [Container Linux][container-linux] release channels and returns an appropriate [AMI][].
It uses [implicit provider inheritance][implicit-provider-inheritance] to access the [AWS provider][AWS-provider].

## Example

Set up a `main.tf` with:

```hcl
provider "aws" {
  region = "us-east-1"
}

module "ami" {
  source = "github.com/openshift/installer//modules/aws/ami"
}

output "ami" {
  value = "${module.ami.id}"
}
```

You can set `release_channel` and `release_version` if you need a specific Container Linux install.

Then run:

```console
$ terraform init
$ terraform apply
$ terraform output ami
ami-00cc4337762ba4a52
```

[AMI]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html
[AWS-provider]: https://www.terraform.io/docs/providers/aws/
[container-linux]: https://coreos.com/os/docs/latest/
[implicit-provider-inheritance]: https://www.terraform.io/docs/modules/usage.html#implicit-provider-inheritance
[module]: https://www.terraform.io/docs/modules/
[Terraform]: https://www.terraform.io/
