resource "azurerm_resource_group" "tectonic_azure_cluster_resource_group" {
   location = "${var.tectonic_region}"
   name = "tectonic-cluster-${var.tectonic_cluster_name}-group"
}

resource "azurerm_resource_group" "tectonic_azure_dns_resource_group" {
   name = "${var.tectonic_azure_dns_resource_group}"
   location = "eastus"
}
