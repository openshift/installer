locals {
  private_zone_id = "${aws_route53_zone.int.zone_id}"

  tags = "${merge(map(
      "openshiftClusterID", "${var.cluster_id}"
    ), var.aws_extra_tags)}"
}

provider "aws" {
  region = "${var.aws_region}"
}

module "bootstrap" {
  source = "./bootstrap"

  ami                      = "${var.aws_ec2_ami_override}"
  cluster_name             = "${var.cluster_name}"
  ignition                 = "${var.ignition_bootstrap}"
  subnet_id                = "${module.vpc.public_subnet_ids[0]}"
  target_group_arns        = "${module.vpc.aws_lb_target_group_arns}"
  target_group_arns_length = "${module.vpc.aws_lb_target_group_arns_length}"
  vpc_id                   = "${module.vpc.vpc_id}"
  vpc_security_group_ids   = "${list(module.vpc.control_plane_sg_id)}"

  tags = "${merge(map(
      "Name", "${var.cluster_name}-bootstrap",
    ), local.tags)}"
}

module "controlplane" {
  source = "./controlplane"

  cluster_id   = "${var.cluster_id}"
  cluster_name = "${var.cluster_name}"
  ec2_type     = "${var.aws_control_plane_ec2_type}"

  tags = "${merge(map(
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
    ), local.tags)}"

  instance_count           = "${var.control_plane_count}"
  machine_pool_name        = "${var.control_plane_machine_pool_name}"
  control_plane_sg_ids     = "${list(module.vpc.control_plane_sg_id)}"
  root_volume_iops         = "${var.aws_control_plane_root_volume_iops}"
  root_volume_size         = "${var.aws_control_plane_root_volume_size}"
  root_volume_type         = "${var.aws_control_plane_root_volume_type}"
  subnet_ids               = "${module.vpc.private_subnet_ids}"
  target_group_arns        = "${module.vpc.aws_lb_target_group_arns}"
  target_group_arns_length = "${module.vpc.aws_lb_target_group_arns_length}"
  ec2_ami                  = "${var.aws_ec2_ami_override}"
  user_data_ign            = "${var.ignition_control_plane}"
}

module "iam" {
  source = "./iam"

  cluster_name = "${var.cluster_name}"

  tags = "${merge(map(
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
    ), local.tags)}"
}

module "dns" {
  source = "./route53"

  api_external_lb_dns_name = "${module.vpc.aws_lb_api_external_dns_name}"
  api_external_lb_zone_id  = "${module.vpc.aws_lb_api_external_zone_id}"
  api_internal_lb_dns_name = "${module.vpc.aws_lb_api_internal_dns_name}"
  api_internal_lb_zone_id  = "${module.vpc.aws_lb_api_internal_zone_id}"
  base_domain              = "${var.base_domain}"
  cluster_name             = "${var.cluster_name}"
  private_zone_id          = "${local.private_zone_id}"
}

module "vpc" {
  source = "./vpc"

  base_domain  = "${var.base_domain}"
  cidr_block   = "${var.machine_cidr}"
  cluster_id   = "${var.cluster_id}"
  cluster_name = "${var.cluster_name}"
  region       = "${var.aws_region}"

  tags = "${merge(map(
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
    ), local.tags)}"
}

resource "aws_route53_record" "etcd_a_nodes" {
  count   = "${var.control_plane_count}"
  type    = "A"
  ttl     = "60"
  zone_id = "${local.private_zone_id}"
  name    = "${var.cluster_name}-etcd-${count.index}"
  records = ["${module.controlplane.ip_addresses[count.index]}"]
}

resource "aws_route53_record" "etcd_cluster" {
  type    = "SRV"
  ttl     = "60"
  zone_id = "${local.private_zone_id}"
  name    = "_etcd-server-ssl._tcp.${var.cluster_name}"
  records = ["${formatlist("0 10 2380 %s", aws_route53_record.etcd_a_nodes.*.fqdn)}"]
}

resource "aws_route53_zone" "int" {
  name          = "${var.base_domain}"
  force_destroy = true

  vpc {
    vpc_id = "${module.vpc.vpc_id}"
  }

  tags = "${merge(map(
      "Name", "${var.cluster_name}_int",
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
    ), local.tags)}"
}
