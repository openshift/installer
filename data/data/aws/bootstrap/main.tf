locals {
  tags = merge(
    {
      "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    },
    var.aws_extra_tags,
  )
  description = "Created By OpenShift Installer"
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

  ami                      = var.aws_region == var.aws_ami_region ? var.aws_ami : data.aws_ami.imported[0].id
  instance_type            = var.aws_bootstrap_instance_type
  cluster_id               = var.cluster_id
  ignition                 = var.ignition_bootstrap_file
  ignition_bucket          = var.aws_ignition_bucket
  ignition_stub            = var.aws_bootstrap_stub_ignition
  subnet_id                = var.aws_publish_strategy == "External" ? module.vpc.az_to_public_subnet_id[var.aws_master_availability_zones[0]] : module.vpc.az_to_private_subnet_id[var.aws_master_availability_zones[0]]
  target_group_arns        = module.vpc.aws_lb_target_group_arns
  target_group_arns_length = module.vpc.aws_lb_target_group_arns_length
  vpc_id                   = module.vpc.vpc_id
  vpc_cidrs                = var.machine_v4_cidrs
  vpc_security_group_ids   = [module.vpc.master_sg_id]
  volume_kms_key_id        = var.aws_master_root_volume_kms_key_id
  publish_strategy         = var.aws_publish_strategy
  iam_role_name            = var.aws_master_iam_role_name

  tags = local.tags
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

data "aws_ami" "imported" {
  count = var.aws_region != var.aws_ami_region ? 1 : 0

  owners      = ["self"]
  most_recent = true

  filter {
    name   = "name"
    values = ["${var.cluster_id}-master"]
  }
}
