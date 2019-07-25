output "ip_addresses" {
  value = [azurerm_public_ip.bootstrap_public_ip.private_ip_address]
}

