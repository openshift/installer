
resource "azurerm_dns_zone" "tectonic_azure_dns_zone" {
   name = "${var.tectonic_base_domain}"
   location = "${var.tectonic_region}"
   resource_group_name = "${azurerm_resource_group.tectonic_azure_dns_resource_group.name}"
}

