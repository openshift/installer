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

data "template_file" "etcd_hostname_list" {
  count    = "${var.tectonic_etcd_count > 0 ? var.tectonic_etcd_count : length(data.aws_availability_zones.azs.names) == 5 ? 5 : 3}"
  template = "${var.tectonic_cluster_name}-etcd-${count.index}.${var.tectonic_base_domain}"
}

module "container_linux" {
  source = "../../modules/container_linux"

  release_channel = "${var.tectonic_container_linux_channel}"
  release_version = "${var.tectonic_container_linux_version}"
}

module "vpc" {
  source     = "../../modules/aws/vpc"
  depends_on = ["${aws_s3_bucket_object.tectonic_assets.id}"]

  base_domain     = "${var.tectonic_base_domain}"
  cidr_block      = "${var.tectonic_aws_vpc_cidr_block}"
  cluster_id      = "${var.tectonic_cluster_id}"
  cluster_name    = "${var.tectonic_cluster_name}"
  custom_dns_name = "${var.tectonic_dns_name}"
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

module "masters" {
  source = "../../modules/aws/master-asg"

  autoscaling_group_extra_tags = "${var.tectonic_autoscaling_group_extra_tags}"
  aws_lbs                      = "${module.vpc.aws_lbs}"
  base_domain                  = "${var.tectonic_base_domain}"
  cluster_id                   = "${var.tectonic_cluster_id}"
  cluster_name                 = "${var.tectonic_cluster_name}"
  container_images             = "${var.tectonic_container_images}"
  container_linux_channel      = "${var.tectonic_container_linux_channel}"
  container_linux_version      = "${module.container_linux.version}"
  ec2_type                     = "${var.tectonic_aws_master_ec2_type}"
  extra_tags                   = "${var.tectonic_aws_extra_tags}"
  instance_count               = "1"
  master_iam_role              = "${var.tectonic_aws_master_iam_role_name}"
  master_sg_ids                = "${concat(var.tectonic_aws_master_extra_sg_ids, list(module.vpc.master_sg_id))}"
  private_endpoints            = "${var.tectonic_aws_private_endpoints}"
  public_endpoints             = "${var.tectonic_aws_public_endpoints}"
  root_volume_iops             = "${var.tectonic_aws_master_root_volume_iops}"
  root_volume_size             = "${var.tectonic_aws_master_root_volume_size}"
  root_volume_type             = "${var.tectonic_aws_master_root_volume_type}"
  ssh_key                      = "${var.tectonic_aws_ssh_key}"
  subnet_ids                   = "${module.vpc.master_subnet_ids}"
  ec2_ami                      = "${var.tectonic_aws_ec2_ami_override}"
  kubeconfig_content           = "${local.kubeconfig_kubelet_content}"
}

module "workers" {
  source = "../../modules/aws/worker-asg"

  autoscaling_group_extra_tags = "${var.tectonic_autoscaling_group_extra_tags}"
  cluster_id                   = "${var.tectonic_cluster_id}"
  cluster_name                 = "${var.tectonic_cluster_name}"
  container_linux_channel      = "${var.tectonic_container_linux_channel}"
  container_linux_version      = "${module.container_linux.version}"
  ec2_type                     = "${var.tectonic_aws_worker_ec2_type}"
  extra_tags                   = "${var.tectonic_aws_extra_tags}"
  instance_count               = "0"
  load_balancers               = "${var.tectonic_aws_worker_load_balancers}"
  root_volume_iops             = "${var.tectonic_aws_worker_root_volume_iops}"
  root_volume_size             = "${var.tectonic_aws_worker_root_volume_size}"
  root_volume_type             = "${var.tectonic_aws_worker_root_volume_type}"
  sg_ids                       = "${concat(var.tectonic_aws_worker_extra_sg_ids, list(module.vpc.worker_sg_id))}"
  ssh_key                      = "${var.tectonic_aws_ssh_key}"
  subnet_ids                   = "${module.vpc.worker_subnet_ids}"
  vpc_id                       = "${module.vpc.vpc_id}"
  worker_iam_role              = "${var.tectonic_aws_worker_iam_role_name}"
  ec2_ami                      = "${var.tectonic_aws_ec2_ami_override}"
  base_domain                  = "${var.tectonic_base_domain}"
  kubeconfig_content           = "${local.kubeconfig_kubelet_content}"
}

module "dns" {
  source = "../../modules/dns/route53"

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
  custom_dns_name                = "${var.tectonic_dns_name}"
  elb_alias_enabled              = true
  etcd_count                     = "${length(data.template_file.etcd_hostname_list.*.id)}"
  external_endpoints             = ["${compact(var.tectonic_etcd_servers)}"]
  master_count                   = "${var.tectonic_master_count}"
  tectonic_external_private_zone = "${join("", aws_route53_zone.tectonic_int.*.zone_id)}"
  tectonic_external_vpc_id       = "${module.vpc.vpc_id}"
  tectonic_extra_tags            = "${var.tectonic_aws_extra_tags}"
  tectonic_private_endpoints     = "${var.tectonic_aws_private_endpoints}"
  tectonic_public_endpoints      = "${var.tectonic_aws_public_endpoints}"
}
