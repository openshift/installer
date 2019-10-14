resource "azurerm_virtual_network" "cluster_vnet" {
  count = var.preexisting_network ? 0 : 1

  name                = var.virtual_network_name
  resource_group_name = var.resource_group_name
  location            = var.region
  address_space       = [var.vnet_cidr]
}

resource "azurerm_route_table" "route_table" {
  name                = "${var.cluster_id}-node-routetable"
  location            = var.region
  resource_group_name = var.resource_group_name
}

resource "azurerm_subnet" "master_subnet" {
  count = var.preexisting_network ? 0 : 1

  resource_group_name  = var.resource_group_name
  address_prefix       = local.master_subnet_cidr
  virtual_network_name = local.virtual_network
  name                 = var.master_subnet
}

resource "azurerm_subnet" "worker_subnet" {
  count = var.preexisting_network ? 0 : 1

  resource_group_name  = var.resource_group_name
  address_prefix       = local.worker_subnet_cidr
  virtual_network_name = local.virtual_network
  name                 = var.worker_subnet
}
