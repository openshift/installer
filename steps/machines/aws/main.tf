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

module "container_linux" {
  source = "../../../modules/container_linux"

  release_channel = "${var.tectonic_container_linux_channel}"
  release_version = "${var.tectonic_container_linux_version}"
}

module "masters" {
  source = "../../../modules/aws/master"

  aws_lbs                 = "${local.aws_lbs}"
  base_domain             = "${var.tectonic_base_domain}"
  cluster_id              = "${var.tectonic_cluster_id}"
  cluster_name            = "${var.tectonic_cluster_name}"
  container_images        = "${var.tectonic_container_images}"
  container_linux_channel = "${var.tectonic_container_linux_channel}"
  container_linux_version = "${module.container_linux.version}"
  ec2_type                = "${var.tectonic_aws_master_ec2_type}"
  extra_tags              = "${var.tectonic_aws_extra_tags}"
  instance_count          = "${var.tectonic_master_count}"
  master_iam_role         = "${var.tectonic_aws_master_iam_role_name}"
  master_sg_ids           = "${concat(var.tectonic_aws_master_extra_sg_ids, list(local.master_sg_id))}"
  private_endpoints       = "${local.private_endpoints}"
  public_endpoints        = "${local.public_endpoints}"
  region                  = "${var.tectonic_aws_region}"
  root_volume_iops        = "${var.tectonic_aws_master_root_volume_iops}"
  root_volume_size        = "${var.tectonic_aws_master_root_volume_size}"
  root_volume_type        = "${var.tectonic_aws_master_root_volume_type}"
  subnet_ids              = "${local.master_subnet_ids}"
  ec2_ami                 = "${var.tectonic_aws_ec2_ami_override}"
  user_data_igns          = "${var.tectonic_ignition_masters}"
}

module "workers" {
  source = "../../../modules/aws/worker"

  cluster_id              = "${var.tectonic_cluster_id}"
  cluster_name            = "${var.tectonic_cluster_name}"
  container_linux_channel = "${var.tectonic_container_linux_channel}"
  container_linux_version = "${module.container_linux.version}"
  ec2_type                = "${var.tectonic_aws_worker_ec2_type}"
  extra_tags              = "${var.tectonic_aws_extra_tags}"
  instance_count          = "${var.tectonic_worker_count}"
  load_balancers          = "${var.tectonic_aws_worker_load_balancers}"
  region                  = "${var.tectonic_aws_region}"
  root_volume_iops        = "${var.tectonic_aws_worker_root_volume_iops}"
  root_volume_size        = "${var.tectonic_aws_worker_root_volume_size}"
  root_volume_type        = "${var.tectonic_aws_worker_root_volume_type}"
  sg_ids                  = "${concat(var.tectonic_aws_worker_extra_sg_ids, list(local.worker_sg_id))}"
  subnet_ids              = "${local.worker_subnet_ids}"
  worker_iam_role         = "${var.tectonic_aws_worker_iam_role_name}"
  ec2_ami                 = "${var.tectonic_aws_ec2_ami_override}"
  base_domain             = "${var.tectonic_base_domain}"
  user_data_ign           = "${file("${path.cwd}/${var.tectonic_ignition_worker}")}"
}

resource "aws_route53_record" "etcd_a_nodes" {
  count   = "${var.tectonic_master_count}"
  type    = "A"
  ttl     = "60"
  zone_id = "${local.private_zone_id}"
  name    = "${var.tectonic_cluster_name}-etcd-${count.index}"
  records = ["${module.masters.ip_addresses[count.index]}"]
}
