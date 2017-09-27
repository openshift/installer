output "vpc_id" {
  value = "${data.aws_vpc.cluster_vpc.id}"
}

# We have to do this join() & split() 'trick' because null_data_source and
# the ternary operator can't output lists or maps
output "master_subnet_ids" {
  value = ["${split(",", var.external_vpc_id == "" ? join(",", aws_subnet.master_subnet.*.id) :  join(",", data.aws_subnet.external_master.*.id))}"]
}

output "worker_subnet_ids" {
  value = ["${split(",", var.external_vpc_id == "" ? join(",", aws_subnet.worker_subnet.*.id) :  join(",", data.aws_subnet.external_worker.*.id))}"]
}

output "etcd_sg_id" {
  value = "${aws_security_group.etcd.id}"
}

output "master_sg_id" {
  value = "${aws_security_group.master.id}"
}

output "worker_sg_id" {
  value = "${aws_security_group.worker.id}"
}

output "api_sg_id" {
  value = "${aws_security_group.api.id}"
}

output "console_sg_id" {
  value = "${aws_security_group.console.id}"
}
