locals {
  private_endpoints = "${var.tectonic_aws_endpoints == "public" ? false : true}"
  public_endpoints  = "${var.tectonic_aws_endpoints == "private" ? false : true}"
}

provider "aws" {
  region  = "${var.tectonic_aws_region}"
  profile = "${var.tectonic_aws_profile}"
  version = "1.8.0"

  assume_role {
    role_arn     = "${var.tectonic_aws_installer_role == "" ? "" : "${var.tectonic_aws_installer_role}"}"
    session_name = "TECTONIC_INSTALLER_${var.tectonic_cluster_name}"
  }
}

data "aws_availability_zones" "azs" {}

# TNC
resource "aws_route53_zone" "tectonic_int" {
  count         = "${local.private_endpoints ? "${var.tectonic_aws_external_private_zone == "" ? 1 : 0 }" : 0}"
  vpc_id        = "${module.vpc.vpc_id}"
  name          = "${var.tectonic_base_domain}"
  force_destroy = true

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}_tectonic_int",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${var.tectonic_cluster_id}"
    ), var.tectonic_aws_extra_tags)}"
}

module "vpc" {
  source = "../../../modules/aws/vpc"

  base_domain     = "${var.tectonic_base_domain}"
  cidr_block      = "${var.tectonic_aws_vpc_cidr_block}"
  cluster_id      = "${var.tectonic_cluster_id}"
  cluster_name    = "${var.tectonic_cluster_name}"
  external_vpc_id = "${var.tectonic_aws_external_vpc_id}"

  external_master_subnet_ids = "${compact(var.tectonic_aws_external_master_subnet_ids)}"
  external_worker_subnet_ids = "${compact(var.tectonic_aws_external_worker_subnet_ids)}"
  extra_tags                 = "${var.tectonic_aws_extra_tags}"

  // empty map subnet_configs will have the vpc module creating subnets in all availabile AZs
  new_master_subnet_configs = "${var.tectonic_aws_master_custom_subnets}"
  new_worker_subnet_configs = "${var.tectonic_aws_worker_custom_subnets}"

  private_master_endpoints = "${local.private_endpoints}"
  public_master_endpoints  = "${local.public_endpoints}"
}

module "dns" {
  source = "../../../modules/dns/route53"

  api_external_elb_dns_name = "${module.vpc.aws_api_external_dns_name}"
  api_external_elb_zone_id  = "${module.vpc.aws_elb_api_external_zone_id}"
  api_internal_elb_dns_name = "${module.vpc.aws_api_internal_dns_name}"
  api_internal_elb_zone_id  = "${module.vpc.aws_elb_api_internal_zone_id}"
  api_ip_addresses          = "${module.vpc.aws_lbs}"
  base_domain               = "${var.tectonic_base_domain}"
  cluster_id                = "${var.tectonic_cluster_id}"
  cluster_name              = "${var.tectonic_cluster_name}"
  console_elb_dns_name      = "${module.vpc.aws_console_dns_name}"
  console_elb_zone_id       = "${module.vpc.aws_elb_console_zone_id}"
  elb_alias_enabled         = true
  master_count              = "${var.tectonic_master_count}"
  private_zone_id           = "${var.tectonic_aws_external_private_zone != "" ? var.tectonic_aws_external_private_zone : join("", aws_route53_zone.tectonic_int.*.zone_id)}"
  external_vpc_id           = "${module.vpc.vpc_id}"
  extra_tags                = "${var.tectonic_aws_extra_tags}"
  private_endpoints         = "${local.private_endpoints}"
  public_endpoints          = "${local.public_endpoints}"
}
