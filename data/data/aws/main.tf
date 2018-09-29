locals {
  private_endpoints = "${var.tectonic_aws_endpoints == "public" ? false : true}"
  public_endpoints  = "${var.tectonic_aws_endpoints == "private" ? false : true}"
  private_zone_id   = "${var.tectonic_aws_external_private_zone != "" ? var.tectonic_aws_external_private_zone : join("", aws_route53_zone.tectonic_int.*.zone_id)}"
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

module "bootstrap" {
  source = "./bootstrap"

  ami                         = "${var.tectonic_aws_ec2_ami_override}"
  associate_public_ip_address = "${var.tectonic_aws_endpoints != "private"}"
  bucket                      = "${aws_s3_bucket.tectonic.bucket}"
  cluster_name                = "${var.tectonic_cluster_name}"
  elbs                        = "${module.vpc.aws_lbs}"
  elbs_length                 = "${module.vpc.aws_lbs_length}"
  iam_role                    = "${var.tectonic_aws_master_iam_role_name}"
  ignition                    = "${var.ignition_bootstrap}"
  subnet_id                   = "${module.vpc.master_subnet_ids[0]}"
  vpc_security_group_ids      = ["${concat(var.tectonic_aws_master_extra_sg_ids, list(module.vpc.master_sg_id))}"]

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-bootstrap",
      "tectonicClusterID", "${var.tectonic_cluster_id}"
    ), var.tectonic_aws_extra_tags)}"
}

module "masters" {
  source = "./master"

  elb_api_internal_id = "${module.vpc.aws_elb_api_internal_id}"
  elb_api_external_id = "${module.vpc.aws_elb_api_external_id}"
  elb_console_id      = "${module.vpc.aws_elb_console_id}"
  base_domain         = "${var.tectonic_base_domain}"
  cluster_id          = "${var.tectonic_cluster_id}"
  cluster_name        = "${var.tectonic_cluster_name}"
  ec2_type            = "${var.tectonic_aws_master_ec2_type}"
  extra_tags          = "${var.tectonic_aws_extra_tags}"
  instance_count      = "${var.tectonic_master_count}"
  master_iam_role     = "${var.tectonic_aws_master_iam_role_name}"
  master_sg_ids       = "${concat(var.tectonic_aws_master_extra_sg_ids, list(module.vpc.master_sg_id))}"
  private_endpoints   = "${local.private_endpoints}"
  public_endpoints    = "${local.public_endpoints}"
  root_volume_iops    = "${var.tectonic_aws_master_root_volume_iops}"
  root_volume_size    = "${var.tectonic_aws_master_root_volume_size}"
  root_volume_type    = "${var.tectonic_aws_master_root_volume_type}"
  subnet_ids          = "${module.vpc.master_subnet_ids}"
  ec2_ami             = "${var.tectonic_aws_ec2_ami_override}"
  user_data_igns      = ["${var.ignition_masters}"]
}

module "iam" {
  source = "./iam"

  cluster_name    = "${var.tectonic_cluster_name}"
  worker_iam_role = "${var.tectonic_aws_worker_iam_role_name}"
}

module "dns" {
  source = "./route53"

  api_external_elb_dns_name = "${module.vpc.aws_elb_api_external_dns_name}"
  api_external_elb_zone_id  = "${module.vpc.aws_elb_api_external_zone_id}"
  api_internal_elb_dns_name = "${module.vpc.aws_elb_api_internal_dns_name}"
  api_internal_elb_zone_id  = "${module.vpc.aws_elb_api_internal_zone_id}"
  api_ip_addresses          = "${module.vpc.aws_lbs}"
  base_domain               = "${var.tectonic_base_domain}"
  cluster_id                = "${var.tectonic_cluster_id}"
  cluster_name              = "${var.tectonic_cluster_name}"
  console_elb_dns_name      = "${module.vpc.aws_console_dns_name}"
  console_elb_zone_id       = "${module.vpc.aws_elb_console_zone_id}"
  elb_alias_enabled         = true
  master_count              = "${var.tectonic_master_count}"
  private_zone_id           = "${local.private_zone_id}"
  external_vpc_id           = "${module.vpc.vpc_id}"
  extra_tags                = "${var.tectonic_aws_extra_tags}"
  private_endpoints         = "${local.private_endpoints}"
  public_endpoints          = "${local.public_endpoints}"
}

module "vpc" {
  source = "./vpc"

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

resource "aws_route53_record" "etcd_a_nodes" {
  count   = "${var.tectonic_master_count}"
  type    = "A"
  ttl     = "60"
  zone_id = "${local.private_zone_id}"
  name    = "${var.tectonic_cluster_name}-etcd-${count.index}"
  records = ["${module.masters.ip_addresses[count.index]}"]
}

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

resource "aws_s3_bucket" "tectonic" {
  bucket = "${lower(var.tectonic_cluster_name)}.${var.tectonic_base_domain}"

  acl = "private"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-tectonic",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${var.tectonic_cluster_id}"
    ), var.tectonic_aws_extra_tags)}"

  lifecycle {
    ignore_changes = ["*"]
  }
}
