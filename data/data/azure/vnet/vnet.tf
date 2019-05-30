resource "azurerm_route_table" "route_table" {
  name                = "${var.cluster_id}-node-routetable"
  location            = var.region
  resource_group_name = var.resource_group_name
}

resource "azurerm_subnet" "master_subnet" {
  resource_group_name  = var.resource_group_name
  address_prefix       = var.master_subnet_cidr
  virtual_network_name = var.vnet_name
  name                 = "${var.cluster_id}-controlplane-subnet"
}

resource "azurerm_subnet" "node_subnet" {
  resource_group_name  = var.resource_group_name
  address_prefix       = var.node_subnet_cidr
  virtual_network_name = var.vnet_name
  name                 = "${var.cluster_id}-node-subnet"
}

