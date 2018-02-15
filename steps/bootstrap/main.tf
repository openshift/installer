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
  source                   = "../../modules/aws/vpc"
  depends_on               = ["${aws_s3_bucket_object.tectonic_assets.id}"]
  base_domain              = "${var.tectonic_base_domain}"
  cidr_block               = "${var.tectonic_aws_vpc_cidr_block}"
  cluster_id               = "${local.cluster_id}"
  cluster_name             = "${var.tectonic_cluster_name}"
  custom_dns_name          = "${var.tectonic_dns_name}"
  enable_etcd_sg           = "${length(compact(var.tectonic_etcd_servers)) == 0 ? 1 : 0}"
  external_master_subnets  = "${compact(var.tectonic_aws_external_master_subnet_ids)}"
  external_vpc_id          = "${var.tectonic_aws_external_vpc_id}"
  external_worker_subnets  = "${compact(var.tectonic_aws_external_worker_subnet_ids)}"
  extra_tags               = "${var.tectonic_aws_extra_tags}"
  private_master_endpoints = "${var.tectonic_aws_private_endpoints}"
  public_master_endpoints  = "${var.tectonic_aws_public_endpoints}"

  # VPC layout settings.
  #
  # The following parameters control the layout of the VPC accross availability zones.
  # Two modes are available:
  # A. Explicitly configure a list of AZs + associated subnet CIDRs
  # B. Let the module calculate subnets accross a set number of AZs
  #
  # To enable mode A, configure a set of AZs + CIDRs for masters and workers using the
  # "tectonic_aws_master_custom_subnets" and "tectonic_aws_worker_custom_subnets" variables.
  #
  # To enable mode B, make sure that "tectonic_aws_master_custom_subnets" and "tectonic_aws_worker_custom_subnets"
  # ARE NOT SET.

  # These counts could be deducted by length(keys(var.tectonic_aws_master_custom_subnets))
  # but there is a restriction on passing computed values as counts. This approach works around that.
  master_az_count = "${length(keys(var.tectonic_aws_master_custom_subnets)) > 0 ? "${length(keys(var.tectonic_aws_master_custom_subnets))}" : "${length(data.aws_availability_zones.azs.names)}"}"
  worker_az_count = "${length(keys(var.tectonic_aws_worker_custom_subnets)) > 0 ? "${length(keys(var.tectonic_aws_worker_custom_subnets))}" : "${length(data.aws_availability_zones.azs.names)}"}"
  # The appending of the "padding" element is required as workaround since the function
  # element() won't work on empty lists. See https://github.com/hashicorp/terraform/issues/11210
  master_subnets = "${concat(values(var.tectonic_aws_master_custom_subnets),list("padding"))}"
  worker_subnets = "${concat(values(var.tectonic_aws_worker_custom_subnets),list("padding"))}"
  # The split() / join() trick works around the limitation of ternary operator expressions
  # only being able to return strings.
  master_azs = "${ split("|", "${length(keys(var.tectonic_aws_master_custom_subnets))}" > 0 ?
    join("|", keys(var.tectonic_aws_master_custom_subnets)) :
    join("|", data.aws_availability_zones.azs.names)
  )}"
  worker_azs = "${ split("|", "${length(keys(var.tectonic_aws_worker_custom_subnets))}" > 0 ?
    join("|", keys(var.tectonic_aws_worker_custom_subnets)) :
    join("|", data.aws_availability_zones.azs.names)
  )}"
}

module "etcd" {
  source = "../../modules/aws/etcd"

  base_domain             = "${var.tectonic_base_domain}"
  cluster_id              = "${local.cluster_id}"
  cluster_name            = "${var.tectonic_cluster_name}"
  container_image         = "${var.tectonic_container_images["etcd"]}"
  container_linux_channel = "${var.tectonic_container_linux_channel}"
  container_linux_version = "${module.container_linux.version}"
  ec2_type                = "${var.tectonic_aws_etcd_ec2_type}"
  external_endpoints      = "${compact(var.tectonic_etcd_servers)}"
  extra_tags              = "${var.tectonic_aws_extra_tags}"
  instance_count          = "${length(data.template_file.etcd_hostname_list.*.id)}"
  root_volume_iops        = "${var.tectonic_aws_etcd_root_volume_iops}"
  root_volume_size        = "${var.tectonic_aws_etcd_root_volume_size}"
  root_volume_type        = "${var.tectonic_aws_etcd_root_volume_type}"
  s3_bucket               = "${aws_s3_bucket.tectonic.bucket}"
  sg_ids                  = "${concat(var.tectonic_aws_etcd_extra_sg_ids, list(module.vpc.etcd_sg_id))}"
  ssh_key                 = "${var.tectonic_aws_ssh_key}"
  subnets                 = "${module.vpc.worker_subnet_ids}"
  etcd_iam_role           = "${var.tectonic_aws_etcd_iam_role_name}"
  ec2_ami                 = "${var.tectonic_aws_ec2_ami_override}"
}

module "masters" {
  source = "../../modules/aws/master-asg"

  autoscaling_group_extra_tags = "${var.tectonic_autoscaling_group_extra_tags}"
  aws_lbs                      = "${module.vpc.aws_lbs}"
  base_domain                  = "${var.tectonic_base_domain}"
  cluster_id                   = "${local.cluster_id}"
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
  cluster_id                   = "${local.cluster_id}"
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
  cluster_id                     = "${local.cluster_id}"
  cluster_name                   = "${var.tectonic_cluster_name}"
  console_elb_dns_name           = "${module.vpc.aws_console_dns_name}"
  console_elb_zone_id            = "${module.vpc.aws_elb_console_zone_id}"
  custom_dns_name                = "${var.tectonic_dns_name}"
  elb_alias_enabled              = true
  etcd_count                     = "${length(data.template_file.etcd_hostname_list.*.id)}"
  etcd_ip_addresses              = "${module.etcd.ip_addresses}"
  external_endpoints             = ["${compact(var.tectonic_etcd_servers)}"]
  master_count                   = "${var.tectonic_master_count}"
  tectonic_external_private_zone = "${aws_route53_zone.tectonic_int.id}"
  tectonic_external_vpc_id       = "${module.vpc.vpc_id}"
  tectonic_extra_tags            = "${var.tectonic_aws_extra_tags}"
  tectonic_private_endpoints     = "${var.tectonic_aws_private_endpoints}"
  tectonic_public_endpoints      = "${var.tectonic_aws_public_endpoints}"
}
