resource "aws_internet_gateway" "igw" {
  count = var.vpc == null ? 1 : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-igw"
    },
    var.tags,
  )
}

resource "aws_route_table" "default" {
  count = var.vpc == null ? 1 : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-public"
    },
    var.tags,
  )
}

resource "aws_main_route_table_association" "main_vpc_routes" {
  count = var.vpc == null ? 1 : 0

  vpc_id         = data.aws_vpc.cluster_vpc.id
  route_table_id = aws_route_table.default[0].id
}

resource "aws_route" "igw_route" {
  count = var.vpc == null ? 1 : 0

  destination_cidr_block = "0.0.0.0/0"
  route_table_id         = aws_route_table.default[0].id
  gateway_id             = aws_internet_gateway.igw[0].id

  timeouts {
    create = "20m"
  }
}

resource "aws_route" "igw_route_v6" {
  count = var.use_ipv6 == true && var.vpc == null ? 1 : 0

  destination_ipv6_cidr_block = "::/0"
  route_table_id              = aws_route_table.default[0].id
  gateway_id                  = aws_internet_gateway.igw[0].id

  timeouts {
    create = "20m"
  }
}

resource "aws_subnet" "public_subnet" {
  count = var.public_subnets == null ? length(var.availability_zones) : 0

  vpc_id            = data.aws_vpc.cluster_vpc.id
  cidr_block        = cidrsubnet(local.new_public_cidr_range, 3, count.index)
  availability_zone = var.availability_zones[count.index]

  ipv6_cidr_block                 = var.use_ipv6 == true ? cidrsubnet(data.aws_vpc.cluster_vpc.ipv6_cidr_block, 8, count.index + length(var.availability_zones)) : ""
  assign_ipv6_address_on_creation = var.use_ipv6

  tags = merge(
    {
      "Name" = "${var.cluster_id}-public-${var.availability_zones[count.index]}"
    },
    var.tags,
  )
}

resource "aws_route_table_association" "route_net" {
  count = var.public_subnets == null ? length(var.availability_zones) : 0

  route_table_id = aws_route_table.default[0].id
  subnet_id      = aws_subnet.public_subnet[count.index].id
}

resource "aws_eip" "nat_eip" {
  count = var.public_subnets == null ? length(var.availability_zones) : 0
  vpc   = true

  tags = merge(
    {
      "Name" = "${var.cluster_id}-eip-${var.availability_zones[count.index]}"
    },
    var.tags,
  )

  # Terraform does not declare an explicit dependency towards the internet gateway.
  # this can cause the internet gateway to be deleted/detached before the EIPs.
  # https://github.com/coreos/tectonic-installer/issues/1017#issuecomment-307780549
  depends_on = [aws_internet_gateway.igw]
}

resource "aws_nat_gateway" "nat_gw" {
  count = var.public_subnets == null ? length(var.availability_zones) : 0

  allocation_id = aws_eip.nat_eip[count.index].id
  subnet_id     = aws_subnet.public_subnet[count.index].id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-nat-${var.availability_zones[count.index]}"
    },
    var.tags,
  )
}
