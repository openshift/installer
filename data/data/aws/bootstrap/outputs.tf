output "ip_addresses" {
  value = aws_network_interface.bootstrap.*.private_ips
}

