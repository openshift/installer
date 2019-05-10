output "ip_addresses" {
  value = aws_network_interface.master.*.private_ips
}

