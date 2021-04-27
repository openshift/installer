locals {
  new_private_cidr_range = cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 1, 1)
  new_public_cidr_range  = cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 1, 0)
}

data "aws_vpc" "new_vpc" {
  count = var.vpc == null ? 1 : 0

  cidr_block = var.cidr_blocks[0]

  tags = merge(
    {
      "Name" = "${var.cluster_id}-vpc"
    },
    var.tags,
  )
}

data "aws_vpc_endpoint" "s3" {
  count = var.vpc == null ? 1 : 0

  vpc_id       = data.aws_vpc.cluster_vpc.id
  service_name = "com.amazonaws.${var.region}.s3"

  tags = var.tags
}

data "aws_vpc_dhcp_options" "main" {
  count = var.vpc == null ? 1 : 0

  filter {
    name   = "tag:Name"
    values = ["${var.cluster_id}-vpc-dhcp-options"]
  }

  tags = var.tags
}
