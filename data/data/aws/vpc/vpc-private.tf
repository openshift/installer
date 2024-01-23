resource "aws_route_table" "private_routes" {
  count = var.aws_private_subnets == null ? length(local.availability_zones) : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-private-${local.availability_zones[count.index]}"
    },
    local.tags,
  )
}

resource "aws_route" "to_nat_gw" {
  count = var.aws_private_subnets == null ? length(local.availability_zones) : 0

  route_table_id         = aws_route_table.private_routes[count.index].id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = element(aws_nat_gateway.nat_gw.*.id, count.index)
  depends_on             = [aws_route_table.private_routes]

  timeouts {
    create = "20m"
  }
}

resource "aws_subnet" "private_subnet" {
  count = var.aws_private_subnets == null ? length(local.availability_zones) : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  cidr_block = cidrsubnet(local.new_private_cidr_range, ceil(log(length(local.availability_zones), 2)), count.index)

  availability_zone = local.availability_zones[count.index]

  tags = merge(
    {
      "Name"                            = "${var.cluster_id}-private-${local.availability_zones[count.index]}"
      "kubernetes.io/role/internal-elb" = ""
    },
    local.tags,
  )
}

resource "aws_subnet" "edge_private_subnet" {
  count = local.edge_zones == null ? 0 : length(local.edge_zones)

  vpc_id            = data.aws_vpc.cluster_vpc.id
  cidr_block        = cidrsubnet(local.new_edge_private_cidr_range, ceil(log(length(local.edge_zones), 2)), count.index)
  availability_zone = local.edge_zones[count.index]

  tags = merge(
    {
      "Name" = "${var.cluster_id}-private-${local.edge_zones[count.index]}"
    },
    local.tags,
  )

}

resource "aws_route_table_association" "private_routing" {
  count = var.aws_private_subnets == null ? length(local.availability_zones) : 0

  route_table_id = aws_route_table.private_routes[count.index].id
  subnet_id      = aws_subnet.private_subnet[count.index].id
}

resource "aws_route_table_association" "edge_private_routing" {
  count = local.edge_zones == null ? 0 : length(local.edge_zones)

  # Lookup the index of the parent zone from a given Local Zone name,
  # getting the index for the route table id for that zone (parent),
  # when not found (parent zone's gateway does not exists), the first
  # route table will be used.
  # Example aws_edge_parent_zones_index = {us-east-1-nyc-1a=0}
  route_table_id = aws_route_table.private_routes[lookup(var.aws_edge_parent_zones_index, aws_subnet.edge_private_subnet[count.index].availability_zone, 0)].id
  subnet_id      = aws_subnet.edge_private_subnet[count.index].id
}
