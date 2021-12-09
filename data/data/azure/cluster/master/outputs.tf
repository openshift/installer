output "ip_addresses" {
  value = azurerm_network_interface.master.*.private_ip_address
}
