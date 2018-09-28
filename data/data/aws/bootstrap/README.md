# Bootstrap Module

This [Terraform][] [module][] manages [AWS][] resources only needed during cluster bootstrapping.
It uses [implicit provider inheritance][implicit-provider-inheritance] to access the [AWS provider][AWS-provider].

## Example

Set up a `main.tf` with:

```hcl
provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "example" {
}

resource "aws_vpc" "example" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_subnet" "example" {
  vpc_id     = "${aws_vpc.example.id}"
  cidr_block = "${aws_vpc.example.cidr_block}"
}

module "bootstrap" {
  source = "github.com/openshift/installer//data/data/aws/bootstrap"

  ami            = "ami-0af8953af3ec06b7c"
  bucket         = "${aws_s3_bucket.example.id}"
  cluster_name   = "my-cluster"
  ignition       = "{\"ignition\": {\"version\": \"2.2.0\"}}",
  subnet_id      = "${aws_subnet.example.id}"
}
```

Then run:

```console
$ terraform init
$ terraform plan
```

[AWS]: https://aws.amazon.com/
[AWS-provider]: https://www.terraform.io/docs/providers/aws/
[implicit-provider-inheritance]: https://www.terraform.io/docs/modules/usage.html#implicit-provider-inheritance
[module]: https://www.terraform.io/docs/modules/
[Terraform]: https://www.terraform.io/
