output "ip_addresses" {
  value = "${aws_instance.master.*.private_ip}"
}
