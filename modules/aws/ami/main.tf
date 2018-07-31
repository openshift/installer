provider "aws" {
  region  = "${var.region}"
  version = "1.8.0"
}

locals {
  ami_owner = "595879546273"
  arn       = "aws"
}

module "container_linux" {
  source = "../../container_linux"

  release_channel = "${var.release_channel}"
  release_version = "${var.release_version}"
}

data "aws_ami" "coreos_ami" {
  filter {
    name   = "name"
    values = ["CoreOS-${var.release_channel}-${module.container_linux.version}-*"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "owner-id"
    values = ["${local.ami_owner}"]
  }
}
