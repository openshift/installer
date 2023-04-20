locals {
  # The VPC CIDR block is split into two main blocks for private and public subnets.
  # Both blocks is sliced depending the number of zones to create subnets on it,
  # using the maximum it can increase bits from the base MachineCIDR (VPC) for each sub-block.
  # When the IPI deployment is created with a single zone, the allow_expansion_zones will be
  # triggered to prevent the expansion to consume all the CIDR blocks, allowing the VPC
  # to be expanded creating new subnets (on regular AZs, Local or other zone types) in Day-2
  # Operations, then creating manifests to run workloads on it.
  # For example: On the deployment with single-zone using MachineCIDR 10.0.0.0/16, the terraform creates:
  # - 2 * subnets /18: 10.0.0.0/18 and 10.0.128.0/18
  # - keep two block ranges for expansion: 10.0.64.0-10.0.127.255/18 and 10.0.192.0-10.0.255.255/18

  allow_expansion_zones  = length(var.availability_zones) == 1 ? 1 : 0
  new_private_cidr_range = cidrsubnet(cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 1, 0), local.allow_expansion_zones, 0)
  new_public_cidr_range  = cidrsubnet(cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 1, 1), local.allow_expansion_zones, 0)
}

resource "aws_vpc" "new_vpc" {
  count = var.vpc == null ? 1 : 0

  cidr_block           = var.cidr_blocks[0]
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = merge(
    {
      "Name" = "${var.cluster_id}-vpc"
    },
    var.tags,
  )
}

resource "aws_vpc_endpoint" "s3" {
  count = var.vpc == null ? 1 : 0

  vpc_id       = data.aws_vpc.cluster_vpc.id
  service_name = "com.amazonaws.${var.region}.s3"
  route_table_ids = concat(
    aws_route_table.private_routes.*.id,
    aws_route_table.default.*.id,
  )

  tags = var.tags
}

resource "aws_vpc_dhcp_options" "main" {
  count = var.vpc == null ? 1 : 0

  domain_name         = var.region == "us-east-1" ? "ec2.internal" : format("%s.compute.internal", var.region)
  domain_name_servers = ["AmazonProvidedDNS"]

  tags = var.tags
}

resource "aws_vpc_dhcp_options_association" "main" {
  count = var.vpc == null ? 1 : 0

  vpc_id          = data.aws_vpc.cluster_vpc.id
  dhcp_options_id = aws_vpc_dhcp_options.main[0].id
}
