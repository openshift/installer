# AWS Etcd Module

This [Terraform][] [module][] makes it easy to create [etcd][] nodes on [AWS][].

Read the [etcd recommended hardware guide][hardware] for best performance.

## Example

```hcl
provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "etcd_ignition" {
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

resource "aws_security_group" "etcd" {
  vpc_id = "${aws_vpc.example.id}"

  ingress {
    from_port   = 2379
    to_port     = 2380
    protocol    = "tcp"
    cidr_blocks = ["${aws_subnet.example.cidr_block}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${aws_subnet.example.cidr_block}"]
  }
}

module "etcd" {
  source = "github.com/openshift/installer//modules/aws/etcd"

  base_domain    = "openshift.example.com"
  cluster_id     = "123"
  cluster_name   = "my-cluster"
  instance_count = "3"
  s3_bucket      = "${aws_s3_bucket.etcd_ignition.id}"
  sg_ids         = ["${aws_security_group.etcd.id}"]
  subnets        = ["${aws_subnet.example.id}"]
}
```

You can set `container_linux_channel` and `container_linux_version` if you need a specific [Container Linux][container-linux] install.
Alternatively, you can set `ec2_ami` directly if you want to use an [AMI][] that is not Container Linux.

[AMI]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html
[AWS]: https://aws.amazon.com/
[container-linux]: https://coreos.com/os/docs/latest/
[etcd]: https://github.com/coreos/etcd
[hardware]: https://github.com/coreos/etcd/blob/v3.3.8/Documentation/op-guide/hardware.md#example-hardware-configurations
[module]: https://www.terraform.io/docs/modules/
[Terraform]: https://www.terraform.io/
