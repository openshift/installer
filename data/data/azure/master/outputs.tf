output "ip_v4_addresses" {
  value = var.use_ipv4 ? azurerm_network_interface.master.*.private_ip_address : []
}

output "ip_v6_addresses" {
  value = var.use_ipv6 ? azurerm_network_interface.master.*.private_ip_addresses.1 : []
}

