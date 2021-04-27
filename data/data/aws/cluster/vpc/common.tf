# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file

locals {
  public_endpoints = var.publish_strategy == "External" ? true : false
  description      = "Created By OpenShift Installer"
}

# all data sources should be input variable-agnostic and used as canonical source for querying "state of resources" and building outputs
# (ie: we don't want "aws.new_vpc" and "data.aws_vpc.cluster_vpc", just "data.aws_vpc.cluster_vpc" used everwhere).

data "aws_vpc" "cluster_vpc" {
  id = var.vpc == null ? data.aws_vpc.new_vpc[0].id : var.vpc
}

data "aws_subnet" "public" {
  count = var.public_subnets == null ? length(var.availability_zones) : length(var.public_subnets)

  id = var.public_subnets == null ? data.aws_subnet.public_subnet[count.index].id : var.public_subnets[count.index]
}

data "aws_subnet" "private" {
  count = var.private_subnets == null ? length(var.availability_zones) : length(var.private_subnets)

  id = var.private_subnets == null ? data.aws_subnet.private_subnet[count.index].id : var.private_subnets[count.index]
}
