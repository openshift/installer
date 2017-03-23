resource "aws_route_table" "private_routes" {
  count  = "${var.external_vpc_id == "" ? var.az_count : 0}"
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags {
    Name              = "private-${data.aws_availability_zones.azs.names[count.index]}"
    KubernetesCluster = "${var.cluster_name}"
  }
}

resource "aws_route" "to_nat_gw" {
  count                  = "${var.external_vpc_id == "" ? var.az_count : 0}"
  route_table_id         = "${aws_route_table.private_routes.*.id[count.index]}"
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = "${aws_nat_gateway.nat_gw.*.id[count.index]}"
  depends_on             = ["aws_route_table.private_routes"]
}

resource "aws_subnet" "worker_subnet" {
  count             = "${var.external_vpc_id == "" ? var.az_count : 0}"
  cidr_block        = "${cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 4, count.index + var.az_count)}"
  vpc_id            = "${data.aws_vpc.cluster_vpc.id}"
  availability_zone = "${data.aws_availability_zones.azs.names[count.index]}"

  tags {
    Name                              = "worker-${data.aws_availability_zones.azs.names[count.index]}"
    KubernetesCluster                 = "${var.cluster_name}"
    "kubernetes.io/role/internal-elb" = ""
  }
}

resource "aws_route_table_association" "worker_routing" {
  count          = "${var.external_vpc_id == "" ? var.az_count : 0}"
  route_table_id = "${aws_route_table.private_routes.*.id[count.index]}"
  subnet_id      = "${aws_subnet.worker_subnet.*.id[count.index]}"
}
