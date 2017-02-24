resource "aws_route_table" "private_routing" {
  count  = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags {
    Name = "private-${data.aws_availability_zones.azs.names[count.index]}"
  }
}

resource "aws_route" "to_nat_gw" {
  count                  = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  route_table_id         = "${aws_route_table.private_routing.*.id[count.index]}"
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = "${aws_nat_gateway.nat_gw.*.id[count.index]}"
  depends_on             = ["aws_route_table.private_routing"]
}

resource "aws_route_table_association" "private_routing" {
  count          = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  route_table_id = "${aws_route_table.private_routing.*.id[count.index]}"
  subnet_id      = "${aws_subnet.etcd_subnet.*.id[count.index]}"
}

resource "aws_subnet" "etcd_subnet" {
  count             = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  cidr_block        = "${cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 12, count.index + 10)}"
  vpc_id            = "${data.aws_vpc.cluster_vpc.id}"
  availability_zone = "${data.aws_availability_zones.azs.names[count.index]}"

  tags {
    Name = "etcd-${data.aws_availability_zones.azs.names[count.index]}"
  }
}

resource "aws_subnet" "master_subnet" {
  count             = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  cidr_block        = "${cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 12, count.index)}"
  vpc_id            = "${data.aws_vpc.cluster_vpc.id}"
  availability_zone = "${data.aws_availability_zones.azs.names[count.index]}"

  tags {
    Name = "master-${data.aws_availability_zones.azs.names[count.index]}"
  }
}

resource "aws_route_table_association" "master_routing" {
  count          = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  route_table_id = "${aws_route_table.private_routing.*.id[count.index]}"
  subnet_id      = "${aws_subnet.master_subnet.*.id[count.index]}"
}

resource "aws_subnet" "worker_subnet" {
  count             = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  cidr_block        = "${cidrsubnet(data.aws_vpc.cluster_vpc.cidr_block, 12, count.index + 5)}"
  vpc_id            = "${data.aws_vpc.cluster_vpc.id}"
  availability_zone = "${data.aws_availability_zones.azs.names[count.index]}"

  tags {
    Name = "worker-${data.aws_availability_zones.azs.names[count.index]}"
  }
}

resource "aws_route_table_association" "worker_routing" {
  count          = "${length(var.external_vpc_id) > 0 ? 0 : var.az_count}"
  route_table_id = "${aws_route_table.private_routing.*.id[count.index]}"
  subnet_id      = "${aws_subnet.worker_subnet.*.id[count.index]}"
}

resource "aws_security_group" "cluster_default" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  ingress {
    protocol  = -1
    self      = true
    from_port = 0
    to_port   = 0
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self        = true
    cidr_blocks = ["0.0.0.0/0"]
  }
}
