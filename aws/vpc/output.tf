output "cluster_vpc_id" {
  value = "${length(var.external_vpc_id) > 0 ? var.external_vpc_id : aws_vpc.new_vpc.id}"
}
