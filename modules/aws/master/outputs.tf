output "subnet_ids" {
  value = "${var.subnet_ids}"
}

output "aws_lbs" {
  value = "${var.aws_lbs}"
}

output "cluster_id" {
  value = "${var.cluster_id}"
}

output "ip_addresses" {
  value = "${aws_instance.master.*.private_ip}"
}
