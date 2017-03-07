resource "aws_internet_gateway" "igw" {
  count  = "${length(var.external_vpc_id) > 0 ? 0 : 1}"
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"
}

resource "aws_route_table" "default" {
  count  = "${length(var.external_vpc_id) > 0 ? 0 : 1}"
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags {
    Name = "public"
    KubernetesCluster = "${var.cluster_name}"
  }
}

resource "aws_main_route_table_association" "a" {
  vpc_id         = "${data.aws_vpc.cluster_vpc.id}"
  route_table_id = "${aws_route_table.default.id}"
}

resource "aws_route" "igw_route" {
  count                  = "${length(var.external_vpc_id) > 0 ? 0 : 1}"
  destination_cidr_block = "0.0.0.0/0"
  route_table_id         = "${aws_route_table.default.id}"
  gateway_id             = "${aws_internet_gateway.igw.id}"
}

resource "aws_subnet" "az_subnet_pub" {
  count             = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  cidr_block        = "${cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 8, count.index + 1)}"
  vpc_id            = "${data.aws_vpc.cluster_vpc.id}"
  availability_zone = "${data.aws_availability_zones.azs.names[count.index]}"

  tags {
    Name = "public-${data.aws_availability_zones.azs.names[count.index]}"
    KubernetesCluster = "${var.cluster_name}"
  }
}

resource "aws_route_table_association" "route_net" {
  count          = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  route_table_id = "${aws_route_table.default.id}"
  subnet_id      = "${aws_subnet.az_subnet_pub.*.id[count.index]}"
}

resource "aws_eip" "nat_eip" {
  count = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  vpc   = true
}

resource "aws_nat_gateway" "nat_gw" {
  count         = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  allocation_id = "${aws_eip.nat_eip.*.id[count.index]}"
  subnet_id     = "${aws_subnet.az_subnet_pub.*.id[count.index]}"
}
