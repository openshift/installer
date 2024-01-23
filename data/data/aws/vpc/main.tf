locals {
  tags = merge(
    {
      "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    },
    var.aws_extra_tags,
  )
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
