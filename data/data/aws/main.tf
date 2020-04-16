locals {
  tags = merge(
    {
      "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    },
    var.aws_extra_tags,
  )
}

provider "aws" {
  region = var.aws_region

  skip_region_validation = var.aws_skip_region_validation

  endpoints {
    ec2     = lookup(var.custom_endpoints, "ec2", null)
    elb     = lookup(var.custom_endpoints, "elasticloadbalancing", null)
    iam     = lookup(var.custom_endpoints, "iam", null)
    route53 = lookup(var.custom_endpoints, "route53", null)
    s3      = lookup(var.custom_endpoints, "s3", null)
    sts     = lookup(var.custom_endpoints, "sts", null)
  }
}

module "bootstrap" {
  source = "./bootstrap"

  ami                      = var.aws_region == var.aws_ami_region ? var.aws_ami : aws_ami_copy.imported[0].id
  instance_type            = var.aws_bootstrap_instance_type
  cluster_id               = var.cluster_id
  ignition                 = var.ignition_bootstrap
  subnet_id                = var.aws_publish_strategy == "External" ? module.vpc.az_to_public_subnet_id[var.aws_master_availability_zones[0]] : module.vpc.az_to_private_subnet_id[var.aws_master_availability_zones[0]]
  target_group_arns        = module.vpc.aws_lb_target_group_arns
  target_group_arns_length = module.vpc.aws_lb_target_group_arns_length
  vpc_id                   = module.vpc.vpc_id
  vpc_cidrs                = var.machine_v4_cidrs
  vpc_security_group_ids   = [module.vpc.master_sg_id]
  volume_kms_key_id        = var.aws_master_root_volume_kms_key_id
  publish_strategy         = var.aws_publish_strategy

  tags = local.tags
}

module "masters" {
  source = "./master"

  cluster_id    = var.cluster_id
  instance_type = var.aws_master_instance_type

  tags = local.tags

  availability_zones       = var.aws_master_availability_zones
  az_to_subnet_id          = module.vpc.az_to_private_subnet_id
  instance_count           = var.master_count
  master_sg_ids            = [module.vpc.master_sg_id]
  root_volume_iops         = var.aws_master_root_volume_iops
  root_volume_size         = var.aws_master_root_volume_size
  root_volume_type         = var.aws_master_root_volume_type
  root_volume_encrypted    = var.aws_master_root_volume_encrypted
  root_volume_kms_key_id   = var.aws_master_root_volume_kms_key_id
  target_group_arns        = module.vpc.aws_lb_target_group_arns
  target_group_arns_length = module.vpc.aws_lb_target_group_arns_length
  ec2_ami                  = var.aws_region == var.aws_ami_region ? var.aws_ami : aws_ami_copy.imported[0].id
  user_data_ign            = var.ignition_master
  publish_strategy         = var.aws_publish_strategy
}

module "iam" {
  source = "./iam"

  cluster_id = var.cluster_id

  tags = local.tags
}

module "dns" {
  source = "./route53"

  api_external_lb_dns_name = module.vpc.aws_lb_api_external_dns_name
  api_external_lb_zone_id  = module.vpc.aws_lb_api_external_zone_id
  api_internal_lb_dns_name = module.vpc.aws_lb_api_internal_dns_name
  api_internal_lb_zone_id  = module.vpc.aws_lb_api_internal_zone_id
  base_domain              = var.base_domain
  cluster_domain           = var.cluster_domain
  cluster_id               = var.cluster_id
  tags                     = local.tags
  vpc_id                   = module.vpc.vpc_id
  publish_strategy         = var.aws_publish_strategy
}

module "vpc" {
  source = "./vpc"

  cidr_blocks      = var.machine_v4_cidrs
  cluster_id       = var.cluster_id
  region           = var.aws_region
  vpc              = var.aws_vpc
  public_subnets   = var.aws_public_subnets
  private_subnets  = var.aws_private_subnets
  publish_strategy = var.aws_publish_strategy

  availability_zones = distinct(
    concat(
      var.aws_master_availability_zones,
      var.aws_worker_availability_zones,
    ),
  )

  tags = local.tags
}

resource "aws_ami_copy" "imported" {
  count             = var.aws_region != var.aws_ami_region ? 1 : 0
  name              = "${var.cluster_id}-master"
  source_ami_id     = var.aws_ami
  source_ami_region = var.aws_ami_region
  encrypted         = true

  tags = merge(
    {
      "Name"         = "${var.cluster_id}-ami-${var.aws_region}"
      "sourceAMI"    = var.aws_ami
      "sourceRegion" = var.aws_ami_region
    },
    local.tags,
  )
}

