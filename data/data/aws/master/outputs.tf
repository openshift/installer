output "ip_addresses" {
  value = aws_network_interface.master.*.private_ips
}

output "ipv6_addresses" {
  value = data.aws_network_interface.master.*.ipv6_addresses
}
