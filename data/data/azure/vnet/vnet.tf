resource "azurerm_virtual_network" "cluster_vnet" {
  count = var.preexisting_network ? 0 : 1

  name                = var.virtual_network_name
  resource_group_name = var.resource_group_name
  location            = var.region
  address_space       = concat(var.vnet_v4_cidrs, var.vnet_v6_cidrs)
}

resource "azurerm_subnet" "master_subnet" {
  count = var.preexisting_network ? 0 : 1

  resource_group_name = var.resource_group_name
  address_prefixes = [for cidr in [
    { value : local.master_subnet_cidr_v4, include : var.use_ipv4 },
    { value : local.master_subnet_cidr_v6, include : var.use_ipv6 }
  ] : cidr.value if cidr.include]
  virtual_network_name = local.virtual_network
  name                 = var.master_subnet
}

resource "azurerm_subnet" "worker_subnet" {
  count = var.preexisting_network ? 0 : 1

  resource_group_name = var.resource_group_name
  address_prefixes = [for cidr in [
    { value : local.worker_subnet_cidr_v4, include : var.use_ipv4 },
    { value : local.worker_subnet_cidr_v6, include : var.use_ipv6 }
  ] : cidr.value if cidr.include]
  virtual_network_name = local.virtual_network
  name                 = var.worker_subnet
}
