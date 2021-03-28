# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file

data "azurestack_subnet" "preexisting_master_subnet" {
  count = var.preexisting_network ? 1 : 0

  resource_group_name  = var.network_resource_group_name
  virtual_network_name = var.virtual_network_name
  name                 = var.master_subnet
}

data "azurestack_subnet" "preexisting_worker_subnet" {
  count = var.preexisting_network ? 1 : 0

  resource_group_name  = var.network_resource_group_name
  virtual_network_name = var.virtual_network_name
  name                 = var.worker_subnet
}

data "azurestack_virtual_network" "preexisting_virtual_network" {
  count = var.preexisting_network ? 1 : 0

  resource_group_name = var.network_resource_group_name
  name                = var.virtual_network_name
}

// Only reference data sources which are guaranteed to exist at any time (above) in this locals{} block
locals {
  master_subnet_cidr_v4 = cidrsubnet(var.vnet_v4_cidrs[0], 1, 0) #master subnet is a smaller subnet within the vnet. i.e from /21 to /22

  worker_subnet_cidr_v4 = cidrsubnet(var.vnet_v4_cidrs[0], 1, 1) #node subnet is a smaller subnet within the vnet. i.e from /21 to /22

  master_subnet_id = var.preexisting_network ? data.azurestack_subnet.preexisting_master_subnet[0].id : azurestack_subnet.master_subnet[0].id
  worker_subnet_id = var.preexisting_network ? data.azurestack_subnet.preexisting_worker_subnet[0].id : azurestack_subnet.worker_subnet[0].id

  virtual_network    = var.preexisting_network ? data.azurestack_virtual_network.preexisting_virtual_network[0].name : azurestack_virtual_network.cluster_vnet[0].name
  virtual_network_id = var.preexisting_network ? data.azurestack_virtual_network.preexisting_virtual_network[0].id : azurestack_virtual_network.cluster_vnet[0].id

  description = "Created By OpenShift Installer"
}
