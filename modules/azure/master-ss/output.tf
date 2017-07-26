output "ip_address" {
  value = ["${azurerm_public_ip.api_ip.ip_address}"]
}

output "console_ip_addresses" {
  value = "${azurerm_public_ip.console_ip.ip_address}"
}

output "ingress_fqdn" {
  value = "${var.base_domain == "" ? "${azurerm_public_ip.console_ip.domain_name_label}.${var.base_domain}" : azurerm_public_ip.console_ip.fqdn}"
}

output "ingress_fqdn" {
  value = "${var.base_domain == "" ? "${azurerm_public_ip.console_ip.domain_name_label}.${var.base_domain}" : azurerm_public_ip.console_ip.fqdn}"
}

output "api_fqdn" {
  value = "${var.base_domain == "" ?  "${azurerm_public_ip.api_ip.domain_name_label}.${var.base_domain}" : azurerm_public_ip.api_ip.fqdn}"
}

output "api_fqdn" {
  value = "${var.base_domain == "" ?  "${azurerm_public_ip.api_ip.domain_name_label}.${var.base_domain}" : azurerm_public_ip.api_ip.fqdn}"
}
