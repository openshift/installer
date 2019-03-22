# resource "aws_route_table" "private_routes" {
#   count  = "${local.new_az_count}"
#   vpc_id = "${data.aws_vpc.cluster_vpc.id}"

#   tags = "${merge(map(
#     "Name","${var.cluster_id}-private-${local.new_subnet_azs[count.index]}",
#   ), var.tags)}"
# }

# resource "aws_route" "to_nat_gw" {
#   count                  = "${local.new_az_count}"
#   route_table_id         = "${aws_route_table.private_routes.*.id[count.index]}"
#   destination_cidr_block = "0.0.0.0/0"
#   nat_gateway_id         = "${element(aws_nat_gateway.nat_gw.*.id, count.index)}"
#   depends_on             = ["aws_route_table.private_routes"]
# }

# resource "aws_subnet" "private_subnet" {
#   count = "${local.new_az_count}"

#   vpc_id = "${data.aws_vpc.cluster_vpc.id}"

#   cidr_block = "${cidrsubnet(local.new_private_cidr_range, 3, count.index)}"

#   availability_zone = "${local.new_subnet_azs[count.index]}"

#   tags = "${merge(map(
#     "Name", "${var.cluster_id}-private-${local.new_subnet_azs[count.index]}",
#     "kubernetes.io/role/internal-elb", "",
#   ), var.tags)}"
# }

# resource "aws_route_table_association" "private_routing" {
#   count          = "${local.new_az_count}"
#   route_table_id = "${aws_route_table.private_routes.*.id[count.index]}"
#   subnet_id      = "${aws_subnet.private_subnet.*.id[count.index]}"
# }
