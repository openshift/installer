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
  tectonic_aws_external_vpc_id = "${var.tectonic_aws_external_vpc_id}"
  tectonic_aws_vpc_cidr_block   = "${var.tectonic_aws_vpc_cidr_block}"
  tectonic_cluster_name    = "${var.tectonic_cluster_name}"
}

data "aws_vpc" "cluster_vpc" {
  id = "${module.vpc.vpc_id}"
}

module "etcd" {
  source = "./etcd"

  vpc_id       = "${data.aws_vpc.cluster_vpc.id}"
  node_count   = "${var.tectonic_aws_az_count == 5 ? 5 : 3}"
  ssh_key      = "${aws_key_pair.ssh-key.id}"
  dns_zone     = "${aws_route53_zone.tectonic-int.zone_id}"
  coreos_ami   = "${data.aws_ami.coreos_ami.id}"
  etcd_subnets = ["${aws_subnet.etcd_subnet.*.id}"]
  tectonic_base_domain  = "${var.tectonic_base_domain}"
  tectonic_cluster_name = "${var.tectonic_cluster_name}"
}
