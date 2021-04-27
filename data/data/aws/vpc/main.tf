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
