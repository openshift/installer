resource "azurerm_resource_group" "tectonic_azure_cluster_resource_group" {
   name = "tectonic-cluster-${var.tectonic_cluster_name}-group"
   location = "eastus"
}

resource "azurerm_resource_group" "tectonic_azure_dns_resource_group" {
   name = "${var.tectonic_azure_dns_resource_group}"
   location = "eastus"
}
