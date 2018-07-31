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
  source = "../../../modules/aws/master-asg"

  autoscaling_group_extra_tags = "${var.tectonic_autoscaling_group_extra_tags}"
  aws_lbs                      = "${local.aws_lbs}"
  base_domain                  = "${var.tectonic_base_domain}"
  cluster_id                   = "${var.tectonic_cluster_id}"
  cluster_name                 = "${var.tectonic_cluster_name}"
  container_images             = "${var.tectonic_container_images}"
  container_linux_channel      = "${var.tectonic_container_linux_channel}"
  container_linux_version      = "${module.container_linux.version}"
  ec2_type                     = "${var.tectonic_aws_master_ec2_type}"
  extra_tags                   = "${var.tectonic_aws_extra_tags}"
  instance_count               = "${var.tectonic_bootstrap == "true" ? 1 : var.tectonic_master_count}"
  master_iam_role              = "${var.tectonic_aws_master_iam_role_name}"
  master_sg_ids                = "${concat(var.tectonic_aws_master_extra_sg_ids, list(local.sg_id))}"
  private_endpoints            = "${local.private_endpoints}"
  public_endpoints             = "${local.public_endpoints}"
  region                       = "${var.tectonic_aws_region}"
  root_volume_iops             = "${var.tectonic_aws_master_root_volume_iops}"
  root_volume_size             = "${var.tectonic_aws_master_root_volume_size}"
  root_volume_type             = "${var.tectonic_aws_master_root_volume_type}"
  ssh_key                      = "${var.tectonic_aws_ssh_key}"
  subnet_ids                   = "${local.subnet_ids}"
  ec2_ami                      = "${var.tectonic_aws_ec2_ami_override}"
  user_data_ign                = "${file("${path.cwd}/${var.tectonic_ignition_master}")}"
}
