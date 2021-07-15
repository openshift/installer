# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file

data "azurestack_subnet" "preexisting_master_subnet" {
  count = var.azure_preexisting_network ? 1 : 0

  resource_group_name  = var.azure_network_resource_group_name
  virtual_network_name = var.azure_virtual_network
  name                 = var.azure_control_plane_subnet
}

data "azurestack_subnet" "preexisting_worker_subnet" {
  count = var.azure_preexisting_network ? 1 : 0

  resource_group_name  = var.azure_network_resource_group_name
  virtual_network_name = var.azure_virtual_network
  name                 = var.azure_compute_subnet
}

data "azurestack_virtual_network" "preexisting_virtual_network" {
  count = var.azure_preexisting_network ? 1 : 0

  resource_group_name = var.azure_network_resource_group_name
  name                = var.azure_virtual_network
}

// Only reference data sources which are guaranteed to exist at any time (above) in this locals{} block
locals {
  master_subnet_cidr_v4 = var.use_ipv4 ? cidrsubnet(var.machine_v4_cidrs[0], 1, 0) : null  #master subnet is a smaller subnet within the vnet. i.e from /16 to /17
  master_subnet_cidr_v6 = var.use_ipv6 ? cidrsubnet(var.machine_v6_cidrs[0], 16, 0) : null #master subnet is a smaller subnet within the vnet. i.e from /48 to /64

  worker_subnet_cidr_v4 = var.use_ipv4 ? cidrsubnet(var.machine_v4_cidrs[0], 1, 1) : null  #node subnet is a smaller subnet within the vnet. i.e from /16 to /17
  worker_subnet_cidr_v6 = var.use_ipv6 ? cidrsubnet(var.machine_v6_cidrs[0], 16, 1) : null #node subnet is a smaller subnet within the vnet. i.e from /48 to /64

  master_subnet_id = var.azure_preexisting_network ? data.azurestack_subnet.preexisting_master_subnet[0].id : azurestack_subnet.master_subnet[0].id
  worker_subnet_id = var.azure_preexisting_network ? data.azurestack_subnet.preexisting_worker_subnet[0].id : azurestack_subnet.worker_subnet[0].id

  virtual_network    = var.azure_preexisting_network ? data.azurestack_virtual_network.preexisting_virtual_network[0].name : azurestack_virtual_network.cluster_vnet[0].name
  virtual_network_id = var.azure_preexisting_network ? data.azurestack_virtual_network.preexisting_virtual_network[0].id : azurestack_virtual_network.cluster_vnet[0].id
}
