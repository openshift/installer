data "aws_availability_zones" "azs" {}

resource "aws_vpc" "new_vpc" {
  count                = "${length(var.external_vpc_id) > 0 ? 0 : 1}"
  cidr_block           = "${var.vpc_cid_block}"
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_default_route_table" "default" {
  count                  = "${length(var.external_vpc_id) > 0 ? 0 : 1}"
  default_route_table_id = "${aws_vpc.new_vpc.default_route_table_id}"
}

resource "aws_subnet" "az_subnet" {
  count             = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  cidr_block        = "${cidrsubnet(aws_vpc.new_vpc.cidr_block, 8, count.index + 1)}"
  vpc_id            = "${aws_vpc.new_vpc.id}"
  availability_zone = "${data.aws_availability_zones.azs.names[count.index]}"
}

resource "aws_route_table_association" "route_net" {
  count          = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  route_table_id = "${aws_default_route_table.default.id}"
  subnet_id      = "${aws_subnet.az_subnet.*.id[count.index]}"
}

resource "aws_internet_gateway" "igw" {
  count  = "${length(var.external_vpc_id) > 0 ? 0 : 1}"
  vpc_id = "${aws_vpc.new_vpc.id}"
}

resource "aws_route" "igw_route" {
  count                  = "${length(var.external_vpc_id) > 0 ? 0 : 1}"
  destination_cidr_block = "0.0.0.0/0"
  route_table_id         = "${aws_default_route_table.default.id}"
  gateway_id             = "${aws_internet_gateway.igw.id}"
}
