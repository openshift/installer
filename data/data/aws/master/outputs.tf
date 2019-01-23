output "cluster_id" {
  value = "${var.cluster_id}"
}

output "ip_addresses" {
  value = "${aws_instance.master.*.private_ip}"
}
