output "ip_addresses" {
  value = azurerm_network_interface.master.*.private_ip_address
}

output "ip_v6_addresses" {
  value = azurerm_network_interface.master.*.private_ip_addresses.1
}

