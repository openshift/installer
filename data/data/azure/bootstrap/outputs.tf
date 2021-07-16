output "bootstrap_ip" {
  value = var.azure_private ? azurerm_network_interface.bootstrap.private_ip_address : azurerm_public_ip.bootstrap_public_ip_v4[0].ip_address
}
