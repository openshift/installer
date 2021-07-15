output "ip_addresses" {
  value = azurestack_network_interface.master.*.private_ip_addresses
}
