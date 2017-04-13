output "ip_address" {
  value = ["${azurerm_public_ip.tectonic_api_ip.ip_address}"]
}

output "console_ip_address" {
  value = "${azurerm_public_ip.tectonic_console_ip.ip_address}"
}

output "ingress_external_fqdn" {
  value = "${azurerm_public_ip.tectonic_console_ip.fqdn}"
}

output "ingress_internal_fqdn" {
  value = "${azurerm_public_ip.tectonic_console_ip.fqdn}"
}

output "api_external_fqdn" {
  value = "${azurerm_public_ip.tectonic_api_ip.fqdn}"
}

output "api_internal_fqdn" {
  value = "${azurerm_public_ip.tectonic_api_ip.fqdn}"
}
