resource "azurestack_virtual_network" "cluster_vnet" {
  count = var.preexisting_network ? 0 : 1

  name                = var.virtual_network_name
  resource_group_name = var.resource_group_name
  location            = var.region
  address_space       = var.vnet_v4_cidrs
}

resource "azurestack_subnet" "master_subnet" {
  count = var.preexisting_network ? 0 : 1

  resource_group_name  = var.resource_group_name
  address_prefix       = local.master_subnet_cidr_v4
  virtual_network_name = local.virtual_network
  name                 = var.master_subnet
  // TODO: Figure out how to associate the subnet with the security group when the subnet is user-supplied.
  network_security_group_id = azurestack_network_security_group.cluster.id
}

resource "azurestack_subnet" "worker_subnet" {
  count = var.preexisting_network ? 0 : 1

  resource_group_name  = var.resource_group_name
  address_prefix       = local.worker_subnet_cidr_v4
  virtual_network_name = local.virtual_network
  name                 = var.worker_subnet
  // TODO: Figure out how to associate the subnet with the security group when the subnet is user-supplied.
  network_security_group_id = azurestack_network_security_group.cluster.id
}
