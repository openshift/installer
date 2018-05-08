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

module "container_linux" {
  source = "../../../modules/container_linux"

  release_channel = "${var.tectonic_container_linux_channel}"
  release_version = "${var.tectonic_container_linux_version}"
}

# TNC
resource "aws_route53_zone" "tectonic_int" {
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
  enable_etcd_sg  = "${length(compact(var.tectonic_etcd_servers)) == 0 ? 1 : 0}"
  external_vpc_id = "${var.tectonic_aws_external_vpc_id}"

  external_master_subnet_ids = "${compact(var.tectonic_aws_external_master_subnet_ids)}"
  external_worker_subnet_ids = "${compact(var.tectonic_aws_external_worker_subnet_ids)}"
  extra_tags                 = "${var.tectonic_aws_extra_tags}"

  // empty map subnet_configs will have the vpc module creating subnets in all availabile AZs
  new_master_subnet_configs = "${var.tectonic_aws_master_custom_subnets}"
  new_worker_subnet_configs = "${var.tectonic_aws_worker_custom_subnets}"

  private_master_endpoints = "${var.tectonic_aws_private_endpoints}"
  public_master_endpoints  = "${var.tectonic_aws_public_endpoints}"
}

module "dns" {
  source = "../../../modules/dns/route53"

  api_external_elb_dns_name      = "${module.vpc.aws_api_external_dns_name}"
  api_external_elb_zone_id       = "${module.vpc.aws_elb_api_external_zone_id}"
  api_internal_elb_dns_name      = "${module.vpc.aws_api_internal_dns_name}"
  api_internal_elb_zone_id       = "${module.vpc.aws_elb_api_internal_zone_id}"
  api_ip_addresses               = "${module.vpc.aws_lbs}"
  base_domain                    = "${var.tectonic_base_domain}"
  cluster_id                     = "${var.tectonic_cluster_id}"
  cluster_name                   = "${var.tectonic_cluster_name}"
  console_elb_dns_name           = "${module.vpc.aws_console_dns_name}"
  console_elb_zone_id            = "${module.vpc.aws_elb_console_zone_id}"
  elb_alias_enabled              = true
  master_count                   = "${var.tectonic_master_count}"
  tectonic_external_private_zone = "${join("", aws_route53_zone.tectonic_int.*.zone_id)}"
  tectonic_external_vpc_id       = "${module.vpc.vpc_id}"
  tectonic_extra_tags            = "${var.tectonic_aws_extra_tags}"
  tectonic_private_endpoints     = "${var.tectonic_aws_private_endpoints}"
  tectonic_public_endpoints      = "${var.tectonic_aws_public_endpoints}"
}
