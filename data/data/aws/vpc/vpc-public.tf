resource "aws_internet_gateway" "igw" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags = "${merge(map(
      "Name", "${var.cluster_name}-igw",
      "kubernetes.io/cluster/${var.cluster_name}", "shared",
      "tectonicClusterID", "${var.cluster_id}",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"
}

resource "aws_route_table" "default" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags = "${merge(map(
      "Name", "${var.cluster_name}-public",
      "kubernetes.io/cluster/${var.cluster_name}", "shared",
      "tectonicClusterID", "${var.cluster_id}",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"
}

resource "aws_main_route_table_association" "main_vpc_routes" {
  vpc_id         = "${data.aws_vpc.cluster_vpc.id}"
  route_table_id = "${aws_route_table.default.id}"
}

resource "aws_route" "igw_route" {
  destination_cidr_block = "0.0.0.0/0"
  route_table_id         = "${aws_route_table.default.id}"
  gateway_id             = "${aws_internet_gateway.igw.id}"
}

resource "aws_subnet" "master_subnet" {
  count  = "${local.new_az_count}"
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  cidr_block = "${cidrsubnet(local.new_master_cidr_range, 3, count.index)}"

  availability_zone = "${local.new_subnet_azs[count.index]}"

  tags = "${merge(map(
    "Name", "${var.cluster_name}-master-${local.new_subnet_azs[count.index]}",
      "kubernetes.io/cluster/${var.cluster_name}", "shared",
      "tectonicClusterID", "${var.cluster_id}",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"
}

resource "aws_route_table_association" "route_net" {
  count          = "${local.new_az_count}"
  route_table_id = "${aws_route_table.default.id}"
  subnet_id      = "${aws_subnet.master_subnet.*.id[count.index]}"
}

resource "aws_eip" "nat_eip" {
  count = "${local.new_az_count}"
  vpc   = true

  tags = "${merge(map(
      "tectonicClusterID", "${var.cluster_id}",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"

  # Terraform does not declare an explicit dependency towards the internet gateway.
  # this can cause the internet gateway to be deleted/detached before the EIPs.
  # https://github.com/coreos/tectonic-installer/issues/1017#issuecomment-307780549
  depends_on = ["aws_internet_gateway.igw"]
}

resource "aws_nat_gateway" "nat_gw" {
  count         = "${local.new_az_count}"
  allocation_id = "${aws_eip.nat_eip.*.id[count.index]}"
  subnet_id     = "${aws_subnet.master_subnet.*.id[count.index]}"

  tags = "${merge(map(
      "tectonicClusterID", "${var.cluster_id}",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"
}
