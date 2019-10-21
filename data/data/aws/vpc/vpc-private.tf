resource "aws_route_table" "private_routes" {
  count = var.private_subnets == null ? length(var.availability_zones) : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-private-${var.availability_zones[count.index]}"
    },
    var.tags,
  )
}

resource "aws_route" "to_nat_gw" {
  count = var.private_subnets == null ? length(var.availability_zones) : 0

  route_table_id         = aws_route_table.private_routes[count.index].id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = element(aws_nat_gateway.nat_gw.*.id, count.index)
  depends_on             = [aws_route_table.private_routes]

  timeouts {
    create = "20m"
  }
}

# We can't target the NAT gw for our "private" IPv6 subnet.  Instead, we target the internet gateway,
# since we want our private IPv6 addresses to be able to talk out to the internet, too.
resource "aws_route" "private_igw_v6" {
  count = var.use_ipv6 == true && var.private_subnets == null ? length(var.availability_zones) : 0

  route_table_id              = aws_route_table.private_routes[count.index].id
  destination_ipv6_cidr_block = "::/0"
  gateway_id                  = aws_internet_gateway.igw[0].id
  depends_on                  = [aws_route_table.private_routes]

  timeouts {
    create = "20m"
  }
}

resource "aws_subnet" "private_subnet" {
  count = var.private_subnets == null ? length(var.availability_zones) : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  cidr_block = cidrsubnet(local.new_private_cidr_range, 3, count.index)

  availability_zone = var.availability_zones[count.index]

  ipv6_cidr_block                 = var.use_ipv6 == true ? cidrsubnet(data.aws_vpc.cluster_vpc.ipv6_cidr_block, 8, count.index) : ""
  assign_ipv6_address_on_creation = var.use_ipv6

  tags = merge(
    {
      "Name"                            = "${var.cluster_id}-private-${var.availability_zones[count.index]}"
      "kubernetes.io/role/internal-elb" = ""
    },
    var.tags,
  )
}

resource "aws_route_table_association" "private_routing" {
  count = var.private_subnets == null ? length(var.availability_zones) : 0

  route_table_id = aws_route_table.private_routes[count.index].id
  subnet_id      = aws_subnet.private_subnet[count.index].id
}
