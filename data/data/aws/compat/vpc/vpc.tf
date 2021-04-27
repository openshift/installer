locals {
  new_private_cidr_range = cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 1, 1)
  new_public_cidr_range  = cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 1, 0)
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
