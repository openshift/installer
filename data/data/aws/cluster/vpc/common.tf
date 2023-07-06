# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file

locals {
  public_endpoints = var.publish_strategy == "External" ? true : false
  byo_lbs          = var.aws_api_ext_lb_name == null ? false : true
  description      = "Created By OpenShift Installer"

  # CIDR block distribution:
  # allow_expansion_* flags checks if is a single-zone deployment, if true (1) the
  # available CIDR block will be split into two to allow user expansion.
  allow_expansion_zones = length(var.availability_zones) == 1 ? 1 : 0
  allow_expansion_edge  = length(var.edge_zones) == 1 ? 1 : 0

  # edge_enabled flag is enabled when edge zones (Local Zone) are provided.
  edge_enabled = length(var.edge_zones) > 0 ? 1 : 0

  # CIDR blocks for default IPI installation
  cidr_dedicated_private = cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 1, 0)
  cidr_dedicated_public  = cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 1, 1)

  # CIDR blocks used when creating subnets into edge zones.
  # The Public CIDR is used to create the CIDR blocks for edge subnets.
  cidr_shared_public = cidrsubnet(local.cidr_dedicated_public, 1, 0)
  cidr_shared_edge   = cidrsubnet(local.cidr_dedicated_public, 1, 1)

  # CIDR blocks for edge subnets
  cidr_edge_private = cidrsubnet(local.cidr_shared_edge, 1, 0)
  cidr_edge_public  = cidrsubnet(local.cidr_shared_edge, 1, 1)

  # CIDR blocks pool used to create subnets for each zone
  new_private_cidr_range      = cidrsubnet(local.cidr_dedicated_private, local.allow_expansion_zones, 0)
  new_public_cidr_range       = local.edge_enabled == 0 ? cidrsubnet(local.cidr_dedicated_public, local.allow_expansion_zones, 0) : cidrsubnet(local.cidr_shared_public, local.allow_expansion_zones, 0)
  new_edge_private_cidr_range = local.allow_expansion_edge == 0 ? local.cidr_edge_private : cidrsubnet(local.cidr_edge_private, local.allow_expansion_edge, 0)
  new_edge_public_cidr_range  = local.allow_expansion_edge == 0 ? local.cidr_edge_public : cidrsubnet(local.cidr_edge_public, local.allow_expansion_edge, 0)
}

# all data sources should be input variable-agnostic and used as canonical source for querying "state of resources" and building outputs
# (ie: we don't want "aws.new_vpc" and "data.aws_vpc.cluster_vpc", just "data.aws_vpc.cluster_vpc" used everwhere).

data "aws_vpc" "cluster_vpc" {
  id = var.vpc == null ? aws_vpc.new_vpc[0].id : var.vpc
}

data "aws_subnet" "public" {
  count = var.public_subnets == null ? length(var.availability_zones) : length(var.public_subnets)

  id = var.public_subnets == null ? aws_subnet.public_subnet[count.index].id : var.public_subnets[count.index]
}

data "aws_subnet" "private" {
  count = var.private_subnets == null ? length(var.availability_zones) : length(var.private_subnets)

  id = var.private_subnets == null ? aws_subnet.private_subnet[count.index].id : var.private_subnets[count.index]
}

data "aws_subnet" "edge_private" {
  count = var.edge_zones == null ? 0 : length(var.edge_zones)

  id = var.edge_zones == null ? null : aws_subnet.edge_private_subnet[count.index].id
}

data "aws_subnet" "edge_public" {
  count = var.edge_zones == null ? 0 : length(var.edge_zones)

  id = var.edge_zones == null ? null : aws_subnet.edge_public_subnet[count.index].id
}

data "aws_lb" "ext" {
  name = local.byo_lbs ? var.api_ext_lb_name : aws_lb.api_external[0].name
}

data "aws_lb" "int" {
  name = local.byo_lbs ? var.api_int_lb_name : aws_lb.api_internal[0].name
}
