locals {
  private_endpoints = "${var.aws_endpoints == "public" ? false : true}"
  public_endpoints  = "${var.aws_endpoints == "private" ? false : true}"
  private_zone_id   = "${var.aws_external_private_zone != "" ? var.aws_external_private_zone : join("", aws_route53_zone.int.*.zone_id)}"
}

provider "aws" {
  region  = "${var.aws_region}"
  version = "1.39.0"

  assume_role {
    role_arn     = "${var.aws_installer_role == "" ? "" : "${var.aws_installer_role}"}"
    session_name = "OPENSHIFT_INSTALLER_${var.cluster_name}"
  }
}

module "bootstrap" {
  source = "./bootstrap"

  ami                         = "${var.aws_ec2_ami_override}"
  associate_public_ip_address = "${var.aws_endpoints != "private"}"
  cluster_name                = "${var.cluster_name}"
  iam_role                    = "${var.aws_master_iam_role_name}"
  ignition                    = "${var.ignition_bootstrap}"
  subnet_id                   = "${module.vpc.master_subnet_ids[0]}"
  target_group_arns           = "${module.vpc.aws_lb_target_group_arns}"
  target_group_arns_length    = "${module.vpc.aws_lb_target_group_arns_length}"
  vpc_security_group_ids      = ["${concat(var.aws_master_extra_sg_ids, list(module.vpc.master_sg_id))}"]

  tags = "${merge(map(
      "Name", "${var.cluster_name}-bootstrap",
      "tectonicClusterID", "${var.cluster_id}"
    ), var.aws_extra_tags)}"
}

module "masters" {
  source = "./master"

  base_domain              = "${var.base_domain}"
  cluster_id               = "${var.cluster_id}"
  cluster_name             = "${var.cluster_name}"
  ec2_type                 = "${var.aws_master_ec2_type}"
  extra_tags               = "${var.aws_extra_tags}"
  instance_count           = "${var.master_count}"
  master_iam_role          = "${var.aws_master_iam_role_name}"
  master_sg_ids            = "${concat(var.aws_master_extra_sg_ids, list(module.vpc.master_sg_id))}"
  public_endpoints         = "${local.public_endpoints}"
  root_volume_iops         = "${var.aws_master_root_volume_iops}"
  root_volume_size         = "${var.aws_master_root_volume_size}"
  root_volume_type         = "${var.aws_master_root_volume_type}"
  subnet_ids               = "${module.vpc.master_subnet_ids}"
  target_group_arns        = "${module.vpc.aws_lb_target_group_arns}"
  target_group_arns_length = "${module.vpc.aws_lb_target_group_arns_length}"
  ec2_ami                  = "${var.aws_ec2_ami_override}"
  user_data_ign            = "${var.ignition_master}"
}

module "iam" {
  source = "./iam"

  cluster_name    = "${var.cluster_name}"
  worker_iam_role = "${var.aws_worker_iam_role_name}"
}

module "dns" {
  source = "./route53"

  api_external_lb_dns_name = "${module.vpc.aws_lb_api_external_dns_name}"
  api_external_lb_zone_id  = "${module.vpc.aws_lb_api_external_zone_id}"
  api_internal_lb_dns_name = "${module.vpc.aws_lb_api_internal_dns_name}"
  api_internal_lb_zone_id  = "${module.vpc.aws_lb_api_internal_zone_id}"
  base_domain              = "${var.base_domain}"
  cluster_name             = "${var.cluster_name}"
  elb_alias_enabled        = true
  master_count             = "${var.master_count}"
  private_zone_id          = "${local.private_zone_id}"
  external_vpc_id          = "${module.vpc.vpc_id}"
  extra_tags               = "${var.aws_extra_tags}"
  private_endpoints        = "${local.private_endpoints}"
  public_endpoints         = "${local.public_endpoints}"
}

module "vpc" {
  source = "./vpc"

  base_domain     = "${var.base_domain}"
  cidr_block      = "${var.aws_vpc_cidr_block}"
  cluster_id      = "${var.cluster_id}"
  cluster_name    = "${var.cluster_name}"
  external_vpc_id = "${var.aws_external_vpc_id}"

  external_master_subnet_ids = "${compact(var.aws_external_master_subnet_ids)}"
  external_worker_subnet_ids = "${compact(var.aws_external_worker_subnet_ids)}"
  extra_tags                 = "${var.aws_extra_tags}"

  // empty map subnet_configs will have the vpc module creating subnets in all availabile AZs
  new_master_subnet_configs = "${var.aws_master_custom_subnets}"
  new_worker_subnet_configs = "${var.aws_worker_custom_subnets}"

  private_master_endpoints = "${local.private_endpoints}"
  public_master_endpoints  = "${local.public_endpoints}"
}

resource "aws_route53_record" "etcd_a_nodes" {
  count   = "${var.master_count}"
  type    = "A"
  ttl     = "60"
  zone_id = "${local.private_zone_id}"
  name    = "${var.cluster_name}-etcd-${count.index}"
  records = ["${module.masters.ip_addresses[count.index]}"]
}

resource "aws_route53_record" "etcd_cluster" {
  type    = "SRV"
  ttl     = "60"
  zone_id = "${local.private_zone_id}"
  name    = "_etcd-server-ssl._tcp.${var.cluster_name}"
  records = ["${formatlist("0 10 2380 %s", aws_route53_record.etcd_a_nodes.*.fqdn)}"]
}

resource "aws_route53_zone" "int" {
  count         = "${local.private_endpoints ? "${var.aws_external_private_zone == "" ? 1 : 0 }" : 0}"
  vpc_id        = "${module.vpc.vpc_id}"
  name          = "${var.base_domain}"
  force_destroy = true

  tags = "${merge(map(
      "Name", "${var.cluster_name}_int",
      "KubernetesCluster", "${var.cluster_name}",
      "tectonicClusterID", "${var.cluster_id}"
    ), var.aws_extra_tags)}"
}
