data "aws_internet_gateway" "igw" {
  count = var.vpc == null ? 1 : 0

  tags = merge(
    {
      "Name" = "${var.cluster_id}-igw"
    },
    var.tags,
  )
}

data "aws_route_table" "default" {
  count = var.vpc == null ? 1 : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-public"
    },
    var.tags,
  )
}

data "aws_route" "igw_route" {
  count = var.vpc == null ? 1 : 0

  destination_cidr_block = "0.0.0.0/0"
  route_table_id         = data.aws_route_table.default[0].id
  gateway_id             = data.aws_internet_gateway.igw[0].id
}

data "aws_subnet" "public_subnet" {
  count = var.public_subnets == null ? length(var.availability_zones) : 0

  vpc_id            = data.aws_vpc.cluster_vpc.id
  cidr_block        = cidrsubnet(local.new_public_cidr_range, ceil(log(length(var.availability_zones), 2)), count.index)
  availability_zone = var.availability_zones[count.index]

  tags = merge(
    {
      "Name" = "${var.cluster_id}-public-${var.availability_zones[count.index]}"
    },
    var.tags,
  )
}

data "aws_eip" "nat_eip" {
  count = var.public_subnets == null ? length(var.availability_zones) : 0

  tags = merge(
    {
      "Name" = "${var.cluster_id}-eip-${var.availability_zones[count.index]}"
    },
    var.tags,
  )
}

data "aws_nat_gateway" "nat_gw" {
  count = var.public_subnets == null ? length(var.availability_zones) : 0

  subnet_id = data.aws_subnet.public_subnet[count.index].id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-nat-${var.availability_zones[count.index]}"
    },
    var.tags,
  )
}
