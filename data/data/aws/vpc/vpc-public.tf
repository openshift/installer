locals {
  has_wavelength_zones = contains(values(var.aws_edge_zones_type), "wavelength-zone")
}

resource "aws_internet_gateway" "igw" {
  count = var.aws_vpc == null ? 1 : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-igw"
    },
    local.tags,
  )
}

resource "aws_route_table" "default" {
  count = var.aws_vpc == null ? 1 : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-public"
    },
    local.tags,
  )
}

resource "aws_main_route_table_association" "main_vpc_routes" {
  count = var.aws_vpc == null ? 1 : 0

  vpc_id         = data.aws_vpc.cluster_vpc.id
  route_table_id = aws_route_table.default[0].id
}

resource "aws_route" "igw_route" {
  count = var.aws_vpc == null ? 1 : 0

  destination_cidr_block = "0.0.0.0/0"
  route_table_id         = aws_route_table.default[0].id
  gateway_id             = aws_internet_gateway.igw[0].id

  timeouts {
    create = "20m"
  }
}

resource "aws_subnet" "public_subnet" {
  count = var.aws_public_subnets == null ? length(local.availability_zones) : 0

  vpc_id            = data.aws_vpc.cluster_vpc.id
  cidr_block        = cidrsubnet(local.new_public_cidr_range, ceil(log(length(local.availability_zones), 2)), count.index)
  availability_zone = local.availability_zones[count.index]

  tags = merge(
    {
      "Name" = "${var.cluster_id}-public-${local.availability_zones[count.index]}"
    },
    local.tags,
  )
}

resource "aws_subnet" "edge_public_subnet" {
  count = local.edge_zones == null ? 0 : length(local.edge_zones)

  vpc_id            = data.aws_vpc.cluster_vpc.id
  cidr_block        = cidrsubnet(local.new_edge_public_cidr_range, ceil(log(length(local.edge_zones), 2)), count.index)
  availability_zone = local.edge_zones[count.index]

  tags = merge(
    {
      "Name" = "${var.cluster_id}-public-${local.edge_zones[count.index]}"
    },
    local.tags,
  )
}

resource "aws_route_table_association" "route_net" {
  count = var.aws_public_subnets == null ? length(local.availability_zones) : 0

  route_table_id = aws_route_table.default[0].id
  subnet_id      = aws_subnet.public_subnet[count.index].id
}

resource "aws_route_table_association" "edge_public_routing" {
  count = local.edge_zones == null ? 0 : length(local.edge_zones)

  route_table_id = lookup(var.aws_edge_zones_type, aws_subnet.edge_public_subnet[count.index].availability_zone, "") == "wavelength-zone" ? aws_route_table.carrier[0].id : aws_route_table.default[0].id
  subnet_id      = aws_subnet.edge_public_subnet[count.index].id
}

resource "aws_eip" "nat_eip" {
  count = var.aws_public_subnets == null ? length(local.availability_zones) : 0
  vpc   = true

  tags = merge(
    {
      "Name" = "${var.cluster_id}-eip-${local.availability_zones[count.index]}"
    },
    local.tags,
  )

  # Terraform does not declare an explicit dependency towards the internet gateway.
  # this can cause the internet gateway to be deleted/detached before the EIPs.
  # https://github.com/coreos/tectonic-installer/issues/1017#issuecomment-307780549
  depends_on = [aws_internet_gateway.igw]
}

resource "aws_nat_gateway" "nat_gw" {
  count = var.aws_public_subnets == null ? length(local.availability_zones) : 0

  allocation_id = aws_eip.nat_eip[count.index].id
  subnet_id     = aws_subnet.public_subnet[count.index].id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-nat-${local.availability_zones[count.index]}"
    },
    local.tags,
  )

  # https://issues.redhat.com/browse/OCPBUGS-891
  depends_on = [aws_eip.nat_eip, aws_subnet.public_subnet]
}

// Carrier Gateway for Wavelength Zones

resource "aws_ec2_carrier_gateway" "carrier" {
  count = local.has_wavelength_zones ? 1 : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-cagw"
    },
    local.tags,
  )
}

resource "aws_route_table" "carrier" {
  count = local.has_wavelength_zones ? 1 : 0

  vpc_id = data.aws_vpc.cluster_vpc.id

  tags = merge(
    {
      "Name" = "${var.cluster_id}-public-carrier"
    },
    local.tags,
  )
}

resource "aws_route" "carrier_default_route" {
  count = local.has_wavelength_zones ? 1 : 0

  destination_cidr_block = "0.0.0.0/0"
  route_table_id         = aws_route_table.carrier[0].id
  carrier_gateway_id     = aws_ec2_carrier_gateway.carrier[0].id

  timeouts {
    create = "20m"
  }
}
