resource "azurestack_virtual_network" "cluster_vnet" {
  count = var.azure_preexisting_network ? 0 : 1

  name                = var.azure_virtual_network
  resource_group_name = data.azurestack_resource_group.main.name
  location            = var.azure_region
  address_space       = var.machine_v4_cidrs
}

resource "azurestack_subnet" "master_subnet" {
  count = var.azure_preexisting_network ? 0 : 1

  resource_group_name       = data.azurestack_resource_group.main.name
  address_prefix            = local.master_subnet_cidr_v4
  virtual_network_name      = local.virtual_network
  name                      = var.azure_control_plane_subnet
  network_security_group_id = azurestack_network_security_group.cluster.id
}

resource "azurestack_subnet" "worker_subnet" {
  count = var.azure_preexisting_network ? 0 : 1

  resource_group_name       = data.azurestack_resource_group.main.name
  address_prefix            = local.worker_subnet_cidr_v4
  virtual_network_name      = local.virtual_network
  name                      = var.azure_compute_subnet
  network_security_group_id = azurestack_network_security_group.cluster.id
}
