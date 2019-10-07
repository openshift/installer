# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file

// Only reference data sources which are guaranteed to exist at any time (above) in this locals{} block
locals {
  master_subnet_cidr = cidrsubnet(var.vnet_cidr, 3, 0) #master subnet is a smaller subnet within the vnet. i.e from /21 to /24
  node_subnet_cidr   = cidrsubnet(var.vnet_cidr, 3, 1) #node subnet is a smaller subnet within the vnet. i.e from /21 to /24
}

data "azurerm_virtual_network" "cluster_vnet" {
  name                = azurerm_virtual_network.cluster_vnet.name
  resource_group_name = var.resource_group_name
}
