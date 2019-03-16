locals {
  tags = "${merge(map(
    "kubernetes.io/cluster/${var.cluster_id}", "owned"
  ), var.aws_extra_tags)}"
}

provider "aws" {
  region = "${var.aws_region}"
}

module "bootstrap" {
  source = "./bootstrap"

  ami                      = "${aws_ami_copy.main.id}"
  instance_type            = "${var.aws_bootstrap_instance_type}"
  cluster_id               = "${var.cluster_id}"
  ignition                 = "${var.ignition_bootstrap}"
  subnet_id                = "${module.vpc.az_to_public_subnet_id[var.aws_master_availability_zones[0]]}"
  target_group_arns        = "${module.vpc.aws_lb_target_group_arns}"
  target_group_arns_length = "${module.vpc.aws_lb_target_group_arns_length}"
  vpc_id                   = "${module.vpc.vpc_id}"
  vpc_security_group_ids   = "${list(module.vpc.master_sg_id, aws_security_group.public_ssh.id)}"

  tags = "${local.tags}"
}

# This SG is only used by the bootstrap node by default. However we leave it around
# even after the bootstrap node in case a someone wants to run worker nodes in the
# public subnets. This SG can also be applied to the existing nodes in the private
# subnets to get connectivity after VPC peering or VPN configuration.
resource "aws_security_group" "public_ssh" {
  vpc_id = "${module.vpc.vpc_id}"

  tags = "${merge(map(
    "Name", "${var.cluster_id}-public-ssh-sg",
  ), local.tags)}"
}

resource "aws_security_group_rule" "ssh" {
  type              = "ingress"
  security_group_id = "${aws_security_group.public_ssh.id}"

  protocol    = "tcp"
  from_port   = 22
  to_port     = 22
  cidr_blocks = ["0.0.0.0/0"]
}

module "masters" {
  source = "./master"

  cluster_id    = "${var.cluster_id}"
  instance_type = "${var.aws_master_instance_type}"

  tags = "${local.tags}"

  availability_zones       = "${var.aws_master_availability_zones}"
  az_to_subnet_id          = "${module.vpc.az_to_private_subnet_id}"
  instance_count           = "${var.master_count}"
  master_sg_ids            = "${list(module.vpc.master_sg_id)}"
  root_volume_iops         = "${var.aws_master_root_volume_iops}"
  root_volume_size         = "${var.aws_master_root_volume_size}"
  root_volume_type         = "${var.aws_master_root_volume_type}"
  target_group_arns        = "${module.vpc.aws_lb_target_group_arns}"
  target_group_arns_length = "${module.vpc.aws_lb_target_group_arns_length}"
  ec2_ami                  = "${aws_ami_copy.main.id}"
  user_data_ign            = "${var.ignition_master}"
}

module "iam" {
  source = "./iam"

  cluster_id = "${var.cluster_id}"

  tags = "${local.tags}"
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
  etcd_count               = "${var.master_count}"
  etcd_ip_addresses        = "${module.masters.ip_addresses}"
  tags                     = "${local.tags}"
  vpc_id                   = "${module.vpc.vpc_id}"
}

module "vpc" {
  source = "./vpc"

  cidr_block = "${var.machine_cidr}"
  cluster_id = "${var.cluster_id}"
  region     = "${var.aws_region}"

  tags = "${local.tags}"
}

resource "aws_ami_copy" "main" {
  name              = "${var.cluster_id}-master"
  source_ami_id     = "${var.aws_ami}"
  source_ami_region = "${var.aws_region}"
  encrypted         = true

  tags = "${merge(map(
    "Name", "${var.cluster_id}-master",
    "sourceAMI", "${var.aws_ami}",
    "sourceRegion", "${var.aws_region}",
  ), local.tags)}"
}
