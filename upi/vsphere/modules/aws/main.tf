locals {
  tags = "${merge(map(
    "kubernetes.io/cluster/${var.cluster_id}", "owned"
  ), var.aws_extra_tags)}"
}

provider "aws" {
  region = "${var.aws_region}"
}

module "dns" {
  source = "./route53"

  api_external_lb_dns_name = "${module.vpc.aws_lb_api_external_dns_name}"
  api_external_lb_zone_id  = "${module.vpc.aws_lb_api_external_zone_id}"
  api_internal_lb_dns_name = "${module.vpc.aws_lb_api_internal_dns_name}"
  api_internal_lb_zone_id  = "${module.vpc.aws_lb_api_internal_zone_id}"
  base_domain              = "${var.base_domain}"
  cluster_domain           = "${var.cluster_domain}"
  cluster_id               = "${var.cluster_id}"
  etcd_count               = "${var.control_plane_count}"
  etcd_ip_addresses        = "${var.control_plane_ip_addresses}"
  tags                     = "${local.tags}"
  vpc_id                   = "${var.vpc_id}"
}

module "vpc" {
  source = "./vpc"

  cidr_block        = "${var.machine_cidr}"
  cluster_id        = "${var.cluster_id}"
  region            = "${var.aws_region}"
  machine_cidr      = "${var.machine_cidr}"
  vpc_id            = "${var.vpc_id}"
  public_subnet_id  = "${var.aws_public_subnet_id}"
  private_subnet_id = "${var.aws_private_subnet_id}"

  availability_zones = "${distinct(concat(var.aws_control_plane_availability_zones, var.aws_compute_availability_zones))}"

  tags = "${local.tags}"
}

module "bootstrap" {
  source = "./bootstrap"

  instance_count           = "${var.bootstrap_count}"
  cluster_id               = "${var.cluster_id}"
  target_group_arns        = "${module.vpc.aws_lb_target_group_arns}"
  target_group_arns_length = "${module.vpc.aws_lb_target_group_arns_length}"
  availability_zone        = "${var.aws_availability_zone}"

  ip_address = "${var.bootstrap_ip_address}"

  tags = "${local.tags}"
}

module "control_plane" {
  source = "./control-plane"

  cluster_id = "${var.cluster_id}"

  tags = "${local.tags}"

  instance_count           = "${var.control_plane_count}"
  ip_addresses             = "${var.control_plane_ip_addresses}"
  target_group_arns        = "${module.vpc.aws_lb_target_group_arns}"
  target_group_arns_length = "${module.vpc.aws_lb_target_group_arns_length}"
  availability_zone        = "${var.aws_availability_zone}"
}

module "compute" {
  source = "./compute"

  cluster_id = "${var.cluster_id}"

  tags = "${local.tags}"

  instance_count           = "${var.compute_count}"
  ip_addresses             = "${var.compute_ip_addresses}"
  target_group_arns        = "${module.vpc.aws_lb_target_group_arns}"
  target_group_arns_length = "${module.vpc.aws_lb_target_group_arns_length}"
  availability_zone        = "${var.aws_availability_zone}"
}
