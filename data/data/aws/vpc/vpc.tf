resource "aws_vpc" "new_vpc" {
  count = var.aws_vpc == null ? 1 : 0

  cidr_block           = var.machine_v4_cidrs[0]
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = merge(
    {
      "Name" = "${var.cluster_id}-vpc"
    },
    local.tags,
  )
}

resource "aws_vpc_endpoint" "s3" {
  count = var.aws_vpc == null ? 1 : 0

  vpc_id       = data.aws_vpc.cluster_vpc.id
  service_name = "com.amazonaws.${var.aws_region}.s3"
  route_table_ids = concat(
    aws_route_table.private_routes.*.id,
    aws_route_table.default.*.id,
    aws_route_table.carrier.*.id,
  )

  tags = local.tags
}

resource "aws_vpc_dhcp_options" "main" {
  count = var.aws_vpc == null ? 1 : 0

  domain_name         = var.aws_region == "us-east-1" ? "ec2.internal" : format("%s.compute.internal", var.aws_region)
  domain_name_servers = ["AmazonProvidedDNS"]

  tags = local.tags
}

resource "aws_vpc_dhcp_options_association" "main" {
  count = var.aws_vpc == null ? 1 : 0

  vpc_id          = data.aws_vpc.cluster_vpc.id
  dhcp_options_id = aws_vpc_dhcp_options.main[0].id
}
