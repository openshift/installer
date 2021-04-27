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

resource "aws_subnet" "private_subnet" {
  count = var.private_subnets == null ? length(var.availability_zones) : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  cidr_block = cidrsubnet(local.new_private_cidr_range, ceil(log(length(var.availability_zones), 2)), count.index)

  availability_zone = var.availability_zones[count.index]

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
