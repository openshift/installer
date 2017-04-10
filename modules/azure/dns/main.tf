resource "azurerm_dns_zone" "tectonic_azure_dns_zone" {
  name                = "${var.base_domain}"
  resource_group_name = "${var.resource_group_name}"
}
