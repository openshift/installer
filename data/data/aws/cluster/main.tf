locals {
  tags = merge(
    {
      "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    },
    var.aws_extra_tags,
  )
  description    = "Created By OpenShift Installer"
  new_tag_len    = 8
  sliced_tag_map = length(var.aws_extra_tags) <= local.new_tag_len ? var.aws_extra_tags : { for k in slice(keys(var.aws_extra_tags), 0, local.new_tag_len) : k => var.aws_extra_tags[k] }
  s3_object_tags = merge(
    {
      "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    },
    local.sliced_tag_map,
  )
}

provider "aws" {
  region = var.aws_region

  skip_region_validation = true

  endpoints {
    ec2     = lookup(var.custom_endpoints, "ec2", null)
    elb     = lookup(var.custom_endpoints, "elasticloadbalancing", null)
    iam     = lookup(var.custom_endpoints, "iam", null)
    route53 = lookup(var.custom_endpoints, "route53", null)
    s3      = lookup(var.custom_endpoints, "s3", null)
    sts     = lookup(var.custom_endpoints, "sts", null)
  }
}

module "masters" {
  source = "./master"

  cluster_id    = var.cluster_id
  instance_type = var.aws_master_instance_type

  tags = local.tags

  availability_zones               = var.aws_master_availability_zones
  az_to_subnet_id                  = module.vpc.az_to_private_subnet_id
  instance_count                   = var.master_count
  master_sg_ids                    = concat([module.vpc.master_sg_id], var.aws_master_security_groups)
  root_volume_iops                 = var.aws_master_root_volume_iops
  root_volume_size                 = var.aws_master_root_volume_size
  root_volume_type                 = var.aws_master_root_volume_type
  root_volume_encrypted            = var.aws_master_root_volume_encrypted
  root_volume_kms_key_id           = var.aws_master_root_volume_kms_key_id
  instance_metadata_authentication = var.aws_master_instance_metadata_authentication
  target_group_arns                = module.vpc.aws_lb_target_group_arns
  target_group_arns_length         = module.vpc.aws_lb_target_group_arns_length
  ec2_ami                          = var.aws_region == var.aws_ami_region ? var.aws_ami : aws_ami_copy.imported[0].id
  user_data_ign                    = var.ignition_master
  publish_strategy                 = var.aws_publish_strategy
  iam_role_name                    = var.aws_master_iam_role_name
}

module "iam" {
  source = "./iam"

  cluster_id           = var.cluster_id
  worker_iam_role_name = var.aws_worker_iam_role_name

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
  internal_zone            = var.aws_internal_zone
  internal_zone_role       = var.aws_internal_zone_role
  vpc_id                   = module.vpc.vpc_id
  region                   = var.aws_region
  publish_strategy         = var.aws_publish_strategy
  custom_endpoints         = var.custom_endpoints
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

  availability_zones = sort(
    distinct(
      concat(
        var.aws_master_availability_zones,
        var.aws_worker_availability_zones,
      ),
    )
  )

  edge_zones         = distinct(var.aws_edge_local_zones)
  edge_parent_gw_map = var.aws_edge_parent_zones_index
  edge_zones_type    = var.aws_edge_zones_type

  tags = local.tags
}

resource "aws_ami_copy" "imported" {
  count             = var.aws_region != var.aws_ami_region ? 1 : 0
  name              = "${var.cluster_id}-master"
  description       = local.description
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

resource "aws_s3_bucket" "ignition" {
  count  = var.aws_preserve_bootstrap_ignition ? 1 : 0
  bucket = var.aws_ignition_bucket

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap"
    },
    local.tags,
  )

  lifecycle {
    ignore_changes = all
  }
}

resource "aws_s3_object" "ignition" {
  count  = var.aws_preserve_bootstrap_ignition ? 1 : 0
  bucket = aws_s3_bucket.ignition[0].id
  key    = "bootstrap.ign"
  source = var.ignition_bootstrap_file

  server_side_encryption = "AES256"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap"
    },
    local.s3_object_tags,
  )

  lifecycle {
    ignore_changes = all
  }
}
