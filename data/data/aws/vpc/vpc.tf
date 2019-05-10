locals {
  new_private_cidr_range = cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 1, 1)
  new_public_cidr_range  = cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 1, 0)
}

resource "aws_vpc" "new_vpc" {
  cidr_block           = var.cidr_block
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
  vpc_id       = aws_vpc.new_vpc.id
  service_name = "com.amazonaws.${var.region}.s3"
  # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
  # force an interpolation expression to be interpreted as a list by wrapping it
  # in an extra set of list brackets. That form was supported for compatibilty in
  # v0.11, but is no longer supported in Terraform v0.12.
  #
  # If the expression in the following list itself returns a list, remove the
  # brackets to avoid interpretation as a list of lists. If the expression
  # returns a single list item then leave it as-is and remove this TODO comment.
  route_table_ids = [concat(
    aws_route_table.private_routes.*.id,
    aws_route_table.default.*.id,
  )]
}

resource "aws_vpc_dhcp_options" "main" {
  domain_name         = var.region == "us-east-1" ? "ec2.internal" : format("%s.compute.internal", var.region)
  domain_name_servers = ["AmazonProvidedDNS"]

  tags = var.tags
}

resource "aws_vpc_dhcp_options_association" "main" {
  vpc_id          = aws_vpc.new_vpc.id
  dhcp_options_id = aws_vpc_dhcp_options.main.id
}

