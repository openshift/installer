data "aws_availability_zones" "azs" {}

data "aws_ami" "coreos_ami" {
  most_recent = true

  filter {
    name   = "name"
    values = ["CoreOS-stable-*"]
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
    values = ["595879546273"]
  }
}

module "vpc" {
  source          = "./vpc"
  external_vpc_id = "${var.external_vpc_id}"
  vpc_cid_block   = "${var.vpc_cid_block}"
  cluster_name    = "${var.cluster_name}"
}

data "aws_vpc" "cluster_vpc" {
  id = "${module.vpc.vpc_id}"
}

module "etcd" {
  source = "./etcd"

  vpc_id          = "${data.aws_vpc.cluster_vpc.id}"
  node_count      = "${var.az_count}"
  ssh_key         = "${aws_key_pair.ssh-key.id}"
  dns_zone        = "${aws_route53_zone.tectonic-int.zone_id}"
  coreos_ami      = "${data.aws_ami.coreos_ami.id}"
  etcd_subnets    = ["${aws_subnet.etcd_subnet.*.id}"]
  tectonic_domain = "${var.tectonic_domain}"
  cluster_name    = "${var.cluster_name}"
}
