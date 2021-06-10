data "azurerm_virtual_network" "cluster_vnet" {
  count = var.preexisting_network ? 0 : 1

  name                = var.virtual_network_name
  resource_group_name = var.resource_group_name
}

data "azurerm_subnet" "master_subnet" {
  count = var.preexisting_network ? 0 : 1

  resource_group_name  = var.resource_group_name
  virtual_network_name = local.virtual_network
  name                 = var.master_subnet
}

data "azurerm_subnet" "worker_subnet" {
  count = var.preexisting_network ? 0 : 1

  resource_group_name  = var.resource_group_name
  virtual_network_name = local.virtual_network
  name                 = var.worker_subnet
}
