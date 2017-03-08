resource "azurerm_resource_group" "tectonic_azure_cluster_resource_group" {
   name = "tectonic-${var.tectonic_cluster_name}"
   location = "West US"
}
