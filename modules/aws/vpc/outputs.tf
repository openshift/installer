output "vpc_id" {
  value = "${length(var.tectonic_aws_external_vpc_id) > 0 ? var.tectonic_aws_external_vpc_id : aws_vpc.new_vpc.id}"
}

output "cluster_default_sg" {
  value = "${aws_security_group.cluster_default.id}"
}

# We have to do this join() & split() 'trick' because null_data_source and 
# the ternary operator can't output lists or maps
#
output "master_subnet_ids" {
  value = ["${split(",", var.tectonic_aws_external_vpc_id == "" ? join(",", aws_subnet.master_subnet.*.id) :  join(",", data.aws_subnet.external_master.*.id))}"]
}

output "worker_subnet_ids" {
  value = ["${split(",", var.tectonic_aws_external_vpc_id == "" ? join(",", aws_subnet.worker_subnet.*.id) :  join(",", data.aws_subnet.external_worker.*.id))}"]
}
