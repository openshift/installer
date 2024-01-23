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

provider "aws" {
  alias = "private_hosted_zone"

  assume_role {
    role_arn = var.aws_internal_zone_role
  }

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
  az_to_subnet_id                  = var.az_to_private_subnet_id
  instance_count                   = var.master_count
  master_sg_ids                    = concat([var.master_sg_id], var.aws_master_security_groups)
  root_volume_iops                 = var.aws_master_root_volume_iops
  root_volume_size                 = var.aws_master_root_volume_size
  root_volume_type                 = var.aws_master_root_volume_type
  root_volume_encrypted            = var.aws_master_root_volume_encrypted
  root_volume_kms_key_id           = var.aws_master_root_volume_kms_key_id
  instance_metadata_authentication = var.aws_master_instance_metadata_authentication
  target_group_arns                = var.aws_lb_target_group_arns
  target_group_arns_length         = var.aws_lb_target_group_arns_length
  ec2_ami                          = var.ami_id
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
  count = var.aws_user_provisioned_dns ? 0 : 1

  source = "./route53"
  providers = {
    aws                     = aws
    aws.private_hosted_zone = aws.private_hosted_zone
  }

  api_external_lb_dns_name = var.aws_lb_api_external_dns_name
  api_external_lb_zone_id  = var.aws_lb_api_external_zone_id
  api_internal_lb_dns_name = var.aws_lb_api_internal_dns_name
  api_internal_lb_zone_id  = var.aws_lb_api_internal_zone_id
  base_domain              = var.base_domain
  cluster_domain           = var.cluster_domain
  cluster_id               = var.cluster_id
  tags                     = local.tags
  internal_zone            = var.aws_internal_zone
  internal_zone_role       = var.aws_internal_zone_role
  vpc_id                   = var.vpc_id
  region                   = var.aws_region
  publish_strategy         = var.aws_publish_strategy
  custom_endpoints         = var.custom_endpoints
}
