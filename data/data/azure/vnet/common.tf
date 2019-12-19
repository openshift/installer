# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file

data "azurerm_subnet" "preexisting_master_subnet" {
  count = var.preexisting_network ? 1 : 0

  resource_group_name  = var.network_resource_group_name
  virtual_network_name = var.virtual_network_name
  name                 = var.master_subnet
}

data "azurerm_subnet" "preexisting_worker_subnet" {
  count = var.preexisting_network ? 1 : 0

  resource_group_name  = var.network_resource_group_name
  virtual_network_name = var.virtual_network_name
  name                 = var.worker_subnet
}

data "azurerm_virtual_network" "preexisting_virtual_network" {
  count = var.preexisting_network ? 1 : 0

  resource_group_name = var.network_resource_group_name
  name                = var.virtual_network_name
}

// Only reference data sources which are guaranteed to exist at any time (above) in this locals{} block
locals {
  master_subnet_cidr_v4 = var.use_ipv4 ? cidrsubnet(var.vnet_v4_cidrs[0], 3, 0) : null  #master subnet is a smaller subnet within the vnet. i.e from /21 to /24
  master_subnet_cidr_v6 = var.use_ipv6 ? cidrsubnet(var.vnet_v6_cidrs[0], 16, 0) : null #master subnet is a smaller subnet within the vnet. i.e from /48 to /64

  worker_subnet_cidr_v4 = var.use_ipv4 ? cidrsubnet(var.vnet_v4_cidrs[0], 3, 1) : null  #node subnet is a smaller subnet within the vnet. i.e from /21 to /24
  worker_subnet_cidr_v6 = var.use_ipv6 ? cidrsubnet(var.vnet_v6_cidrs[0], 16, 1) : null #node subnet is a smaller subnet within the vnet. i.e from /48 to /64

  master_subnet_id = var.preexisting_network ? data.azurerm_subnet.preexisting_master_subnet[0].id : azurerm_subnet.master_subnet[0].id
  worker_subnet_id = var.preexisting_network ? data.azurerm_subnet.preexisting_worker_subnet[0].id : azurerm_subnet.worker_subnet[0].id

  virtual_network    = var.preexisting_network ? data.azurerm_virtual_network.preexisting_virtual_network[0].name : azurerm_virtual_network.cluster_vnet[0].name
  virtual_network_id = var.preexisting_network ? data.azurerm_virtual_network.preexisting_virtual_network[0].id : azurerm_virtual_network.cluster_vnet[0].id
}
