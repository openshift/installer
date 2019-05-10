resource "aws_internet_gateway" "igw" {
  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-igw"
    },
    var.tags,
  )
}

resource "aws_route_table" "default" {
  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-public"
    },
    var.tags,
  )
}

resource "aws_main_route_table_association" "main_vpc_routes" {
  vpc_id         = data.aws_vpc.cluster_vpc.id
  route_table_id = aws_route_table.default.id
}

resource "aws_route" "igw_route" {
  destination_cidr_block = "0.0.0.0/0"
  route_table_id         = aws_route_table.default.id
  gateway_id             = aws_internet_gateway.igw.id

  timeouts {
    create = "20m"
  }
}

resource "aws_subnet" "public_subnet" {
  count  = local.new_az_count
  vpc_id = data.aws_vpc.cluster_vpc.id

  cidr_block = cidrsubnet(local.new_public_cidr_range, 3, count.index)

  availability_zone = var.availability_zones[count.index]

  tags = merge(
    {
      "Name" = "${var.cluster_id}-public-${var.availability_zones[count.index]}"
    },
    var.tags,
  )
}

resource "aws_route_table_association" "route_net" {
  count          = local.new_az_count
  route_table_id = aws_route_table.default.id
  subnet_id      = aws_subnet.public_subnet[count.index].id
}

resource "aws_eip" "nat_eip" {
  count = local.new_az_count
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
  count         = local.new_az_count
  allocation_id = aws_eip.nat_eip[count.index].id
  subnet_id     = aws_subnet.public_subnet[count.index].id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-nat-${var.availability_zones[count.index]}"
    },
    var.tags,
  )
}

