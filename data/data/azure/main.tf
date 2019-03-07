# Configure the Azure Provider
provider "azurerm" {
  # whilst the `version` attribute is optional, we recommend pinning to a given version of the Provider
  version = "=1.22.0"
}

# Create a resource group
resource "azurerm_resource_group" "main" {
  name     = "${var.cluster_id}"
  location = "${var.region}"
}

# Create a virtual network within the resource group
resource "azurerm_virtual_network" "main" {
  name                = "${var.cluster_id}"
  resource_group_name = "${azurerm_resource_group.main.name}"
  location            = "${azurerm_resource_group.main.location}"
  address_space       = ["192.0.0.0/8"]
}
