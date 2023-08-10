# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file

data "azurerm_subnet" "preexisting_master_subnet" {
  count = var.azure_preexisting_network ? 1 : 0

  resource_group_name  = var.azure_network_resource_group_name
  virtual_network_name = var.azure_virtual_network
  name                 = var.azure_control_plane_subnet[0]
}

data "azurerm_subnet" "preexisting_worker_subnet" {
  count = var.azure_preexisting_network ? 1 : 0

  resource_group_name  = var.azure_network_resource_group_name
  virtual_network_name = var.azure_virtual_network
  name                 = var.azure_compute_subnet[0]
}

data "azurerm_virtual_network" "preexisting_virtual_network" {
  count = var.azure_preexisting_network ? 1 : 0

  resource_group_name = var.azure_network_resource_group_name
  name                = var.azure_virtual_network
}

// Only reference data sources which are guaranteed to exist at any time (above) in this locals{} block
locals {
  cluster_zones     = distinct(compact(concat(var.azure_master_availability_zones, var.azure_worker_availability_zones)))
  nat_gateway_count = var.azure_outbound_routing_type == "NatGateway" ? max(length(local.cluster_zones), 1) : 0

  master_zones                = distinct(compact(var.azure_master_availability_zones))
  master_subnet_count         = var.azure_outbound_routing_type == "NatGateway" ? max(length(local.master_zones), 1) : 1
  master_subnet_cidr_v4       = var.use_ipv4 ? cidrsubnet(var.machine_v4_cidrs[0], 1, 0) : null  #master subnet is a smaller subnet within the vnet. i.e from /16 to /17
  master_subnet_cidr_v6       = var.use_ipv6 ? cidrsubnet(var.machine_v6_cidrs[0], 16, 0) : null #master subnet is a smaller subnet within the vnet. i.e from /48 to /64
  master_subnet_cidr_v4_zonal = [for i in range(local.master_subnet_count) : cidrsubnet(local.master_subnet_cidr_v4, ceil(log(local.master_subnet_count, 2)), i)]

  worker_zones                = distinct(compact(var.azure_worker_availability_zones))
  worker_subnet_count         = var.azure_outbound_routing_type == "NatGateway" ? max(length(local.worker_zones), 1) : 1
  worker_subnet_cidr_v4       = var.use_ipv4 ? cidrsubnet(var.machine_v4_cidrs[0], 1, 1) : null  #node subnet is a smaller subnet within the vnet. i.e from /16 to /17
  worker_subnet_cidr_v6       = var.use_ipv6 ? cidrsubnet(var.machine_v6_cidrs[0], 16, 1) : null #node subnet is a smaller subnet within the vnet. i.e from /48 to /64
  worker_subnet_cidr_v4_zonal = [for i in range(local.worker_subnet_count) : cidrsubnet(local.worker_subnet_cidr_v4, ceil(log(local.worker_subnet_count, 2)), i)]

  master_subnet_id = var.azure_preexisting_network ? data.azurerm_subnet.preexisting_master_subnet[*].id : azurerm_subnet.master_subnet[*].id
  worker_subnet_id = var.azure_preexisting_network ? data.azurerm_subnet.preexisting_worker_subnet[*].id : azurerm_subnet.worker_subnet[*].id

  virtual_network    = var.azure_preexisting_network ? data.azurerm_virtual_network.preexisting_virtual_network[0].name : azurerm_virtual_network.cluster_vnet[0].name
  virtual_network_id = var.azure_preexisting_network ? data.azurerm_virtual_network.preexisting_virtual_network[0].id : azurerm_virtual_network.cluster_vnet[0].id
}
